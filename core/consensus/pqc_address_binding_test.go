package consensus

import (
	"bytes"
	"fmt"
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

type heightTestPqcShim struct{}

func (heightTestPqcShim) PqcGenPrivate() ([]byte, []byte, []byte, error) { return nil, nil, nil, nil }
func (heightTestPqcShim) PqcToPublic([]byte) ([]byte, error)             { return nil, nil }
func (heightTestPqcShim) PqcSign([]byte, []byte) ([]byte, error)         { return nil, nil }
func (heightTestPqcShim) PqcVerify(signature, _, message []byte) error {
	if !bytes.Equal(signature, message) {
		return fmt.Errorf("invalid test signature")
	}
	return nil
}

func TestPqcAddressBindingHeight(t *testing.T) {
	previous := sigs.Pqcsigs
	sigs.Pqcsigs = map[crypto.SigType]sigs.PqcSigShim{crypto.SigTypeDilithium3: heightTestPqcShim{}}
	t.Cleanup(func() { sigs.Pqcsigs = previous })
	cert := crypto.SignPQCCert{Version: 0, Pubkeys: []crypto.SignPqcCertPubkey{{Typ: "dilithium3", Pubkey: []byte{1, 2, 3}}}}
	candidates, err := pqcaddress.SignerAddressCandidates(cert, address.PQC)
	if err != nil {
		t.Fatal(err)
	}
	wrongSigner, err := address.NewPqcAddress([]byte{99})
	if err != nil {
		t.Fatal(err)
	}
	msg := &types.SignedMessage{
		Message:   types.Message{From: wrongSigner, To: candidates[0], Nonce: 15, Value: big.Zero(), GasLimit: 30000000, GasFeeCap: big.Zero(), GasPremium: big.Zero()},
		Signature: crypto.Signature{Type: crypto.SigTypeMultiPqc, PqcCert: cert},
	}
	msg.Signature.PqcSignatures = []crypto.PqcSignature{{Type: crypto.SigTypeDilithium3, Data: msg.Message.Cid().Bytes()}}
	h := build.UpgradePqcAddressBindingHeight
	for _, test := range []struct {
		height  abi.ChainEpoch
		allowed bool
	}{{h - 1, true}, {h, true}, {h + 1, false}} {
		err := AuthenticateMessageAtHeight(msg, wrongSigner, test.height)
		if (err == nil) != test.allowed {
			t.Fatalf("height %d: allowed=%v, error=%v", test.height, test.allowed, err)
		}
	}
	if err := AuthenticateMessage(msg, wrongSigner); err == nil {
		t.Fatal("strict message-pool path accepted wrong signer")
	}
	if err := AuthenticateMessageAtHeight(msg, candidates[0], h+1); err != nil {
		t.Fatalf("valid signer rejected after activation: %v", err)
	}
	msg.Signature.PqcSignatures[0].Data = []byte("invalid")
	if err := AuthenticateMessageAtHeight(msg, wrongSigner, h); err == nil {
		t.Fatal("historical rule skipped cryptographic verification")
	}
	msg.Signature.PqcSignatures = nil
	if err := AuthenticateMessageAtHeight(msg, wrongSigner, h); err == nil {
		t.Fatal("historical rule accepted missing signature")
	}
}
