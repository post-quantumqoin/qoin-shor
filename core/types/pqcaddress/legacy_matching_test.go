package pqcaddress

import (
	"reflect"
	"testing"

	"github.com/post-quantumqoin/address"
	"github.com/post-quantumqoin/core-types/crypto"
)

func TestLegacyCertificateMatchesEveryKey(t *testing.T) {
	cert := testCertificate(SchemeLegacyCompatible,
		crypto.SignPqcCertPubkey{Typ: "dilithium3", Pubkey: []byte{1, 2, 3}},
		crypto.SignPqcCertPubkey{Typ: "falcon512", Pubkey: []byte{4, 5, 6}})
	before := append([]crypto.SignPqcCertPubkey(nil), cert.Pubkeys...)
	for _, protocol := range []address.Protocol{address.PQC, address.Delegated} {
		for _, key := range cert.Pubkeys {
			candidates, err := SignerAddressCandidates(testCertificate(SchemeLegacyCompatible, key), protocol)
			if err != nil {
				t.Fatal(err)
			}
			for _, candidate := range candidates {
				matched, err := LegacyCertificateMatchesAddress(cert, candidate)
				if err != nil || !matched {
					t.Fatalf("key %s candidate %s: matched=%v err=%v", key.Typ, candidate, matched, err)
				}
			}
		}
		unknown, err := fromBytes([]byte{99}, protocol)
		if err != nil {
			t.Fatal(err)
		}
		if matched, err := LegacyCertificateMatchesAddress(cert, unknown); err != nil || matched {
			t.Fatalf("unknown signer: matched=%v err=%v", matched, err)
		}
	}
	if !reflect.DeepEqual(cert.Pubkeys, before) {
		t.Fatal("certificate mutated")
	}
	cert.Pubkeys = append(cert.Pubkeys, cert.Pubkeys[0])
	candidates, _ := SignerAddressCandidates(testCertificate(SchemeLegacyCompatible, before[0]), address.PQC)
	if _, err := LegacyCertificateMatchesAddress(cert, candidates[0]); err == nil {
		t.Fatal("duplicate key accepted after first-key match")
	}
	cert.Pubkeys = before
	cert.Version = 1
	if _, err := LegacyCertificateMatchesAddress(cert, candidates[0]); err == nil {
		t.Fatal("version 1 accepted by legacy matcher")
	}
}
