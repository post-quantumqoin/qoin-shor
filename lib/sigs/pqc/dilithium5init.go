package pqc

import (
	"fmt"

	"github.com/post-quantumqoin/core-types/crypto"
	"github.com/post-quantumqoin/qoin-shor/pqccrypto/nistround3/dilithium/dilithium5"

	"github.com/post-quantumqoin/qoin-shor/lib/sigs"
)

type dilithium5Signer struct{}

func (dilithium5Signer) PqcGenPrivate() ([]byte, []byte, []byte, error) {
	seedbytes, skbytes, pkbytes, err := dilithium5.GenerateKey()
	if err != nil {
		return nil, nil, nil, err
	}
	return seedbytes, skbytes, pkbytes, nil
}

func (dilithium5Signer) PqcToPublic(seed []byte) ([]byte, error) {
	pkbytes := dilithium5.GenPkBySeed(seed)

	return pkbytes, nil
}

// msg:"hello"
func (dilithium5Signer) PqcSign(sk []byte, msg []byte) ([]byte, error) {
	sig, err := dilithium5.Sign(sk, msg)
	if err != nil {
		return nil, err
	}

	return sig, nil
}

// msg:"hello"
func (dilithium5Signer) PqcVerify(sig []byte, pk []byte, msg []byte) error {
	if dilithium5.Verify(pk, msg, sig) {
		return nil
	}

	return fmt.Errorf("PqcVerify is fail")
}

func init() {
	sigs.PqcRegisterSignature(crypto.SigTypeDilithium5, dilithium5Signer{})
}
