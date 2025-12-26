package pqc

import (
	"fmt"

	"github.com/post-quantumqoin/core-types/crypto"
	"github.com/post-quantumqoin/qoin-shor/pqccrypto/nistround3/dilithium/dilithium3"

	"github.com/post-quantumqoin/qoin-shor/lib/sigs"
)

type dilithium3Signer struct{}

func (dilithium3Signer) PqcGenPrivate() ([]byte, []byte, []byte, error) {
	seedbytes, skbytes, pkbytes, err := dilithium3.GenerateKey()
	if err != nil {
		return nil, nil, nil, err
	}
	return seedbytes, skbytes, pkbytes, nil
}

func (dilithium3Signer) PqcToPublic(seed []byte) ([]byte, error) {
	pkbytes := dilithium3.GenPkBySeed(seed)

	return pkbytes, nil
}

// msg:"hello"
func (dilithium3Signer) PqcSign(sk []byte, msg []byte) ([]byte, error) {
	sig, err := dilithium3.Sign(sk, msg)
	if err != nil {
		return nil, err
	}

	return sig, nil
}

// msg:"hello"
func (dilithium3Signer) PqcVerify(sig []byte, pk []byte, msg []byte) error {
	if dilithium3.Verify(pk, msg, sig) {
		return nil
	}

	return fmt.Errorf("PqcVerify is fail")
}

func init() {
	sigs.PqcRegisterSignature(crypto.SigTypeDilithium3, dilithium3Signer{})
}
