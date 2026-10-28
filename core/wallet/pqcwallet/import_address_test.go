package pqcwallet_test

import (
	"context"
	"strings"
	"testing"

	"github.com/post-quantumqoin/address"
	"github.com/post-quantumqoin/core-types/crypto"
	"github.com/post-quantumqoin/qoin-shor/core/types"
	"github.com/post-quantumqoin/qoin-shor/core/types/pqcaddress"
	"github.com/post-quantumqoin/qoin-shor/core/wallet"
	"github.com/post-quantumqoin/qoin-shor/core/wallet/key"
	"github.com/post-quantumqoin/qoin-shor/core/wallet/pqcwallet"
)

func TestImportRequiresLegacyAddressAndPreservesSelection(t *testing.T) {
	generated, err := key.GeneratePqcKeyWithAlgs(types.KTPqc, []types.SigAlg{types.Falcon512, types.Dilithium3})
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	for _, protocol := range []address.Protocol{address.PQC, address.Delegated} {
		info := generated.KeyInfo
		info.PqcAddressVersion, info.PqcAddress = nil, ""
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
		store := wallet.NewMemKeyStore()
		w, err := pqcwallet.NewWallet(store)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.WalletImport(ctx, &info); err == nil || !strings.Contains(err.Error(), "ambiguous legacy PQC import") {
			t.Fatalf("expected ambiguous import rejection: %v", err)
		}
		if names, err := store.List(); err != nil || len(names) != 0 {
			t.Fatal("rejected import wrote keystore")
		}
		if info.PqcAddress != "" || info.PqcAddressVersion != nil {
			t.Fatal("input mutated")
		}
		for _, candidate := range candidates {
			selected := info
			selected.PqcAddress = candidate.Address.String()
			got, err := w.WalletImport(ctx, &selected)
			if err != nil || got != candidate.Address {
				t.Fatalf("import drift: %v", err)
			}
			ki, err := store.Get(pqcwallet.KNamePrefix + got.String())
			if err != nil || ki.PqcAddress != got.String() || ki.PqcAddressVersion == nil || *ki.PqcAddressVersion != 0 {
				t.Fatalf("missing persisted metadata: %v", err)
			}
			reloaded, err := key.NewPqcKey(ki)
			if err != nil || reloaded.Address != got {
				t.Fatalf("reload drift: %v", err)
			}
			if reloaded.PQCCert.Pubkeys[0].Typ != string(info.PqcKeypairs[0].PqcType) {
				t.Fatal("certificate reordered")
			}
		}
		before, _ := store.List()
		bad, _ := address.NewPqcAddress([]byte("unrelated import address"))
		info.PqcAddress = bad.String()
		if _, err := w.WalletImport(ctx, &info); err == nil {
			t.Fatal("unrelated address accepted")
		}
		after, _ := store.List()
		if len(before) != len(after) {
			t.Fatal("failed import wrote keystore")
		}
	}
	store := wallet.NewMemKeyStore()
	w, err := pqcwallet.NewWallet(store)
	if err != nil {
		t.Fatal(err)
	}
	modern := generated.KeyInfo
	modern.PqcAddress = ""
	got, err := w.WalletImport(ctx, &modern)
	if err != nil || got != generated.Address {
		t.Fatalf("V1 import changed: %v", err)
	}
	if _, err := w.WalletImport(ctx, nil); err == nil {
		t.Fatal("nil import accepted")
	}
	// A single legacy key is still ambiguous between CBOR and raw encoding.
	legacy := generated.KeyInfo
	legacy.PqcAddressVersion, legacy.PqcAddress = nil, ""
	legacy.PqcKeypairs = legacy.PqcKeypairs[:1]
	if _, err := w.WalletImport(ctx, &legacy); err == nil {
		t.Fatal("single-key legacy ambiguity ignored")
	}
}
