package key

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/post-quantumqoin/address"
	"github.com/post-quantumqoin/qoin-shor/core/types"
)

func versionTestKeyInfo(version *uint8) types.KeyInfo {
	return types.KeyInfo{Type: types.KTPqc, PqcAddressVersion: version, PqcKeypairs: []types.PqcKeypair{
		{PqcType: types.KeyType(types.Falcon512), PqcPublicKey: []byte{1, 2, 3}},
		{PqcType: types.KeyType(types.Dilithium3), PqcPublicKey: []byte{4, 5, 6}},
	}}
}

func TestPqcWalletVersionedAddresses(t *testing.T) {
	for _, version := range []uint8{0, 1} {
		info := versionTestKeyInfo(&version)
		if version == 0 {
			info.PqcAddressVersion = nil
		}
		key, err := NewPqcKey(info)
		if err != nil {
			t.Fatal(err)
		}
		wantFirst, encoded := "falcon512", "816966616c636f6e35313243010203"
		if version == 1 {
			wantFirst, encoded = "dilithium3", "826a64696c69746869756d3343040506"
		}
		payload, err := hex.DecodeString(encoded)
		if err != nil {
			t.Fatal(err)
		}
		want, err := address.NewPqcAddress(payload)
		if err != nil {
			t.Fatal(err)
		}
		if key.PQCCert.Version != version || key.PQCCert.Pubkeys[0].Typ != wantFirst || key.Address != want {
			t.Fatalf("version %d: wrong certificate order or address", version)
		}
		if info.PqcKeypairs[0].PqcType != types.KeyType(types.Falcon512) {
			t.Fatal("input keypairs mutated")
		}
		data, err := json.Marshal(key.KeyInfo)
		if err != nil {
			t.Fatal(err)
		}
		var restored types.KeyInfo
		if err := json.Unmarshal(data, &restored); err != nil {
			t.Fatal(err)
		}
		loaded, err := NewPqcKey(restored)
		if err != nil {
			t.Fatal(err)
		}
		if loaded.Address != key.Address || loaded.PQCCert.Version != version {
			t.Fatal("reload changed address or version")
		}
	}
}

func TestLegacyRawAddressLoadAndExport(t *testing.T) {
	info := versionTestKeyInfo(nil)
	raw, err := address.NewPqcAddress(info.PqcKeypairs[0].PqcPublicKey)
	if err != nil {
		t.Fatal(err)
	}
	key, err := NewPqcKeyForAddress(info, raw)
	if err != nil {
		t.Fatal(err)
	}
	loaded, err := NewPqcKey(key.KeyInfo)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Address != raw || loaded.PQCCert.Version != 0 {
		t.Fatal("legacy raw account changed on reload")
	}
	modernVersion := uint8(1)
	modern, err := NewPqcKey(versionTestKeyInfo(&modernVersion))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := NewPqcKeyForAddress(info, modern.Address); err == nil {
		t.Fatal("mismatched stored key accepted")
	}
	modernInfo := modern.KeyInfo
	modernInfo.PqcAddress = raw.String()
	if _, err := NewPqcKey(modernInfo); err == nil {
		t.Fatal("version 1 accepted raw legacy address")
	}
}

func TestPqcWalletRejectsUnknownVersion(t *testing.T) {
	version := uint8(2)
	if _, err := NewPqcKey(versionTestKeyInfo(&version)); err == nil {
		t.Fatal("unknown version accepted")
	}
}
