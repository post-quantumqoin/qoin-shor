package pqcaddress

import (
	"reflect"
	"testing"

	"github.com/post-quantumqoin/address"
	"github.com/post-quantumqoin/core-types/crypto"
)

func TestLegacyWalletHistoricalSelections(t *testing.T) {
	cert := crypto.SignPQCCert{Version: 0, Pubkeys: []crypto.SignPqcCertPubkey{
		{Typ: "falcon512", Pubkey: []byte{1, 2, 3}},
		{Typ: "dilithium3", Pubkey: []byte{4, 5, 6}},
		{Typ: "falcon1024", Pubkey: []byte{7, 8, 9}},
	}}
	before := append([]crypto.SignPqcCertPubkey(nil), cert.Pubkeys...)
	for _, protocol := range []address.Protocol{address.PQC, address.Delegated} {
		candidates, err := LegacyWalletAddressCandidates(cert, protocol)
		if err != nil {
			t.Fatal(err)
		}
		if len(candidates) != 4 {
			t.Fatalf("expected four historical candidates, got %d", len(candidates))
		}
		for index, key := range cert.Pubkeys[:2] {
			want, err := SignerAddressCandidates(crypto.SignPQCCert{Version: 0, Pubkeys: []crypto.SignPqcCertPubkey{key}}, protocol)
			if err != nil {
				t.Fatal(err)
			}
			selection := "original_first"
			if index == 1 {
				selection = "sorted_first"
			}
			for i, addr := range want {
				got := candidates[index*2+i]
				encoding := "legacy_cbor"
				if i == 1 {
					encoding = "raw"
				}
				if got.Address != addr || got.Selection != selection || got.Encoding != encoding {
					t.Fatalf("incorrect candidate: %+v", got)
				}
			}
		}
		other, err := SignerAddressCandidates(crypto.SignPQCCert{Version: 0, Pubkeys: cert.Pubkeys[2:]}, protocol)
		if err != nil {
			t.Fatal(err)
		}
		for _, candidate := range candidates {
			if Contains(candidate.Address, other) {
				t.Fatal("arbitrary third key accepted as wallet history")
			}
		}
		single := crypto.SignPQCCert{Version: 0, Pubkeys: cert.Pubkeys[:1]}
		dedup, err := LegacyWalletAddressCandidates(single, protocol)
		if err != nil || len(dedup) != 2 {
			t.Fatalf("duplicate selections not removed: %+v %v", dedup, err)
		}
	}
	if !reflect.DeepEqual(before, cert.Pubkeys) {
		t.Fatal("input certificate reordered")
	}
	for _, invalid := range []crypto.SignPQCCert{
		{Version: 1, Pubkeys: cert.Pubkeys[:1]},
		{Version: 0},
		{Version: 0, Pubkeys: []crypto.SignPqcCertPubkey{cert.Pubkeys[0], cert.Pubkeys[0]}},
	} {
		if _, err := LegacyWalletAddressCandidates(invalid, address.PQC); err == nil {
			t.Fatal("invalid certificate accepted")
		}
	}
}
