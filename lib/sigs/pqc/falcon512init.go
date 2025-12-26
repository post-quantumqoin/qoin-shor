package pqc

import (
	"fmt"

	"github.com/post-quantumqoin/core-types/crypto"
	"github.com/post-quantumqoin/qoin-shor/pqccrypto/nistround3/falcon/falcon512"

	"github.com/post-quantumqoin/qoin-shor/lib/sigs"
)

type falcon512Signer struct{}

func (falcon512Signer) PqcGenPrivate() ([]byte, []byte, []byte, error) {
	// fmt.Println("falcon512.PqcGenPrivate")
	seedbytes, skbytes, pkbytes, err := falcon512.GenerateKey()
	if err != nil {
		return nil, nil, nil, err
	}
	// fmt.Println("falcon512.PqcGenPrivate skbytes:",skbytes)
	return seedbytes, skbytes, pkbytes, nil
}

func (falcon512Signer) PqcToPublic(seed []byte) ([]byte, error) {
	pkbytes := falcon512.GenPkBySeed(seed)

	return pkbytes, nil
}

// msg:"hello"
func (falcon512Signer) PqcSign(sk []byte, msg []byte) ([]byte, error) {
	sig, err := falcon512.Sign(sk, msg)
	if err != nil {
		return nil, err
	}

	return sig, nil
}

// msg:"hello"
func (falcon512Signer) PqcVerify(sig []byte, pk []byte, msg []byte) error {
	if falcon512.Verify(pk, msg, sig) {
		return nil
	}

	return fmt.Errorf("PqcVerify is fail")
}

func init() {
	// fmt.Println("falcon512 init")
	sigs.PqcRegisterSignature(crypto.SigTypeFalcon512, falcon512Signer{})
}
