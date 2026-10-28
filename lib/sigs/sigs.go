package sigs

import (
	"context"
	"fmt"

	"go.opencensus.io/trace"
	"golang.org/x/xerrors"

	"github.com/post-quantumqoin/address"
	"github.com/post-quantumqoin/core-types/crypto"

	// pqccrypto "github.com/post-quantumqoin/qoin-shor/pqccrypto"

	"github.com/post-quantumqoin/qoin-shor/core/types"
	"github.com/post-quantumqoin/qoin-shor/core/types/pqcaddress"
)

// Sign takes in signature type, private key and message. Returns a signature for that message.
// Valid sigTypes are: "secp256k1" and "bls"
func Sign(sigType crypto.SigType, privkey []byte, msg []byte) (*crypto.Signature, error) {
	sv, ok := sigs[sigType]
	if !ok {
		return nil, fmt.Errorf("cannot sign message with signature of unsupported type: %v", sigType)
	}
	//bls secp
	sb, err := sv.Sign(privkey, msg)
	if err != nil {
		return nil, err
	}
	return &crypto.Signature{
		Type: sigType,
		Data: sb,
	}, nil
}

// Verify verifies signatures
func Verify(sig *crypto.Signature, addr address.Address, msg []byte) error {
	if sig == nil {
		return xerrors.Errorf("signature is nil")
	}

	if addr.Protocol() == address.ID {
		return fmt.Errorf("must resolve ID addresses before using them to verify a signature")
	}

	sv, ok := sigs[sig.Type]
	if !ok {
		return fmt.Errorf("cannot verify signature of unsupported type: %v", sig.Type)
	}

	return sv.Verify(sig.Data, addr, msg)
}

// Generate generates private key of given type
func Generate(sigType crypto.SigType) ([]byte, error) {
	sv, ok := sigs[sigType]
	if !ok {
		return nil, fmt.Errorf("cannot generate private key of unsupported type: %v", sigType)
	}

	return sv.GenPrivate()
}

// ToPublic converts private key to public key
func ToPublic(sigType crypto.SigType, pk []byte) ([]byte, error) {
	sv, ok := sigs[sigType]
	if !ok {
		return nil, fmt.Errorf("cannot generate public key of unsupported type: %v", sigType)
	}

	return sv.ToPublic(pk)
}

func CheckBlockSignature(ctx context.Context, blk *types.BlockHeader, worker address.Address) error {
	_, span := trace.StartSpan(ctx, "checkBlockSignature")
	defer span.End()

	if blk.IsValidated() {
		return nil
	}

	if blk.BlockSig == nil {
		return xerrors.New("block signature not present")
	}

	sigb, err := blk.SigningBytes()
	if err != nil {
		return xerrors.Errorf("failed to get block signing bytes: %w", err)
	}

	// err = Verify(blk.BlockSig, worker, sigb)
	// if err == nil {
	// 	blk.SetValidated()
	// }
	err = MultiPqcVerify(blk.BlockSig, worker, sigb)
	if err == nil {
		blk.SetValidated()
	}

	return err
}

// SigShim is used for introducing signature functions
type SigShim interface {
	GenPrivate() ([]byte, error)
	ToPublic(pk []byte) ([]byte, error)
	Sign(pk []byte, msg []byte) ([]byte, error)
	Verify(sig []byte, a address.Address, msg []byte) error
}

var sigs map[crypto.SigType]SigShim

// RegisterSignature should be only used during init
func RegisterSignature(typ crypto.SigType, vs SigShim) {
	if sigs == nil {
		sigs = make(map[crypto.SigType]SigShim)
	}
	sigs[typ] = vs
}

func toSignPqcCert(cert types.PQCCert) crypto.SignPQCCert {
	var pkeys []crypto.SignPqcCertPubkey
	for _, pk := range cert.Pubkeys {
		pkeys = append(pkeys, crypto.SignPqcCertPubkey{Typ: pk.Typ, Pubkey: pk.Pubkey})
	}

	signcert := crypto.SignPQCCert{
		Pubkeys: pkeys,
		Version: cert.Version,
	}
	return signcert
}

func MultiPqcSign(cert types.PQCCert, kps []types.PqcKeypair, msg []byte) (*crypto.Signature, error) {
	// signs []crypto.PqcSignature
	var pqcSigns []crypto.PqcSignature
	for _, kp := range kps {
		fmt.Println("MultiPqcSign kp.PqcType:", kp.PqcType)
		sgin, err := PqcSign(crypto.GetTypeByName(string(kp.PqcType)), kp.PqcPrivateKey, msg)
		if err != nil {
			return nil, err
		}
		pqcSigns = append(pqcSigns, crypto.PqcSignature{Type: crypto.GetTypeByName(string(kp.PqcType)), Data: sgin})
	}

	pqcsign := &crypto.Signature{
		Type:          crypto.SigTypeMultiPqc,
		PqcCert:       toSignPqcCert(cert),
		PqcSignatures: pqcSigns,
	}
	return pqcsign, nil
}

func PqcSign(typ crypto.SigType, privkey []byte, msg []byte) ([]byte, error) {
	sv, ok := Pqcsigs[typ]
	if !ok {
		return nil, fmt.Errorf("cannot sign message with signature of unsupported type: %v", typ)
	}

	signature, err := sv.PqcSign(privkey, msg)
	if err != nil {
		return nil, err
	}
	return signature, nil
}

func MultiPqcVerify(sig *crypto.Signature, addr address.Address, msg []byte) error {
	return MultiPqcVerifyWithAddressBinding(sig, addr, msg, true)
}

// MultiPqcVerifyWithAddressBinding supports historical block authentication.
// Only address matching is optional; envelope and cryptographic checks always run.
func MultiPqcVerifyWithAddressBinding(sig *crypto.Signature, addr address.Address, msg []byte, requireAddressMatch bool) error {
	if sig == nil {
		return xerrors.Errorf("signature is nil")
	}

	if addr.Protocol() == address.ID {
		return fmt.Errorf("must resolve ID addresses before using them to verify a signature")
	}

	signatures, err := validatePqcSignatureEnvelope(sig)
	if err != nil {
		return err
	}

	for _, publicKey := range sig.PqcCert.Pubkeys {
		typ := crypto.GetTypeByName(publicKey.Typ)
		if err := PqcVerify(typ, signatures[typ], publicKey.Pubkey, msg); err != nil {
			return xerrors.Errorf("verify PQC signature type %s: %w", publicKey.Typ, err)
		}
	}

	if requireAddressMatch {
		if sig.PqcCert.Version == uint8(pqcaddress.SchemeLegacyCompatible) {
			matched, err := pqcaddress.LegacyCertificateMatchesAddress(sig.PqcCert, addr)
			if err != nil {
				return xerrors.Errorf("derive signer address from legacy PQC certificate: %w", err)
			}
			if !matched {
				return xerrors.Errorf("signer address %s does not match any legacy PQC certificate public key", addr)
			}
		} else {
			candidates, err := pqcaddress.SignerAddressCandidates(sig.PqcCert, addr.Protocol())
			if err != nil {
				return xerrors.Errorf("derive signer address from PQC certificate: %w", err)
			}
			if !pqcaddress.Contains(addr, candidates) {
				return xerrors.Errorf("signer address %s does not match the first PQC certificate public key", addr)
			}
		}
	}

	return nil
}

func validatePqcSignatureEnvelope(sig *crypto.Signature) (map[crypto.SigType][]byte, error) {
	if sig.Type != crypto.SigTypeDelegated && sig.Type != crypto.SigTypeMultiPqc {
		return nil, xerrors.Errorf("invalid outer PQC signature type: %d", sig.Type)
	}
	if len(sig.Data) != 0 {
		return nil, xerrors.Errorf("PQC signature envelope must not contain legacy signature data")
	}
	if err := pqcaddress.ValidateCertificate(sig.PqcCert); err != nil {
		return nil, err
	}
	if len(sig.PqcCert.Pubkeys) != len(sig.PqcSignatures) {
		return nil, xerrors.Errorf(
			"PQC certificate/signature count mismatch: %d != %d",
			len(sig.PqcCert.Pubkeys),
			len(sig.PqcSignatures),
		)
	}

	certificateTypes := make(map[crypto.SigType]struct{}, len(sig.PqcCert.Pubkeys))
	for _, publicKey := range sig.PqcCert.Pubkeys {
		certificateTypes[crypto.GetTypeByName(publicKey.Typ)] = struct{}{}
	}

	signatures := make(map[crypto.SigType][]byte, len(sig.PqcSignatures))
	for index, signature := range sig.PqcSignatures {
		if _, ok := certificateTypes[signature.Type]; !ok {
			return nil, xerrors.Errorf("PQC signature at index %d has no matching certificate key", index)
		}
		if len(signature.Data) == 0 {
			return nil, xerrors.Errorf("empty PQC signature at index %d", index)
		}
		if _, ok := signatures[signature.Type]; ok {
			return nil, xerrors.Errorf("duplicate PQC signature type: %d", signature.Type)
		}
		signatures[signature.Type] = signature.Data
	}

	return signatures, nil
}

func PqcVerify(typ crypto.SigType, sig []byte, pk []byte, msg []byte) error {
	if sig == nil {
		return xerrors.Errorf("signature is nil")
	}

	sv, ok := Pqcsigs[typ]
	if !ok {
		return fmt.Errorf("cannot verify signature of unsupported type: %v", typ)
	}

	return sv.PqcVerify(sig, pk, msg)
}

func PqcGenerate(typ crypto.SigType) ([]byte, []byte, []byte, error) {
	// fmt.Println("sigs.PqcGenerate  sigType:",sigType)
	sv, ok := Pqcsigs[typ]
	if !ok {
		return nil, nil, nil, fmt.Errorf("cannot generate private key of unsupported type: %v", typ)
	}
	// fmt.Println("sigs.PqcGenPrivate")
	return sv.PqcGenPrivate()
}

func PqcToPublic(typ crypto.SigType, seed []byte) ([]byte, error) {
	sv, ok := Pqcsigs[typ]
	if !ok {
		return nil, fmt.Errorf("cannot generate pqc public key of unsupported type: %v", typ)
	}

	return sv.PqcToPublic(seed)
}

// SigShim is used for introducing signature functions
type PqcSigShim interface {
	PqcGenPrivate() ([]byte, []byte, []byte, error)
	PqcToPublic(seed []byte) ([]byte, error)
	PqcSign(sK []byte, msg []byte) ([]byte, error)
	PqcVerify(sig []byte, pk []byte, msg []byte) error
}

var Pqcsigs map[crypto.SigType]PqcSigShim

// RegisterSignature should be only used during init
func PqcRegisterSignature(typ crypto.SigType, vs PqcSigShim) {
	if Pqcsigs == nil {
		Pqcsigs = make(map[crypto.SigType]PqcSigShim)
	}
	Pqcsigs[typ] = vs
}
