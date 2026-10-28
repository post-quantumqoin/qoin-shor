package ethtypes

import (
	"testing"

	"github.com/post-quantumqoin/address"
	typescrypto "github.com/post-quantumqoin/core-types/crypto"

	"github.com/post-quantumqoin/qoin-shor/core/types/pqcaddress"
)

func pqcSenderTestTransaction(version pqcaddress.Scheme) *EthTxArgs {
	signature := typescrypto.Signature{
		Type: typescrypto.SigTypeDelegated,
		PqcCert: typescrypto.SignPQCCert{
			Version: uint8(version),
			Pubkeys: []typescrypto.SignPqcCertPubkey{
				{Typ: "dilithium3", Pubkey: []byte{1, 2, 3}},
				{Typ: "falcon512", Pubkey: []byte{4, 5, 6}},
			},
		},
		PqcSignatures: []typescrypto.PqcSignature{
			{Type: typescrypto.SigTypeDilithium3, Data: []byte{7}},
			{Type: typescrypto.SigTypeFalcon512, Data: []byte{8}},
		},
	}
	return &EthTxArgs{PQCSig: (*PQCSignature)(&signature)}
}

func TestPqcSenderUsesExplicitCanonicalScheme(t *testing.T) {
	tx := pqcSenderTestTransaction(pqcaddress.SchemeCanonicalCBOR)
	sender, err := tx.Sender()
	if err != nil {
		t.Fatal(err)
	}

	signature := (*typescrypto.Signature)(tx.PQCSig)
	candidates, err := pqcaddress.SignerAddressCandidates(signature.PqcCert, address.Delegated)
	if err != nil {
		t.Fatal(err)
	}
	if len(candidates) != 1 || sender != candidates[0] {
		t.Fatal("PQC sender does not match the explicit canonical certificate scheme")
	}
}

func TestPqcSenderRejectsAmbiguousLegacyScheme(t *testing.T) {
	tx := pqcSenderTestTransaction(pqcaddress.SchemeLegacyCompatible)
	if _, err := tx.Sender(); err == nil {
		t.Fatal("legacy scheme without an explicit CBOR/raw choice must not recover an Ethereum sender")
	}
}
