package key

import (
	"fmt"
	"strings"

	"golang.org/x/xerrors"

	// "github.com/post-quantumqoin/qoin-shor/lib/address"
	"github.com/post-quantumqoin/address"
	"github.com/post-quantumqoin/core-types/crypto"

	"github.com/post-quantumqoin/qoin-shor/core/types"
	"github.com/post-quantumqoin/qoin-shor/core/types/ethtypes"
	"github.com/post-quantumqoin/qoin-shor/lib/sigs"
	// pqccrypto "github.com/post-quantumqoin/qoin-shor/pqccrypto"
)

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

		k.Address, err = ea.ToFilecoinAddress()
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
	case types.KTDelegated:
		return crypto.SigTypeDelegated
	default:
		return crypto.SigTypeUnknown
	}
}

// //type SigType byte
func PqcActSigType(typ types.KeyType) crypto.SigType {
	switch typ {
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

func PqcGenerateKey(typ types.KeyType) (*PqcKey, error) {
	ki := types.KeyInfo{
		Type: types.KTPqc,
	}

	if strings.Contains(string(typ), string(types.KTDelegated)) {
		typ = "falcon512 dilithium3"
		ki.Type = types.KTDelegated
	}

	if strings.Contains(string(typ), string(types.KTPqc)) {
		typ = "falcon512 dilithium3"
		ki.Type = types.KTPqc
	}

	if !strings.Contains(string(typ), string(types.Falcon512)) &&
		!strings.Contains(string(typ), string(types.Falcon1024)) &&
		!strings.Contains(string(typ), string(types.Dilithium3)) &&
		!strings.Contains(string(typ), string(types.Dilithium5)) {
		return nil, xerrors.Errorf("unknown sig type: %s", typ)
	}
	fd := strings.Fields(string(typ))
	var kprs []types.PqcKeypair
	for _, tp := range fd {
		ctyp := PqcActSigType(types.KeyType(tp))
		if ctyp == crypto.SigTypeUnknown {
			return nil, xerrors.Errorf("unknown sig type: %s", typ)
		}
		// fmt.Println("key PqcGenerateKey ctyp:",ctyp)
		seed, sk, pk, err := sigs.PqcGenerate(ctyp)
		if err != nil {
			return nil, err
		}
		kpr := types.PqcKeypair{
			PqcVersion: 0,
			PqcSeed:    seed,
			PqcType:    types.KeyType(tp),

			PqcPrivateKey: sk,
			PqcPublicKey:  pk,
		}
		kprs = append(kprs, kpr)
	}
	// ki := types.KeyInfo{
	// 	PqcKeypairs: kprs,
	// }
	ki.PqcKeypairs = kprs
	return PqcNewKey(ki)
}

type PqcKey struct {
	types.KeyInfo
	types.PQCCert

	// Hcert   []byte
	Address address.Address
}

func PqcNewKey(KeyInfo types.KeyInfo) (*PqcKey, error) {
	var Pbks []types.PqcCertPubkey

	for _, kp := range KeyInfo.PqcKeypairs {
		Pbks = append(Pbks, types.PqcCertPubkey{Typ: string(kp.PqcType), Pubkey: kp.PqcPublicKey})
	}
	fmt.Println("PqcNewKey Pbks len:", len(Pbks))
	cert := types.PQCCert{
		Pubkeys: Pbks,
		Version: 0,
	}

	k := &PqcKey{
		KeyInfo: KeyInfo,
		PQCCert: cert,
	}
	if KeyInfo.Type == types.KTDelegated {
		//Use the first public key as the address
		pb, err := Pbks[0].Serialize()
		if err != nil {
			return nil, fmt.Errorf("pqc public key serialization failure")
		}
		// Transitory Delegated signature verification as per FIP-0055
		ethAddr, err := ethtypes.EthAddressFromPubKey(pb)
		if err != nil {
			return nil, xerrors.Errorf("failed to calculate Eth address from public key: %w", err)
		}

		ea, err := ethtypes.CastEthAddress(ethAddr)
		if err != nil {
			return nil, xerrors.Errorf("failed to create ethereum address from bytes: %w", err)
		}

		k.Address, err = ea.ToFilecoinAddress()
		if err != nil {
			return nil, xerrors.Errorf("converting Delegated to address: %w", err)
		}
		return k, nil
	}

	//Use the first public key as the address
	pb, err := Pbks[0].Serialize()
	if err != nil {
		return nil, fmt.Errorf("pqc public key serialization failure")
	}

	k.Address, err = address.NewPqcAddress(pb)
	return k, nil
}
