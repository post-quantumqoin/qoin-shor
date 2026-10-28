package sigs

import (
	"testing"

	"github.com/post-quantumqoin/address"
	"github.com/post-quantumqoin/core-types/crypto"

	"github.com/post-quantumqoin/qoin-shor/core/types/pqcaddress"
)

type acceptingPqcShim struct{}

func (acceptingPqcShim) PqcGenPrivate() ([]byte, []byte, []byte, error) {
	return nil, nil, nil, nil
}

func (acceptingPqcShim) PqcToPublic([]byte) ([]byte, error) { return nil, nil }
func (acceptingPqcShim) PqcSign([]byte, []byte) ([]byte, error) {
	return []byte{1}, nil
}
func (acceptingPqcShim) PqcVerify([]byte, []byte, []byte) error { return nil }

func installAcceptingPqcShims(t *testing.T) {
	t.Helper()
	previous := Pqcsigs
	Pqcsigs = map[crypto.SigType]PqcSigShim{
		crypto.SigTypeDilithium3: acceptingPqcShim{},
		crypto.SigTypeFalcon512:  acceptingPqcShim{},
	}
	t.Cleanup(func() { Pqcsigs = previous })
}

func validCanonicalSignature() *crypto.Signature {
	return &crypto.Signature{
		Type: crypto.SigTypeDelegated,
		PqcCert: crypto.SignPQCCert{
			Version: uint8(pqcaddress.SchemeCanonicalCBOR),
			Pubkeys: []crypto.SignPqcCertPubkey{
				{Typ: "dilithium3", Pubkey: []byte{1, 2, 3}},
				{Typ: "falcon512", Pubkey: []byte{4, 5, 6}},
			},
		},
		PqcSignatures: []crypto.PqcSignature{
			{Type: crypto.SigTypeDilithium3, Data: []byte{7}},
			{Type: crypto.SigTypeFalcon512, Data: []byte{8}},
		},
	}
}

func TestMultiPqcVerifyBindsCanonicalAddress(t *testing.T) {
	installAcceptingPqcShims(t)
	signature := validCanonicalSignature()
	candidates, err := pqcaddress.SignerAddressCandidates(signature.PqcCert, address.Delegated)
	if err != nil {
		t.Fatal(err)
	}

	if err := MultiPqcVerify(signature, candidates[0], []byte("message")); err != nil {
		t.Fatalf("verify canonical signer: %v", err)
	}

	secondKey := crypto.SignPQCCert{
		Version: uint8(pqcaddress.SchemeLegacyCompatible),
		Pubkeys: []crypto.SignPqcCertPubkey{
			{Typ: "falcon512", Pubkey: signature.PqcCert.Pubkeys[1].Pubkey},
		},
	}
	secondCandidates, err := pqcaddress.SignerAddressCandidates(secondKey, address.Delegated)
	if err != nil {
		t.Fatal(err)
	}
	if len(secondCandidates) != 2 {
		t.Fatalf("expected two legacy candidates, got %d", len(secondCandidates))
	}
	if err := MultiPqcVerify(signature, secondCandidates[1], []byte("message")); err == nil {
		t.Fatal("address derived from a non-first certificate key was accepted")
	}
}

func TestMultiPqcVerifyLegacySupportsFirstKeyCborAndRaw(t *testing.T) {
	installAcceptingPqcShims(t)
	signature := validCanonicalSignature()
	signature.PqcCert.Version = uint8(pqcaddress.SchemeLegacyCompatible)
	candidates, err := pqcaddress.SignerAddressCandidates(signature.PqcCert, address.Delegated)
	if err != nil {
		t.Fatal(err)
	}
	if len(candidates) != 2 {
		t.Fatalf("expected two legacy candidates, got %d", len(candidates))
	}
	for _, candidate := range candidates {
		if err := MultiPqcVerify(signature, candidate, []byte("message")); err != nil {
			t.Fatalf("verify legacy candidate %s: %v", candidate, err)
		}
	}
}

func TestMultiPqcVerifyRejectsInvalidEnvelope(t *testing.T) {
	installAcceptingPqcShims(t)
	tests := []struct {
		name   string
		mutate func(*crypto.Signature)
	}{
		{
			name: "legacy data mixed with PQC",
			mutate: func(signature *crypto.Signature) {
				signature.Data = []byte{1}
			},
		},
		{
			name: "missing signature",
			mutate: func(signature *crypto.Signature) {
				signature.PqcSignatures = signature.PqcSignatures[:1]
			},
		},
		{
			name: "duplicate signature",
			mutate: func(signature *crypto.Signature) {
				signature.PqcSignatures[1].Type = crypto.SigTypeDilithium3
			},
		},
		{
			name: "signature without certificate key",
			mutate: func(signature *crypto.Signature) {
				signature.PqcSignatures[1].Type = crypto.SigTypeFalcon1024
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			signature := validCanonicalSignature()
			test.mutate(signature)
			candidates, err := pqcaddress.SignerAddressCandidates(signature.PqcCert, address.Delegated)
			if err != nil {
				t.Fatal(err)
			}
			if err := MultiPqcVerify(signature, candidates[0], []byte("message")); err == nil {
				t.Fatal("expected invalid envelope to be rejected")
			}
		})
	}
}

func TestMultiPqcVerifyRejectsCrossVersionAddresses(t *testing.T) {
	installAcceptingPqcShims(t)
	for _, protocol := range []address.Protocol{address.PQC, address.Delegated} {
		signature := validCanonicalSignature()
		modern, err := pqcaddress.SignerAddressCandidates(signature.PqcCert, protocol)
		if err != nil {
			t.Fatal(err)
		}
		legacyCert := signature.PqcCert
		legacyCert.Version = uint8(pqcaddress.SchemeLegacyCompatible)
		legacy, err := pqcaddress.SignerAddressCandidates(legacyCert, protocol)
		if err != nil {
			t.Fatal(err)
		}
		for _, candidate := range legacy {
			if err := MultiPqcVerify(signature, candidate, []byte("message")); err == nil {
				t.Fatal("version 1 accepted legacy CBOR or raw address")
			}
		}
		signature.PqcCert = legacyCert
		if err := MultiPqcVerify(signature, modern[0], []byte("message")); err == nil {
			t.Fatal("version 0 accepted version 1 CBOR address")
		}
	}
}
