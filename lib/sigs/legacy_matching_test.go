package sigs

import (
	"fmt"
	"testing"

	"github.com/post-quantumqoin/address"
	"github.com/post-quantumqoin/core-types/crypto"
	"github.com/post-quantumqoin/qoin-shor/core/types/pqcaddress"
)

type rejectingLegacyShim struct{ acceptingPqcShim }

func (rejectingLegacyShim) PqcVerify([]byte, []byte, []byte) error {
	return fmt.Errorf("invalid test signature")
}

func TestLegacySecondaryKeyBindingRequiresAllSignatures(t *testing.T) {
	installAcceptingPqcShims(t)
	for _, protocol := range []address.Protocol{address.PQC, address.Delegated} {
		sig := validCanonicalSignature()
		sig.PqcCert.Version = 0
		second := crypto.SignPQCCert{Version: 0, Pubkeys: sig.PqcCert.Pubkeys[1:]}
		candidates, err := pqcaddress.SignerAddressCandidates(second, protocol)
		if err != nil {
			t.Fatal(err)
		}
		for _, candidate := range candidates {
			if err := MultiPqcVerify(sig, candidate, []byte("message")); err != nil {
				t.Fatalf("secondary signer %s: %v", candidate, err)
			}
			for _, typ := range []crypto.SigType{crypto.SigTypeDilithium3, crypto.SigTypeFalcon512} {
				Pqcsigs[typ] = rejectingLegacyShim{}
				if err := MultiPqcVerify(sig, candidate, []byte("message")); err == nil {
					t.Fatalf("invalid signature type %d accepted", typ)
				}
				if err := MultiPqcVerifyWithAddressBinding(sig, candidate, []byte("message"), false); err == nil {
					t.Fatal("historical bypass skipped signature verification")
				}
				Pqcsigs[typ] = acceptingPqcShim{}
			}
		}
		sig.PqcSignatures = sig.PqcSignatures[:1]
		if err := MultiPqcVerify(sig, candidates[0], []byte("message")); err == nil {
			t.Fatal("missing matched-key signature accepted")
		}
		modern := validCanonicalSignature()
		modernSecond := crypto.SignPQCCert{Version: 1, Pubkeys: modern.PqcCert.Pubkeys[1:]}
		modernCandidates, err := pqcaddress.SignerAddressCandidates(modernSecond, protocol)
		if err != nil {
			t.Fatal(err)
		}
		if err := MultiPqcVerify(modern, modernCandidates[0], []byte("message")); err == nil {
			t.Fatal("version 1 accepted secondary key")
		}
	}
}
