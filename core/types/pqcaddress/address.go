package pqcaddress

import (
	"bytes"
	"encoding/binary"

	"golang.org/x/crypto/sha3"
	"golang.org/x/xerrors"

	"github.com/post-quantumqoin/address"
	builtintypes "github.com/post-quantumqoin/core-types/builtin"
	typescrypto "github.com/post-quantumqoin/core-types/crypto"

	"github.com/post-quantumqoin/qoin-shor/core/types"
)

// Scheme defines how the address key is selected and encoded. The value is
// stored in SignPQCCert.Version to preserve the existing wire structure.
type Scheme uint8

const (
	// SchemeLegacyCompatible keeps the certificate wire order, uses the first
	// key, and accepts both CBOR and raw encodings. This is only for addresses
	// created by legacy shor wallet and Qwallet releases.
	SchemeLegacyCompatible Scheme = iota

	// SchemeCanonicalCBOR requires canonical certificate order and derives the
	// address from the CBOR-serialized first certificate key. All new wallets
	// and signatures must use this scheme.
	SchemeCanonicalCBOR
)

// ValidateCertificate validates the certificate structure and, for modern
// schemes, requires strict lexicographic key-type ordering.
func ValidateCertificate(cert typescrypto.SignPQCCert) error {
	if len(cert.Pubkeys) == 0 {
		return xerrors.New("PQC certificate contains no public keys")
	}

	scheme := Scheme(cert.Version)
	if scheme > SchemeCanonicalCBOR {
		return xerrors.Errorf("unsupported PQC address scheme: %d", cert.Version)
	}

	seen := make(map[typescrypto.SigType]struct{}, len(cert.Pubkeys))
	for index, publicKey := range cert.Pubkeys {
		typ := typescrypto.GetTypeByName(publicKey.Typ)
		if !isPqcAlgorithmType(typ) {
			return xerrors.Errorf("unsupported PQC certificate key type at index %d: %q", index, publicKey.Typ)
		}
		if len(publicKey.Pubkey) == 0 {
			return xerrors.Errorf("empty PQC certificate public key at index %d", index)
		}
		if _, ok := seen[typ]; ok {
			return xerrors.Errorf("duplicate PQC certificate key type: %q", publicKey.Typ)
		}
		seen[typ] = struct{}{}

		if scheme != SchemeLegacyCompatible && index > 0 && cert.Pubkeys[index-1].Typ >= publicKey.Typ {
			return xerrors.Errorf("PQC certificate public keys are not in canonical order")
		}
	}

	return nil
}

// SignerAddressCandidates returns the only addresses accepted by the selected
// certificate scheme. Legacy scheme 0 returns at most two candidates for the
// same first key (CBOR first, raw second); modern schemes return exactly one.
func SignerAddressCandidates(cert typescrypto.SignPQCCert, protocol address.Protocol) ([]address.Address, error) {
	if err := ValidateCertificate(cert); err != nil {
		return nil, err
	}

	first := cert.Pubkeys[0]
	serialized, err := serializeCertPublicKey(first, Scheme(cert.Version))
	if err != nil {
		return nil, xerrors.Errorf("serialize first PQC certificate public key: %w", err)
	}

	var payloads [][]byte
	switch Scheme(cert.Version) {
	case SchemeLegacyCompatible:
		payloads = [][]byte{serialized, first.Pubkey}
	case SchemeCanonicalCBOR:
		payloads = [][]byte{serialized}
	default:
		return nil, xerrors.Errorf("unsupported PQC address scheme: %d", cert.Version)
	}

	result := make([]address.Address, 0, len(payloads))
	for _, payload := range payloads {
		candidate, err := fromBytes(payload, protocol)
		if err != nil {
			return nil, err
		}
		if !Contains(candidate, result) {
			result = append(result, candidate)
		}
	}

	return result, nil
}

// Contains reports whether expected is one of the address candidates.
func Contains(expected address.Address, candidates []address.Address) bool {
	for _, candidate := range candidates {
		if candidate == expected {
			return true
		}
	}
	return false
}

func serializeCertPublicKey(publicKey typescrypto.SignPqcCertPubkey, scheme Scheme) ([]byte, error) {
	switch scheme {
	case SchemeLegacyCompatible:
		value := types.PqcCertPubkey{Typ: publicKey.Typ, Pubkey: publicKey.Pubkey}
		return value.Serialize()
	case SchemeCanonicalCBOR:
		value := types.PqcCertPubkey{Typ: publicKey.Typ, Pubkey: publicKey.Pubkey}
		return value.Serialize()
	default:
		return nil, xerrors.Errorf("unsupported PQC address scheme: %d", scheme)
	}
}

func fromBytes(publicKey []byte, protocol address.Protocol) (address.Address, error) {
	switch protocol {
	case address.Delegated:
		hasher := sha3.NewLegacyKeccak256()
		_, _ = hasher.Write(publicKey)
		ethAddress := hasher.Sum(nil)[12:]

		maskedIDPrefix := [12]byte{0xff}
		if bytes.HasPrefix(ethAddress, maskedIDPrefix[:]) {
			return address.NewIDAddress(binary.BigEndian.Uint64(ethAddress[12:]))
		}
		result, err := address.NewDelegatedAddress(builtintypes.EthereumAddressManagerActorID, ethAddress)
		if err != nil {
			return address.Undef, xerrors.Errorf("derive delegated PQC address: %w", err)
		}
		return result, nil
	case address.PQC:
		result, err := address.NewPqcAddress(publicKey)
		if err != nil {
			return address.Undef, xerrors.Errorf("derive PQC address: %w", err)
		}
		return result, nil
	default:
		return address.Undef, xerrors.Errorf("unsupported PQC signer address protocol: %d", protocol)
	}
}

func isPqcAlgorithmType(typ typescrypto.SigType) bool {
	switch typ {
	case typescrypto.SigTypeFalcon512,
		typescrypto.SigTypeFalcon1024,
		typescrypto.SigTypeDilithium3,
		typescrypto.SigTypeDilithium5:
		return true
	default:
		return false
	}
}
