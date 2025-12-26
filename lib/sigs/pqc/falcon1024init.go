package pqc

import (
	"fmt"

	"github.com/post-quantumqoin/core-types/crypto"
	"github.com/post-quantumqoin/qoin-shor/pqccrypto/nistround3/falcon/falcon1024"

	"github.com/post-quantumqoin/qoin-shor/lib/sigs"
)

type falcon1024Signer struct{}

func (falcon1024Signer) PqcGenPrivate() ([]byte, []byte, []byte, error) {
	seedbytes, skbytes, pkbytes, err := falcon1024.GenerateKey()
	if err != nil {
		return nil, nil, nil, err
	}
	return seedbytes, skbytes, pkbytes, nil
}

func (falcon1024Signer) PqcToPublic(seed []byte) ([]byte, error) {
	pkbytes := falcon1024.GenPkBySeed(seed)

	return pkbytes, nil
}

// msg:"hello"
func (falcon1024Signer) PqcSign(sk []byte, msg []byte) ([]byte, error) {
	sig, err := falcon1024.Sign(sk, msg)
	if err != nil {
		return nil, err
	}

	return sig, nil
}

// msg:"hello"
func (falcon1024Signer) PqcVerify(sig []byte, pk []byte, msg []byte) error {
	if falcon1024.Verify(pk, msg, sig) {
		return nil
	}

	return fmt.Errorf("PqcVerify is fail")
}

func init() {
	sigs.PqcRegisterSignature(crypto.SigTypeFalcon1024, falcon1024Signer{})
}
