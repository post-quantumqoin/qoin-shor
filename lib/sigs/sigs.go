package sigs

import (
	"context"
	"fmt"

	"go.opencensus.io/trace"
	"golang.org/x/xerrors"

	"github.com/post-quantumqoin/address"
	"github.com/post-quantumqoin/core-types/crypto"

	// pqccrypto "github.com/post-quantumqoin/qoin-shor/pqccrypto"

	"github.com/post-quantumqoin/qoin-shor/chain/types"
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
	if sig == nil {
		return xerrors.Errorf("signature is nil")
	}

	if addr.Protocol() == address.ID {
		return fmt.Errorf("must resolve ID addresses before using them to verify a signature")
	}
	fmt.Println("MultiPqcVerify msg:", msg)
	for _, pk := range sig.PqcCert.Pubkeys {
		fmt.Println("MultiPqcVerify Typ:", crypto.GetTypeByName(pk.Typ))
		sig, err := findSign(crypto.GetTypeByName(pk.Typ), sig.PqcSignatures)
		if err != nil {
			return err
		}
		fmt.Println("MultiPqcVerify len(sig):", len(sig))
		err = PqcVerify(crypto.GetTypeByName(pk.Typ), sig, pk.Pubkey, msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func findSign(typ crypto.SigType, signatures []crypto.PqcSignature) ([]byte, error) {
	for _, sg := range signatures {
		if typ == sg.Type {
			return sg.Data, nil
		}
	}
	return nil, xerrors.Errorf("fail to find a valid signature")
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
