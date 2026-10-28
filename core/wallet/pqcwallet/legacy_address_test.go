package pqcwallet_test

import (
	"context"
	"testing"

	"github.com/post-quantumqoin/address"
	"github.com/post-quantumqoin/core-types/crypto"
	"github.com/post-quantumqoin/qoin-shor/api"
	"github.com/post-quantumqoin/qoin-shor/core/types"
	"github.com/post-quantumqoin/qoin-shor/core/types/pqcaddress"
	"github.com/post-quantumqoin/qoin-shor/core/wallet"
	"github.com/post-quantumqoin/qoin-shor/core/wallet/key"
	"github.com/post-quantumqoin/qoin-shor/core/wallet/pqcwallet"
	"github.com/post-quantumqoin/qoin-shor/lib/sigs"
)

func TestLegacyDefaultExportReloadAndSigning(t *testing.T) {
	generated, err := key.GeneratePqcKeyWithAlgs(types.KTPqc, []types.SigAlg{types.Falcon512, types.Dilithium3})
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	for _, protocol := range []address.Protocol{address.PQC, address.Delegated} {
		info := generated.KeyInfo
		info.PqcAddressVersion = nil
		info.PqcAddress = ""
		if protocol == address.Delegated {
			info.Type = types.KTDelegated
		}
		cert := crypto.SignPQCCert{Version: 0}
		for _, kp := range info.PqcKeypairs {
			cert.Pubkeys = append(cert.Pubkeys, crypto.SignPqcCertPubkey{Typ: string(kp.PqcType), Pubkey: kp.PqcPublicKey})
		}
		candidates, err := pqcaddress.LegacyWalletAddressCandidates(cert, protocol)
		if err != nil {
			t.Fatal(err)
		}
		for _, candidate := range candidates {
			store := wallet.NewMemKeyStore()
			if err := store.Put(pqcwallet.KNamePrefix+candidate.Address.String(), info); err != nil {
				t.Fatal(err)
			}
			w, err := pqcwallet.NewWallet(store)
			if err != nil {
				t.Fatal(err)
			}
			if err := w.SetDefault(candidate.Address); err != nil {
				t.Fatal(err)
			}
			persisted, err := store.Get(pqcwallet.KDefault)
			if err != nil || persisted.PqcAddress != candidate.Address.String() || persisted.PqcAddressVersion == nil || *persisted.PqcAddressVersion != 0 {
				t.Fatalf("default metadata not persisted: %v", err)
			}
			// Rebuild the wallet to discard its in-memory cache.
			w, err = pqcwallet.NewWallet(store)
			if err != nil {
				t.Fatal(err)
			}
			defaultAddress, err := w.GetDefault()
			if err != nil || defaultAddress != candidate.Address {
				t.Fatalf("default changed on restart: %v", err)
			}
			exported, err := w.WalletExport(ctx, candidate.Address)
			if err != nil || exported.PqcAddress != candidate.Address.String() {
				t.Fatalf("export identity changed: %v", err)
			}
			imported, err := pqcwallet.NewWallet(wallet.NewMemKeyStore())
			if err != nil {
				t.Fatal(err)
			}
			importAddress, err := imported.WalletImport(ctx, exported)
			if err != nil || importAddress != candidate.Address {
				t.Fatalf("import identity changed: %v", err)
			}
			message := []byte("legacy wallet compatibility regression")
			signature, err := imported.WalletSign(ctx, importAddress, message, api.MsgMeta{Type: api.MTUnknown})
			if err != nil {
				t.Fatal(err)
			}
			if signature.PqcCert.Version != 0 || signature.PqcCert.Pubkeys[0].Typ != "falcon512" {
				t.Fatal("signing reordered the legacy certificate")
			}
			if err := sigs.MultiPqcVerify(signature, importAddress, message); err != nil {
				t.Fatal(err)
			}
			for index := range signature.PqcSignatures {
				original := signature.PqcSignatures[index].Data
				signature.PqcSignatures[index].Data = append([]byte(nil), original...)
				signature.PqcSignatures[index].Data[100] ^= 1
				if err := sigs.MultiPqcVerify(signature, importAddress, message); err == nil {
					t.Fatalf("invalid required signature %d accepted", index)
				}
				signature.PqcSignatures[index].Data = original
			}
			wrong, err := address.NewPqcAddress([]byte("not a historical key"))
			if err != nil {
				t.Fatal(err)
			}
			if err := store.Put(pqcwallet.KNamePrefix+wrong.String(), info); err != nil {
				t.Fatal(err)
			}
			if err := w.SetDefault(wrong); err == nil {
				t.Fatal("unrelated default address accepted")
			}
			unchanged, err := w.GetDefault()
			if err != nil || unchanged != candidate.Address {
				t.Fatal("failed validation changed default")
			}
		}
	}
}
