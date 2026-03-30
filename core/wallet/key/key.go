package key

import (
	"fmt"
	"sort"
	// "strings"

	"golang.org/x/xerrors"

	logging "github.com/ipfs/go-log/v2"
	// "github.com/post-quantumqoin/qoin-shor/lib/address"
	"github.com/post-quantumqoin/address"
	"github.com/post-quantumqoin/core-types/crypto"

	"github.com/post-quantumqoin/qoin-shor/core/types"
	"github.com/post-quantumqoin/qoin-shor/core/types/ethtypes"
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
	
	fmt.Println("GeneratePqcKeyWithAlgs generating key for keyType:", typ, "algs:", list)
	var kprs []types.PqcKeypair
    for _, a := range list {
        ctyp := PqcActSigType(a)
		fmt.Println("GeneratePqcKeyWithAlgs generating key for alg:", a, "ctyp:", ctyp)
        seed, sk, pk, err := sigs.PqcGenerate(ctyp)
        if err != nil {
            return nil, xerrors.Errorf("generate %s failed: %w", a, err)
        }
        kprs = append(kprs, types.PqcKeypair{
            PqcVersion:     0,
            PqcSeed:        seed,
            PqcType:        types.KeyType(a),
            PqcPrivateKey:  sk,
            PqcPublicKey:   pk,
        })
    }

    ki := types.KeyInfo{
        Type:        typ, // This should be set to a type that indicates it's a multi-algorithm key, e.g., "pqc-multi"
        PqcKeypairs: kprs,
    }

    return NewPqcKey(ki)
	// if strings.Contains(string(sg), string(types.KTDelegated)) {
	// 	sg = "falcon512 dilithium3"
	// 	ki.Type = typ
	// }

	// if strings.Contains(string(sg), string(types.KTPqc)) {
	// 	typ = "falcon512 dilithium3"
	// 	ki.Type = typ
	// }

	// if !strings.Contains(string(sg), string(types.Falcon512)) &&
	// 	!strings.Contains(string(sg), string(types.Falcon1024)) &&
	// 	!strings.Contains(string(sg), string(types.Dilithium3)) &&
	// 	!strings.Contains(string(sg), string(types.Dilithium5)) {
	// 	return nil, xerrors.Errorf("unknown sig type: %s", typ)
	// }
	// fd := strings.Fields(string(sg))
	// var kprs []types.PqcKeypair
	// for _, tp := range fd {
	// 	ctyp := PqcActSigType(types.KeyType(tp))
	// 	if ctyp == crypto.SigTypeUnknown {
	// 		return nil, xerrors.Errorf("unknown sig type: %s", typ)
	// 	}
	// 	// fmt.Println("key GeneratePqcKey ctyp:",ctyp)
	// 	seed, sk, pk, err := sigs.PqcGenerate(ctyp)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	kpr := types.PqcKeypair{
	// 		PqcVersion: 0,
	// 		PqcSeed:    seed,
	// 		PqcType:    types.KeyType(tp),

	// 		PqcPrivateKey: sk,
	// 		PqcPublicKey:  pk,
	// 	}
	// 	kprs = append(kprs, kpr)
	// }
	// // ki := types.KeyInfo{
	// // 	PqcKeypairs: kprs,
	// // }
	// ki.PqcKeypairs = kprs
	// return NewPqcKey(ki)
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
	fmt.Println("NewPqcKey KeyInfo Type:", keyInfo.Type, "NumKeypairs:", len(keyInfo.PqcKeypairs))
	// Sort keypairs by type to ensure deterministic order in cert
	kpairs := make([]types.PqcKeypair, len(keyInfo.PqcKeypairs))
    copy(kpairs, keyInfo.PqcKeypairs)
    sort.Slice(kpairs, func(i, j int) bool {
        return string(kpairs[i].PqcType) < string(kpairs[j].PqcType)
    })

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
        Version: 0,
    }

    pk := &PqcKey{
        KeyInfo: keyInfo,
        PQCCert: cert,
    }

	// ensure at least one pubkey exists
    if len(pubkeys) == 0 {
        return nil, xerrors.Errorf("no cert pubkeys built")
    }

	// If the key type is Delegated, we derive the address from the first public key using Ethereum's method
    if keyInfo.Type == types.KTDelegated {
        pb, err := pubkeys[0].Serialize()
        if err != nil {
            return nil, xerrors.Errorf("serialize delegated pubkey: %w", err)
        }
        ethAddr, err := ethtypes.EthAddressFromPubKey(pb)
        if err != nil {
            return nil, xerrors.Errorf("eth address from pubkey: %w", err)
        }
        ea, err := ethtypes.CastEthAddress(ethAddr)
        if err != nil {
            return nil, xerrors.Errorf("cast eth address: %w", err)
        }
        pk.Address, err = ea.ToFilecoinAddress()
        if err != nil {
            return nil, xerrors.Errorf("convert delegated to filecoin address: %w", err)
        }
        return pk, nil
    }

	// For non-delegated keys, we derive the address from the first public key using Filecoin's method
    pb, err := pubkeys[0].Serialize()
    if err != nil {
        return nil, xerrors.Errorf("serialize pqc pubkey: %w", err)
    }
    pk.Address, err = address.NewPqcAddress(pb)
    if err != nil {
        return nil, xerrors.Errorf("new pqc address: %w", err)
    }

	return pk, nil
	// var Pbks []types.PqcCertPubkey

	// for _, kp := range KeyInfo.PqcKeypairs {
	// 	Pbks = append(Pbks, types.PqcCertPubkey{Typ: string(kp.PqcType), Pubkey: kp.PqcPublicKey})
	// }
	// fmt.Println("NewPqcKey Pbks len:", len(Pbks))
	// cert := types.PQCCert{
	// 	Pubkeys: Pbks,
	// 	Version: 0,
	// }

	// k := &PqcKey{
	// 	KeyInfo: KeyInfo,
	// 	PQCCert: cert,
	// }
	// if KeyInfo.Type == types.KTDelegated {
	// 	//Use the first public key as the address
	// 	pb, err := Pbks[0].Serialize()
	// 	if err != nil {
	// 		return nil, fmt.Errorf("pqc public key serialization failure")
	// 	}
	// 	// Transitory Delegated signature verification as per FIP-0055
	// 	ethAddr, err := ethtypes.EthAddressFromPubKey(pb)
	// 	if err != nil {
	// 		return nil, xerrors.Errorf("failed to calculate Eth address from public key: %w", err)
	// 	}

	// 	ea, err := ethtypes.CastEthAddress(ethAddr)
	// 	if err != nil {
	// 		return nil, xerrors.Errorf("failed to create ethereum address from bytes: %w", err)
	// 	}

	// 	k.Address, err = ea.ToFilecoinAddress()
	// 	if err != nil {
	// 		return nil, xerrors.Errorf("converting Delegated to address: %w", err)
	// 	}
	// 	return k, nil
	// }

	//Use the first public key as the address
	// pb, err := Pbks[0].Serialize()
	// if err != nil {
	// 	return nil, fmt.Errorf("pqc public key serialization failure")
	// }

	// k.Address, err = address.NewPqcAddress(pb)
	// return k, nil
}
