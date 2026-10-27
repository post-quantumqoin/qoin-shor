package cli

import (
	"testing"

	"github.com/post-quantumqoin/address"
	"github.com/post-quantumqoin/qoin-shor/core/types"
	"github.com/post-quantumqoin/qoin-shor/core/types/ethtypes"
)

func TestBindImportAddress(t *testing.T) {
	a, _ := address.NewPqcAddress([]byte("import expected"))
	b, _ := address.NewPqcAddress([]byte("import conflicting"))
	ki := types.KeyInfo{Type: types.KTPqc}
	if err := bindImportAddress(&ki, ""); err != nil || ki.PqcAddress != "" {
		t.Fatal("empty flag changed input")
	}
	if err := bindImportAddress(&ki, a.String()); err != nil || ki.PqcAddress != a.String() {
		t.Fatal("expected address not bound")
	}
	if err := bindImportAddress(&ki, a.String()); err != nil {
		t.Fatal(err)
	}
	if err := bindImportAddress(&ki, b.String()); err == nil || ki.PqcAddress != a.String() {
		t.Fatal("conflicting metadata accepted")
	}
	if err := bindImportAddress(&ki, "invalid"); err == nil {
		t.Fatal("invalid address accepted")
	}
	other := types.KeyInfo{Type: types.KTBLS}
	if err := bindImportAddress(&other, a.String()); err == nil {
		t.Fatal("non-PQC key accepted")
	}
	delegated := types.KeyInfo{Type: types.KTDelegated}
	eth, err := ethtypes.ParseEthAddress("0x1234567890123456789012345678901234567890")
	if err != nil {
		t.Fatal(err)
	}
	expected, err := eth.ToQoinAddress()
	if err != nil {
		t.Fatal(err)
	}
	if err := bindImportAddress(&delegated, "0x1234567890123456789012345678901234567890"); err != nil || delegated.PqcAddress != expected.String() {
		t.Fatalf("hex expected address conversion failed: %v", err)
	}
}
