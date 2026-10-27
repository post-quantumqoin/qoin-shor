package ethtypes

import (
	"testing"

	"github.com/post-quantumqoin/address"
)

func TestQoinAddressConversionRoundTrip(t *testing.T) {
	for _, input := range []EthAddress{
		{},
		EthAddressFromActorID(1840),
		{0x12, 0x34, 0x56, 0x78, 0x90},
	} {
		converted, err := input.ToQoinAddress()
		if err != nil {
			t.Fatal(err)
		}
		wantProtocol := address.Delegated
		if input.IsMaskedID() {
			wantProtocol = address.ID
		}
		if converted.Protocol() != wantProtocol {
			t.Fatalf("unexpected address protocol for %s: got %v, want %v", input, converted.Protocol(), wantProtocol)
		}
		restored, err := EthAddressFromQoinAddress(converted)
		if err != nil {
			t.Fatal(err)
		}
		if restored != input {
			t.Fatalf("Qoin address round trip changed %s to %s", input, restored)
		}
	}
}

func TestQoinAddressConversionRejectsUnsupportedProtocols(t *testing.T) {
	input, err := address.NewPqcAddress([]byte{1, 2, 3})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := EthAddressFromQoinAddress(input); err != ErrInvalidAddress {
		t.Fatalf("expected ErrInvalidAddress for a PQC address, got %v", err)
	}
}
