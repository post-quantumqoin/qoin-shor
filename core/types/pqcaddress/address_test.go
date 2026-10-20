package pqcaddress

import (
	"encoding/hex"
	"testing"

	"github.com/post-quantumqoin/address"
	typescrypto "github.com/post-quantumqoin/core-types/crypto"
)

func testCertificate(version Scheme, keys ...typescrypto.SignPqcCertPubkey) typescrypto.SignPQCCert {
	return typescrypto.SignPQCCert{Pubkeys: keys, Version: uint8(version)}
}

func TestSignerAddressCandidatesLegacyUsesOnlyFirstKey(t *testing.T) {
	cert := testCertificate(
		SchemeLegacyCompatible,
		typescrypto.SignPqcCertPubkey{Typ: "falcon512", Pubkey: []byte{1, 2, 3}},
		typescrypto.SignPqcCertPubkey{Typ: "dilithium3", Pubkey: []byte{4, 5, 6}},
	)

	candidates, err := SignerAddressCandidates(cert, address.Delegated)
	if err != nil {
		t.Fatal(err)
	}
	if len(candidates) != 2 {
		t.Fatalf("expected CBOR and raw legacy candidates, got %d", len(candidates))
	}

	serialized, err := serializeCertPublicKey(cert.Pubkeys[0], SchemeLegacyCompatible)
	if err != nil {
		t.Fatal(err)
	}
	wantCBOR, err := fromBytes(serialized, address.Delegated)
	if err != nil {
		t.Fatal(err)
	}
	wantRaw, err := fromBytes(cert.Pubkeys[0].Pubkey, address.Delegated)
	if err != nil {
		t.Fatal(err)
	}
	secondRaw, err := fromBytes(cert.Pubkeys[1].Pubkey, address.Delegated)
	if err != nil {
		t.Fatal(err)
	}

	if candidates[0] != wantCBOR || candidates[1] != wantRaw {
		t.Fatal("legacy candidates do not use first-key CBOR then first-key raw order")
	}
	if Contains(secondRaw, candidates) {
		t.Fatal("address derived from a non-first certificate key must not be accepted")
	}
}

func TestSerializeCertPublicKeyMatchesLegacyShorWireFormat(t *testing.T) {
	serialized, err := serializeCertPublicKey(typescrypto.SignPqcCertPubkey{
		Typ:    "falcon512",
		Pubkey: []byte{1, 2, 3},
	}, SchemeLegacyCompatible)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := hex.EncodeToString(serialized), "816966616c636f6e35313243010203"; got != want {
		t.Fatalf("unexpected certificate public-key bytes: got %s, want %s", got, want)
	}
}

func TestSignerAddressCandidatesCanonicalScheme(t *testing.T) {
	keys := []typescrypto.SignPqcCertPubkey{
		{Typ: "dilithium3", Pubkey: []byte{4, 5, 6}},
		{Typ: "falcon512", Pubkey: []byte{1, 2, 3}},
	}

	cborCandidates, err := SignerAddressCandidates(testCertificate(SchemeCanonicalCBOR, keys...), address.PQC)
	if err != nil {
		t.Fatal(err)
	}
	if len(cborCandidates) != 1 {
		t.Fatal("the modern address scheme must produce exactly one CBOR candidate")
	}
}

func TestVersionedAddressEncoding(t *testing.T) {
	key := typescrypto.SignPqcCertPubkey{Typ: "falcon512", Pubkey: []byte{1, 2, 3}}
	for _, test := range []struct {
		scheme Scheme
		want   string
	}{
		{SchemeLegacyCompatible, "816966616c636f6e35313243010203"},
		{SchemeCanonicalCBOR, "826966616c636f6e35313243010203"},
	} {
		encoded, err := serializeCertPublicKey(key, test.scheme)
		if err != nil {
			t.Fatal(err)
		}
		if got := hex.EncodeToString(encoded); got != test.want {
			t.Fatalf("version %d: got %s, want %s", test.scheme, got, test.want)
		}
	}
	if _, err := serializeCertPublicKey(key, Scheme(2)); err == nil {
		t.Fatal("unsupported version accepted")
	}
	for _, protocol := range []address.Protocol{address.PQC, address.Delegated} {
		legacy, err := SignerAddressCandidates(testCertificate(SchemeLegacyCompatible, key), protocol)
		if err != nil {
			t.Fatal(err)
		}
		modern, err := SignerAddressCandidates(testCertificate(SchemeCanonicalCBOR, key), protocol)
		if err != nil {
			t.Fatal(err)
		}
		encoded, _ := hex.DecodeString("826966616c636f6e35313243010203")
		want, err := fromBytes(encoded, protocol)
		if err != nil {
			t.Fatal(err)
		}
		if len(modern) != 1 || modern[0] != want {
			t.Fatal("version 1 must use standard CBOR")
		}
		if Contains(legacy[0], modern) || Contains(legacy[1], modern) || Contains(want, legacy) {
			t.Fatal("versioned addresses must not accept the other encoding")
		}
	}
}

func TestLegacyPreservesCertificateOrder(t *testing.T) {
	first := typescrypto.SignPqcCertPubkey{Typ: "falcon512", Pubkey: []byte{1, 2, 3}}
	second := typescrypto.SignPqcCertPubkey{Typ: "dilithium3", Pubkey: []byte{4, 5, 6}}
	before, err := SignerAddressCandidates(testCertificate(SchemeLegacyCompatible, first, second), address.PQC)
	if err != nil {
		t.Fatal(err)
	}
	after, err := SignerAddressCandidates(testCertificate(SchemeLegacyCompatible, second, first), address.PQC)
	if err != nil {
		t.Fatal(err)
	}
	if before[0] == after[0] {
		t.Fatal("legacy first key was reordered")
	}
	if _, err := SignerAddressCandidates(testCertificate(SchemeCanonicalCBOR, first, second), address.PQC); err == nil {
		t.Fatal("version 1 must reject noncanonical order")
	}
}

func TestValidateCertificateRejectsInvalidModernCertificates(t *testing.T) {
	tests := []typescrypto.SignPQCCert{
		testCertificate(
			SchemeCanonicalCBOR,
			typescrypto.SignPqcCertPubkey{Typ: "falcon512", Pubkey: []byte{1}},
			typescrypto.SignPqcCertPubkey{Typ: "dilithium3", Pubkey: []byte{2}},
		),
		testCertificate(
			SchemeLegacyCompatible,
			typescrypto.SignPqcCertPubkey{Typ: "falcon512", Pubkey: []byte{1}},
			typescrypto.SignPqcCertPubkey{Typ: "falcon512", Pubkey: []byte{2}},
		),
		testCertificate(
			SchemeCanonicalCBOR,
			typescrypto.SignPqcCertPubkey{Typ: "unknown", Pubkey: []byte{1}},
		),
		testCertificate(
			SchemeCanonicalCBOR,
			typescrypto.SignPqcCertPubkey{Typ: "falcon512"},
		),
		{
			Version: 2,
			Pubkeys: []typescrypto.SignPqcCertPubkey{{Typ: "falcon512", Pubkey: []byte{1}}},
		},
	}

	for index, cert := range tests {
		if err := ValidateCertificate(cert); err == nil {
			t.Fatalf("test %d: expected certificate validation to fail", index)
		}
	}
}
