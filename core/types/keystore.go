package types

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/post-quantumqoin/core-types/crypto"
	// "github.com/ipfs/go-cid"
)

var (
	ErrKeyInfoNotFound = fmt.Errorf("key info not found")
	ErrKeyExists       = fmt.Errorf("key already exists")
)

// KeyType defines a type of a key
type KeyType string

func (kt *KeyType) UnmarshalJSON(bb []byte) error {
	{
		// first option, try unmarshaling as string
		var s string
		err := json.Unmarshal(bb, &s)
		if err == nil {
			*kt = KeyType(s)
			return nil
		}
	}

	{
		var b byte
		err := json.Unmarshal(bb, &b)
		if err != nil {
			return fmt.Errorf("could not unmarshal KeyType either as string nor integer: %w", err)
		}
		bst := crypto.SigType(b)

		switch bst {
		case crypto.SigTypeBLS:
			*kt = KTBLS
		case crypto.SigTypeSecp256k1:
			*kt = KTSecp256k1
		case crypto.SigTypeDelegated:
			*kt = KTDelegated
		default:
			return fmt.Errorf("unknown sigtype: %d", bst)
		}
		log.Warnf("deprecation: integer style 'KeyType' is deprecated, switch to string style")
		return nil
	}
}

const (
	KTBLS             KeyType = "bls"
	KTSecp256k1       KeyType = "secp256k1"
	KTSecp256k1Ledger KeyType = "secp256k1-ledger"

	KTDelegated KeyType = "delegated"
	KTPqc       KeyType = "pqc"

	Falcon512  KeyType = "falcon512"
	Falcon1024 KeyType = "falcon1024"
	Dilithium3 KeyType = "dilithium3"
	Dilithium5 KeyType = "dilithium5"
)

type PqcKeypair struct {
	PqcVersion uint8
	PqcSeed    []byte
	PqcType    KeyType

	PqcPrivateKey []byte
	PqcPublicKey  []byte
}

type PqcCertPubkey struct {
	Typ    string
	Pubkey []byte
}

func (t *PqcCertPubkey) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := t.MarshalCBOR(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type PQCCert struct {
	Pubkeys []PqcCertPubkey
	Version uint8
	// Nonce   []byte
}

func (pc *PQCCert) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := pc.MarshalCBOR(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// KeyInfo is used for storing keys in KeyStore
type KeyInfo struct {
	Type        KeyType
	PrivateKey  []byte
	PqcKeypairs []PqcKeypair
	// PQCCert
}

type PqcKeyInfo struct {
	Keypairs map[string]PqcKeypair
	PQCCert
}

// KeyStore is used for storing secret keys
type KeyStore interface {
	// List lists all the keys stored in the KeyStore
	List() ([]string, error)
	// Get gets a key out of keystore and returns KeyInfo corresponding to named key
	Get(string) (KeyInfo, error)
	// Put saves a key info under given name
	Put(string, KeyInfo) error
	// Delete removes a key from keystore
	Delete(string) error
}
