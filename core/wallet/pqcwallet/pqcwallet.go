package pqcwallet

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	// "fmt"
	logging "github.com/ipfs/go-log/v2"
	"golang.org/x/xerrors"

	"github.com/post-quantumqoin/address"
	"github.com/post-quantumqoin/core-types/crypto"

	"github.com/post-quantumqoin/qoin-shor/api"
	"github.com/post-quantumqoin/qoin-shor/core/types"
	"github.com/post-quantumqoin/qoin-shor/core/wallet/key"

	// "github.com/filecoin-project/go-address"
	"github.com/post-quantumqoin/qoin-shor/lib/sigs"
	_ "github.com/post-quantumqoin/qoin-shor/lib/sigs/bls" // enable bls signatures
	_ "github.com/post-quantumqoin/qoin-shor/lib/sigs/delegated"
	_ "github.com/post-quantumqoin/qoin-shor/lib/sigs/pqc"  // enable pqc signatures
	_ "github.com/post-quantumqoin/qoin-shor/lib/sigs/secp" // enable secp signatures
	// "github.com/post-quantumqoin/qoin-shor/pqccrypto"
)

var log = logging.Logger("pqcwallet")

const (
	KNamePrefix  = "pqc-wallet-"
	KTrashPrefix = "pqc-trash-"
	KDefault     = "pqc-default"
	// KNamePrefix  = "wallet-"
	// KTrashPrefix = "trash-"
	// KDefault     = "default"

	KpPrefix    = "pqc-keypair"
	KFalcon512  = "falcon512"
	KFalcon1024 = "falcon1024"
	KDilithium3 = "dilithium3"
	KDilithium5 = "dilithium5"

	Secp256k1 = "secp256k1"
)

//type Address struct{ str string }

type PqcWallet struct {
	keys    map[address.Address]*key.Key
	Pqckeys map[address.Address]*key.PqcKey
	// Keypairs map[string]types.PqcKeypair
	keystore types.KeyStore

	lk sync.Mutex
}

type Default interface {
	GetDefault() (address.Address, error)
	SetDefault(a address.Address) error
}

func NewWallet(keystore types.KeyStore) (*PqcWallet, error) {
	w := &PqcWallet{
		keys:    make(map[address.Address]*key.Key),
		Pqckeys: make(map[address.Address]*key.PqcKey),
		// Keypairs: make(map[string]types.PqcKeypair),
		keystore: keystore,
	}

	return w, nil
}

func KeyWallet(keys ...*key.Key) *PqcWallet {
	m := make(map[address.Address]*key.Key)
	for _, key := range keys {
		m[key.Address] = key
	}

	return &PqcWallet{
		keys: m,
	}
}

func (w *PqcWallet) WalletSign(ctx context.Context, addr address.Address, msg []byte, meta api.MsgMeta) (*crypto.Signature, error) {
	pkey, err := w.findPqcKey(addr)
	if err != nil {
		return nil, err
	}

	for _, kp := range pkey.KeyInfo.PqcKeypairs {
		log.Debug("PqcSign kps  PqcType:", kp.PqcType)
		log.Debug("PqcSign kps  PqcVersion:", kp.PqcVersion)
		log.Debug("PqcSign kps  seed:", kp.PqcSeed)
	}

	if pkey.KeyInfo.Type == types.KTDelegated {
		sg, err := sigs.MultiPqcSign(pkey.PQCCert, pkey.KeyInfo.PqcKeypairs, msg)
		if err != nil {
			return nil, err
		}
		sg.Type = crypto.SigTypeDelegated
		return sg, nil
	}

	return sigs.MultiPqcSign(pkey.PQCCert, pkey.KeyInfo.PqcKeypairs, msg)
}

func (w *PqcWallet) findPqcKey(addr address.Address) (*key.PqcKey, error) {
	w.lk.Lock()
	defer w.lk.Unlock()
	fmt.Println("PqcWallet findPqcKey:", addr)
	pkey, ok := w.Pqckeys[addr]
	if ok {
		return pkey, nil
	}
	fmt.Println("w.keystore.Get addr:", KNamePrefix+addr.String())
	ki, err := w.keystore.Get(KNamePrefix + addr.String())
	if err != nil {
		if xerrors.Is(err, types.ErrKeyInfoNotFound) {
			return nil, nil
		}
		return nil, err
	}

	pkey, err = key.PqcNewKey(ki)
	if err != nil {
		return nil, xerrors.Errorf("failed to new pkey: %w", err)
	}
	log.Debug("findPqcKey PqcKey by addr Version", pkey.PQCCert.Version)
	w.Pqckeys[pkey.Address] = pkey

	return pkey, nil
}

func (w *PqcWallet) findKey(addr address.Address) (*key.Key, error) {
	w.lk.Lock()
	defer w.lk.Unlock()

	k, ok := w.keys[addr]
	if ok {
		return k, nil
	}
	if w.keystore == nil {
		log.Warn("findKey didn't find the key in in-memory wallet")
		return nil, nil
	}

	ki, err := w.tryFind(addr)
	if err != nil {
		if xerrors.Is(err, types.ErrKeyInfoNotFound) {
			return nil, nil
		}
		return nil, xerrors.Errorf("getting from keystore: %w", err)
	}
	k, err = key.NewKey(ki)
	if err != nil {
		return nil, xerrors.Errorf("decoding from keystore: %w", err)
	}
	w.keys[k.Address] = k
	return k, nil
}

func (w *PqcWallet) tryFind(addr address.Address) (types.KeyInfo, error) {

	ki, err := w.keystore.Get(KNamePrefix + addr.String())
	if err == nil {
		return ki, err
	}

	if !xerrors.Is(err, types.ErrKeyInfoNotFound) {
		return types.KeyInfo{}, err
	}

	// We got an ErrKeyInfoNotFound error
	// Try again, this time with the testnet prefix

	tAddress, err := swapMainnetForTestnetPrefix(addr.String())
	if err != nil {
		return types.KeyInfo{}, err
	}

	ki, err = w.keystore.Get(KNamePrefix + tAddress)
	if err != nil {
		return types.KeyInfo{}, err
	}

	// We found it with the testnet prefix
	// Add this KeyInfo with the mainnet prefix address string
	err = w.keystore.Put(KNamePrefix+addr.String(), ki)
	if err != nil {
		return types.KeyInfo{}, err
	}

	return ki, nil
}

func (w *PqcWallet) WalletExport(ctx context.Context, addr address.Address) (*types.KeyInfo, error) {
	pkey, err := w.findPqcKey(addr)
	if err != nil {
		return nil, xerrors.Errorf("failed to find cert to export: %w", err)
	}

	// pki := &types.KeyInfo{
	// 	PqcKeypairs: pkey.PqcKeypairs,
	// }

	return &pkey.KeyInfo, nil
}

func (w *PqcWallet) WalletImport(ctx context.Context, pki *types.KeyInfo) (address.Address, error) {
	w.lk.Lock()
	defer w.lk.Unlock()

	log.Debug("WalletImport pki typ:", pki.Type)
	for _, kp := range pki.PqcKeypairs {
		log.Debug("WalletImport Keypairs keyInfo PqcKeypair PqcType:", kp.PqcType)
		log.Debug("WalletImport Keypairs keyInfo PqcVersion:", kp.PqcVersion)
	}
	k, err := key.PqcNewKey(*pki)
	if err != nil {
		return address.Undef, xerrors.Errorf("failed to make key: %w", err)
	}

	if err := w.keystore.Put(KNamePrefix+k.Address.String(), k.KeyInfo); err != nil {
		return address.Undef, xerrors.Errorf("saving to keystore: %w", err)
	}
	log.Debug("WalletImport k.Address:", k.Address.String())
	return k.Address, nil
}

func (w *PqcWallet) WalletList(ctx context.Context) ([]address.Address, error) {
	all, err := w.keystore.List()
	if err != nil {
		return nil, xerrors.Errorf("listing keystore: %w", err)
	}

	sort.Strings(all)

	seen := map[address.Address]struct{}{}
	out := make([]address.Address, 0, len(all))
	for _, a := range all {
		if strings.HasPrefix(a, KNamePrefix) {
			name := strings.TrimPrefix(a, KNamePrefix)
			log.Debug("WalletList name:", name)
			// if err := w.keystore.Delete(KNamePrefix + name); err != nil {
			// 	return out, xerrors.Errorf("failed to delete key %s: %w", name, err)
			// }
			addr, err := address.NewFromString(name)
			log.Debug("WalletList addr:", addr.String())
			if err != nil {
				return nil, xerrors.Errorf("converting name to address: %w", err)
			}
			if _, ok := seen[addr]; ok {
				continue // got duplicate with a different prefix
			}
			seen[addr] = struct{}{}

			out = append(out, addr)
		}
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].String() < out[j].String()
	})

	return out, nil
}

func (w *PqcWallet) GetDefault() (address.Address, error) {
	w.lk.Lock()
	defer w.lk.Unlock()
	ki, err := w.keystore.Get(KDefault)
	if err != nil {
		return address.Undef, xerrors.Errorf("failed to get default key: %w", err)
	}
	if ki.Type == types.KTBLS || ki.Type == types.KTSecp256k1 {
		k, err := key.NewKey(ki)
		if err != nil {
			return address.Undef, xerrors.Errorf("failed to read default key from keystore: %w", err)
		}

		return k.Address, nil
	}

	if len(ki.PrivateKey) == 0 && len(ki.PqcKeypairs) == 0 {
		if err := w.keystore.Delete(KDefault); err != nil {
			if !xerrors.Is(err, types.ErrKeyInfoNotFound) {
				log.Warnf("failed to unregister current default key: %s", err)
			}
		}
		return address.Undef, xerrors.Errorf("Key info is nil")
	}

	PqcKey, err := key.PqcNewKey(ki)

	// log.Debug("GetDefault Address:", PqcKey.Address)
	// log.Debug("GetDefault PQCCert nonce:", PqcKey.PQCCert.Version)

	// for _, pubk := range PqcKey.PQCCert.Pubkeys {
	// 	log.Debugf("GetDefault pqccert pubk:%s\n", pubk.Typ)
	// }

	if err != nil {
		return address.Undef, xerrors.Errorf("failed to read default key from keystore: %w", err)
	}

	return PqcKey.Address, nil
}

func (w *PqcWallet) SetDefault(a address.Address) error {
	w.lk.Lock()
	defer w.lk.Unlock()
	log.Debug("SetDefault Address:", a.String())

	ki, err := w.keystore.Get(KNamePrefix + a.String())
	if err != nil {
		if !xerrors.Is(err, types.ErrKeyInfoNotFound) {
			return err
		}
	}

	if len(ki.PrivateKey) == 0 && len(ki.PqcKeypairs) == 0 {
		ki, err = w.keystore.Get("wallet-" + a.String())
		if err != nil {
			if !xerrors.Is(err, types.ErrKeyInfoNotFound) {
				return err
			}
		}
	}

	if len(ki.PrivateKey) == 0 && len(ki.PqcKeypairs) == 0 {
		return types.ErrKeyInfoNotFound
	}

	log.Debug("SetDefault len(ki.PqcKeypairs):", len(ki.PqcKeypairs))
	PqcKey, err := key.PqcNewKey(ki)
	log.Debug("SetDefault Address:", PqcKey.Address)

	log.Debug("SetDefault PQCCert Version:", PqcKey.PQCCert.Version)
	for _, pubk := range PqcKey.PQCCert.Pubkeys {
		log.Debugf("SetDefault pqccert pubk:%s\n", pubk.Typ)
	}

	if err := w.keystore.Delete(KDefault); err != nil {
		if !xerrors.Is(err, types.ErrKeyInfoNotFound) {
			log.Warnf("failed to unregister current default key: %s", err)
		}
	}

	if err := w.keystore.Put(KDefault, ki); err != nil {
		return err
	}

	return nil
}

func (w *PqcWallet) WalletNew(ctx context.Context, typ types.KeyType) (address.Address, error) {
	w.lk.Lock()
	defer w.lk.Unlock()

	log.Debug("PqcWallet WalletNew typ:", typ)
	k, err := key.PqcGenerateKey(typ)
	if err != nil {
		return address.Undef, err
	}

	log.Debug("k.PQCCert.Version:", k.PQCCert.Version)
	for _, pk := range k.PQCCert.Pubkeys {
		log.Debug("k.PQCCert.pk.Typ:", pk.Typ)
		log.Debug("k.PQCCert.pk.Pubkey:", len(pk.Pubkey))
	}

	for _, kps := range k.KeyInfo.PqcKeypairs {
		log.Debug("PqcKeypairs .PqcType:", kps.PqcType)
		log.Debug("PqcKeypairs.PqcVersion:", kps.PqcVersion)
		log.Debug("PqcKeypairs .PqcSeed:", kps.PqcSeed)
	}

	if err := w.keystore.Put(KNamePrefix+k.Address.String(), k.KeyInfo); err != nil {
		return address.Undef, xerrors.Errorf("saving to keystore: %w", err)
	}

	w.Pqckeys[k.Address] = k

	_, err = w.keystore.Get(KDefault)
	if err != nil {
		if !xerrors.Is(err, types.ErrKeyInfoNotFound) {
			return address.Undef, err
		}

		if err := w.keystore.Put(KDefault, k.KeyInfo); err != nil {
			return address.Undef, xerrors.Errorf("failed to set new key as default: %w", err)
		}
	}

	return k.Address, nil
}

func (w *PqcWallet) WalletHas(ctx context.Context, addr address.Address) (bool, error) {
	log.Debug("PqcWallet WalletHas addr:", addr)
	k, err := w.findPqcKey(addr)
	if err != nil {
		return false, err
	}
	log.Debug("PqcWallet WalletHas existence:", k != nil)
	return k != nil, nil
}

func (w *PqcWallet) walletDelete(ctx context.Context, addr address.Address) error {
	k, err := w.findPqcKey(addr)
	if err != nil {
		return xerrors.Errorf("failed to delete key %s : %w", addr, err)
	}
	if k == nil {
		return nil // already not there
	}
	log.Debug("walletDelete findPqcKey:", addr)
	w.lk.Lock()
	defer w.lk.Unlock()

	if err := w.keystore.Delete(KTrashPrefix + k.Address.String()); err != nil && !xerrors.Is(err, types.ErrKeyInfoNotFound) {
		return xerrors.Errorf("failed to purge trashed key %s: %w", addr, err)
	}

	if err := w.keystore.Put(KTrashPrefix+k.Address.String(), k.KeyInfo); err != nil {
		return xerrors.Errorf("failed to mark key %s as trashed: %w", addr, err)
	}

	if err := w.keystore.Delete(KNamePrefix + k.Address.String()); err != nil {
		return xerrors.Errorf("failed to delete key %s: %w", addr, err)
	}

	// tAddr, err := swapMainnetForTestnetPrefix(addr.String())
	// if err != nil {
	// 	return xerrors.Errorf("failed to swap prefixes: %w", err)
	// }

	// // TODO: Does this always error in the not-found case? Just ignoring an error return for now.
	// _ = w.keystore.Delete(KNamePrefix + tAddr)

	delete(w.Pqckeys, addr)

	return nil
}

func (w *PqcWallet) deleteDefault() {
	w.lk.Lock()
	defer w.lk.Unlock()
	if err := w.keystore.Delete(KDefault); err != nil {
		if !xerrors.Is(err, types.ErrKeyInfoNotFound) {
			log.Warnf("failed to unregister current default key: %s", err)
		}
	}
}

func (w *PqcWallet) WalletDelete(ctx context.Context, addr address.Address) error {
	if err := w.walletDelete(ctx, addr); err != nil {
		return xerrors.Errorf("wallet delete: %w", err)
	}

	if def, err := w.GetDefault(); err == nil {
		if def == addr {
			w.deleteDefault()
		}
	}
	return nil
}

func (w *PqcWallet) Get() api.Wallet {
	if w == nil {
		return nil
	}

	return w
}

var _ api.Wallet = &PqcWallet{}

func swapMainnetForTestnetPrefix(addr string) (string, error) {
	aChars := []rune(addr)
	prefixRunes := []rune(address.TestnetPrefix)
	if len(prefixRunes) != 1 {
		return "", xerrors.Errorf("unexpected prefix length: %d", len(prefixRunes))
	}

	aChars[0] = prefixRunes[0]
	return string(aChars), nil
}

type nilDefault struct{}

func (n nilDefault) GetDefault() (address.Address, error) {
	return address.Undef, nil
}

func (n nilDefault) SetDefault(a address.Address) error {
	return xerrors.Errorf("not supported; local wallet disabled")
}

var NilDefault nilDefault
var _ Default = NilDefault
