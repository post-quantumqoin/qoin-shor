package key

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/post-quantumqoin/address"
	"github.com/post-quantumqoin/core-types/crypto"
	"github.com/post-quantumqoin/qoin-shor/core/types"
	"github.com/post-quantumqoin/qoin-shor/core/types/pqcaddress"
)

func TestLoadLegacySortedAndOriginalWalletAddresses(t *testing.T) {
	for _, protocol := range []address.Protocol{address.PQC, address.Delegated} {
		info := versionTestKeyInfo(nil)
		if protocol == address.Delegated {
			info.Type = types.KTDelegated
		}
		before, _ := json.Marshal(info)
		cert := crypto.SignPQCCert{Version: 0}
		for _, kp := range info.PqcKeypairs {
			cert.Pubkeys = append(cert.Pubkeys, crypto.SignPqcCertPubkey{Typ: string(kp.PqcType), Pubkey: kp.PqcPublicKey})
		}
		candidates, err := pqcaddress.LegacyWalletAddressCandidates(cert, protocol)
		if err != nil {
			t.Fatal(err)
		}
		initial, err := NewPqcKey(info)
		if err != nil || initial.Address != candidates[0].Address {
			t.Fatalf("default generation changed: %v", err)
		}
		for index, candidate := range candidates {
			loaded, err := NewPqcKeyForAddress(info, candidate.Address)
			if err != nil {
				t.Fatal(err)
			}
			if loaded.Address != candidate.Address || loaded.PqcAddress != candidate.Address.String() || *loaded.PqcAddressVersion != 0 {
				t.Fatal("historical identity not preserved")
			}
			if !reflect.DeepEqual(loaded.PQCCert.Pubkeys, initial.PQCCert.Pubkeys) {
				t.Fatal("loading changed certificate order")
			}
			data, _ := json.Marshal(loaded.KeyInfo)
			var restored types.KeyInfo
			if err := json.Unmarshal(data, &restored); err != nil {
				t.Fatal(err)
			}
			reloaded, err := NewPqcKey(restored)
			if err != nil || reloaded.Address != candidate.Address {
				t.Fatalf("metadata reload changed address: %v", err)
			}
			if _, err := NewPqcKeyForAddress(restored, candidates[(index+1)%len(candidates)].Address); err == nil {
				t.Fatal("conflicting requested and stored addresses accepted")
			}
		}
		after, _ := json.Marshal(info)
		if string(before) != string(after) {
			t.Fatal("input KeyInfo mutated")
		}
		version := uint8(1)
		modernInfo := info
		modernInfo.PqcAddressVersion = &version
		for _, candidate := range candidates {
			if _, err := NewPqcKeyForAddress(modernInfo, candidate.Address); err == nil {
				t.Fatal("version 1 accepted legacy wallet candidate")
			}
		}
	}
}
