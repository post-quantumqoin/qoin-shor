package consensus

import (
	"testing"

	"github.com/post-quantumqoin/address"
	"github.com/post-quantumqoin/core-types/abi"
	"github.com/post-quantumqoin/core-types/big"
	"github.com/post-quantumqoin/core-types/crypto"
	"github.com/post-quantumqoin/qoin-shor/build"
	"github.com/post-quantumqoin/qoin-shor/core/types"
	"github.com/post-quantumqoin/qoin-shor/core/types/pqcaddress"
	"github.com/post-quantumqoin/qoin-shor/lib/sigs"
)

func TestLegacySecondarySignerAcrossBindingHeight(t *testing.T) {
	previous := sigs.Pqcsigs
	sigs.Pqcsigs = map[crypto.SigType]sigs.PqcSigShim{
		crypto.SigTypeDilithium3: heightTestPqcShim{},
		crypto.SigTypeFalcon512:  heightTestPqcShim{},
	}
	t.Cleanup(func() { sigs.Pqcsigs = previous })
	cert := crypto.SignPQCCert{Version: 0, Pubkeys: []crypto.SignPqcCertPubkey{
		{Typ: "dilithium3", Pubkey: []byte{1, 2, 3}},
		{Typ: "falcon512", Pubkey: []byte{4, 5, 6}},
	}}
	for _, protocol := range []address.Protocol{address.PQC, address.Delegated} {
		candidates, err := pqcaddress.SignerAddressCandidates(crypto.SignPQCCert{Version: 0, Pubkeys: cert.Pubkeys[1:]}, protocol)
		if err != nil {
			t.Fatal(err)
		}
		for _, signer := range candidates {
			msg := &types.SignedMessage{
				Message:   types.Message{From: signer, To: signer, Value: big.Zero(), GasLimit: 30000000, GasFeeCap: big.Zero(), GasPremium: big.Zero()},
				Signature: crypto.Signature{Type: crypto.SigTypeMultiPqc, PqcCert: cert},
			}
			digest := msg.Message.Cid().Bytes()
			msg.Signature.PqcSignatures = []crypto.PqcSignature{
				{Type: crypto.SigTypeDilithium3, Data: digest},
				{Type: crypto.SigTypeFalcon512, Data: digest},
			}
			h := build.UpgradePqcAddressBindingHeight
			for _, height := range []abi.ChainEpoch{h - 1, h, h + 1} {
				if err := AuthenticateMessageAtHeight(msg, signer, height); err != nil {
					t.Fatalf("secondary signer rejected at %d: %v", height, err)
				}
			}
			if err := AuthenticateMessage(msg, signer); err != nil {
				t.Fatalf("message-pool authentication rejected secondary signer: %v", err)
			}
			for index := range msg.Signature.PqcSignatures {
				msg.Signature.PqcSignatures[index].Data = []byte("invalid")
				for _, height := range []abi.ChainEpoch{h - 1, h, h + 1} {
					if err := AuthenticateMessageAtHeight(msg, signer, height); err == nil {
						t.Fatalf("invalid signature %d accepted at %d", index, height)
					}
				}
				if err := AuthenticateMessage(msg, signer); err == nil {
					t.Fatalf("message-pool authentication accepted invalid signature %d", index)
				}
				msg.Signature.PqcSignatures[index].Data = digest
			}
		}
	}
}
