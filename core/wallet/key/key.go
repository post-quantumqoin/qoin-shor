package key

import (
	// "fmt"
	"sort"
	// "strings"

	"golang.org/x/xerrors"

	logging "github.com/ipfs/go-log/v2"
	// "github.com/post-quantumqoin/qoin-shor/lib/address"
	"github.com/post-quantumqoin/address"
	"github.com/post-quantumqoin/core-types/crypto"

	"github.com/post-quantumqoin/qoin-shor/core/types"
	"github.com/post-quantumqoin/qoin-shor/core/types/ethtypes"
	"github.com/post-quantumqoin/qoin-shor/core/types/pqcaddress"
	"github.com/post-quantumqoin/qoin-shor/lib/sigs"
	// pqccrypto "github.com/post-quantumqoin/qoin-shor/pqccrypto"
)

var log = logging.Logger("key")

func GenerateKey(typ types.KeyType) (*Key, error) {
	ctyp := ActSigType(typ)
	if ctyp == crypto.SigTypeUnknown {
		return nil, xerrors.Errorf("unknown sig type: %s", typ)
	}
	pk, err := sigs.Generate(ctyp)
	if err != nil {
		return nil, err
	}
	ki := types.KeyInfo{
		Type:       typ,
		PrivateKey: pk,
	}
	return NewKey(ki)
}

type Key struct {
	types.KeyInfo

	PublicKey []byte
	Address   address.Address
}

func NewKey(keyinfo types.KeyInfo) (*Key, error) {
	k := &Key{
		KeyInfo: keyinfo,
	}

	var err error
	k.PublicKey, err = sigs.ToPublic(ActSigType(k.Type), k.PrivateKey)
	if err != nil {
		return nil, err
	}

	switch k.Type {
	case types.KTSecp256k1:
		k.Address, err = address.NewSecp256k1Address(k.PublicKey)
		if err != nil {
			return nil, xerrors.Errorf("converting Secp256k1 to address: %w", err)
		}
	case types.KTDelegated:
		// Transitory Delegated signature verification as per FIP-0055
		ethAddr, err := ethtypes.EthAddressFromPubKey(k.PublicKey)
		if err != nil {
			return nil, xerrors.Errorf("failed to calculate Eth address from public key: %w", err)
		}

		ea, err := ethtypes.CastEthAddress(ethAddr)
		if err != nil {
			return nil, xerrors.Errorf("failed to create ethereum address from bytes: %w", err)
		}

		k.Address, err = ea.ToQoinAddress()
		if err != nil {
			return nil, xerrors.Errorf("converting Delegated to address: %w", err)
		}
	case types.KTBLS:
		k.Address, err = address.NewBLSAddress(k.PublicKey)
		if err != nil {
			return nil, xerrors.Errorf("converting BLS to address: %w", err)
		}
	default:
		return nil, xerrors.Errorf("unsupported key type: %s", k.Type)
	}

	return k, nil

}

func ActSigType(typ types.KeyType) crypto.SigType {
	switch typ {
	case types.KTBLS:
		return crypto.SigTypeBLS
	case types.KTSecp256k1:
		return crypto.SigTypeSecp256k1
	default:
		return crypto.SigTypeUnknown
	}
}

// //type SigType byte
func PqcActSigType(sa types.SigAlg) crypto.SigType {
	switch sa {
	case types.Falcon512:
		return crypto.SigTypeFalcon512
	case types.Falcon1024:
		return crypto.SigTypeFalcon1024
	case types.Dilithium3:
		return crypto.SigTypeDilithium3
	case types.Dilithium5:
		return crypto.SigTypeDilithium5
	default:
		return crypto.SigTypeUnknown
	}
}

func GeneratePqcKeyWithAlgs(typ types.KeyType, algs []types.SigAlg) (*PqcKey, error) {
	if len(algs) == 0 {
		// Default to all supported algorithms if none specified
		algs = []types.SigAlg{types.Falcon512, types.Dilithium3}
	}
	// Validate algorithms and remove duplicates
	seen := map[types.SigAlg]struct{}{}
	var list []types.SigAlg
	for _, a := range algs {
		if _, ok := seen[a]; ok {
			continue
		}
		switch a {
		case types.Falcon512, types.Falcon1024, types.Dilithium3, types.Dilithium5:
			// ok
		default:
			return nil, xerrors.Errorf("unknown pqc algorithm: %s", a)
		}
		seen[a] = struct{}{}
		list = append(list, a)
	}

	var kprs []types.PqcKeypair
	for _, a := range list {
		ctyp := PqcActSigType(a)
		seed, sk, pk, err := sigs.PqcGenerate(ctyp)
		if err != nil {
			return nil, xerrors.Errorf("generate %s failed: %w", a, err)
		}
		kprs = append(kprs, types.PqcKeypair{
			PqcVersion:    0,
			PqcSeed:       seed,
			PqcType:       types.KeyType(a),
			PqcPrivateKey: sk,
			PqcPublicKey:  pk,
		})
	}

	ki := types.KeyInfo{
		Type:        typ, // This should be set to a type that indicates it's a multi-algorithm key, e.g., "pqc-multi"
		PqcKeypairs: kprs,
	}
	version := uint8(pqcaddress.SchemeCanonicalCBOR)
	ki.PqcAddressVersion = &version

	return NewPqcKey(ki)
}

type PqcKey struct {
	types.KeyInfo
	types.PQCCert

	// Hcert   []byte
	Address address.Address
}

func NewPqcKey(keyInfo types.KeyInfo) (*PqcKey, error) {
	if len(keyInfo.PqcKeypairs) == 0 {
		return nil, xerrors.Errorf("no pqc keypairs in KeyInfo")
	}
	// Missing metadata identifies legacy wallets; never reinterpret them as version 1.
	version := uint8(pqcaddress.SchemeLegacyCompatible)
	if keyInfo.PqcAddressVersion != nil {
		version = *keyInfo.PqcAddressVersion
	}
	if version > uint8(pqcaddress.SchemeCanonicalCBOR) {
		return nil, xerrors.Errorf("unsupported PQC address scheme: %d", version)
	}
	kpairs := make([]types.PqcKeypair, len(keyInfo.PqcKeypairs))
	copy(kpairs, keyInfo.PqcKeypairs)
	if version == uint8(pqcaddress.SchemeCanonicalCBOR) {
		sort.Slice(kpairs, func(i, j int) bool {
			return string(kpairs[i].PqcType) < string(kpairs[j].PqcType)
		})
	}

	var pubkeys []types.PqcCertPubkey
	for _, kp := range kpairs {
		if len(kp.PqcPublicKey) == 0 {
			return nil, xerrors.Errorf("empty public key for pqc type %s", kp.PqcType)
		}
		pubkeys = append(pubkeys, types.PqcCertPubkey{
			Typ:    string(kp.PqcType),
			Pubkey: kp.PqcPublicKey,
		})
	}

	cert := types.PQCCert{
		Pubkeys: pubkeys,
		Version: version,
	}

	pk := &PqcKey{
		KeyInfo: keyInfo,
		PQCCert: cert,
	}

	// ensure at least one pubkey exists
	if len(pubkeys) == 0 {
		return nil, xerrors.Errorf("no cert pubkeys built")
	}

	signPubkeys := make([]crypto.SignPqcCertPubkey, len(pubkeys))
	for index, publicKey := range pubkeys {
		signPubkeys[index] = crypto.SignPqcCertPubkey{
			Typ:    publicKey.Typ,
			Pubkey: publicKey.Pubkey,
		}
	}

	addressProtocol := address.PQC
	if keyInfo.Type == types.KTDelegated {
		addressProtocol = address.Delegated
	}
	addresses, err := pqcaddress.SignerAddressCandidates(crypto.SignPQCCert{
		Pubkeys: signPubkeys,
		Version: cert.Version,
	}, addressProtocol)
	if err != nil {
		return nil, xerrors.Errorf("derive PQC wallet address: %w", err)
	}
	pk.Address = addresses[0]
	if keyInfo.PqcAddress != "" {
		expected, err := address.NewFromString(keyInfo.PqcAddress)
		if err != nil {
			return nil, xerrors.Errorf("invalid stored PQC address: %w", err)
		}
		if !pqcaddress.Contains(expected, addresses) {
			return nil, xerrors.Errorf("stored PQC address %s does not match certificate version %d and first public key", expected, version)
		}
		pk.Address = expected
	}
	pk.KeyInfo.PqcAddressVersion = &version
	pk.KeyInfo.PqcAddress = pk.Address.String()

	return pk, nil
}

// NewPqcKeyForAddress binds a keystore entry to the address used to load it.
func NewPqcKeyForAddress(keyInfo types.KeyInfo, expected address.Address) (*PqcKey, error) {
	if keyInfo.PqcAddress != "" {
		stored, err := address.NewFromString(keyInfo.PqcAddress)
		if err != nil || stored != expected {
			return nil, xerrors.Errorf("keystore PQC address metadata does not match requested address %s", expected)
		}
	}
	keyInfo.PqcAddress = expected.String()
	return NewPqcKey(keyInfo)
}
