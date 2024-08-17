package dilithium3

/*
#cgo CFLAGS: -I../core -DDILITHIUM_PREFIX=pqcrystals_dilithium3
#cgo LDFLAGS: -L. -L../core/build3 -lpqcrystals_dilithium3_ref 

#include <stdlib.h>
#include "api.h"
*/
import "C"
import (
	"bytes"
	// "io"
	"fmt"
	"unsafe"
	"errors"

)

// Dilithium3CryptoNonceBytes is the size of a serialized nonce key.
const Dilithium3CryptoNonceBytes = 48 

// Dilithium3CryptoSeedBytes is the size of a serialized seed key.
const Dilithium3CryptoSeedBytes = 32

// Dilithium3CryptoSkBytes is the size of a serialized private key.
const Dilithium3CryptoSkBytes = 4016

// Dilithium3CryptoPkBytes is the size of a serialized public key.
const Dilithium3CryptoPkBytes = 1952

// Dilithium3CryptoSignBytes is the size of a serialized sign key.
const Dilithium3CryptoSignBytes = 3293

// Private returns the public key for this  seed.
func GenSkBySeed(seed []byte) []byte {
	if(Dilithium3CryptoSeedBytes != len(seed)){ 
		fmt.Println(errors.New("err:Dilithium3CryptoSeedBytes is not match"))
		return nil
	}
	var sk []C.uchar = make([]C.uchar, Dilithium3CryptoSkBytes)
	seedptr := C.CBytes(seed)
	// const int gensk = crypto_generate_sk( sk, seed );
	err := int32(C.crypto_generate_sk(&sk[0], (*C.uchar)(seedptr)))
	if err != 0 {
		fmt.Println(errors.New("err: crypto_generate_sk"))
		return nil
	}
	C.free(seedptr)

	return C.GoBytes(unsafe.Pointer(&sk[0]), C.int(Dilithium3CryptoSkBytes))
}

// PublicKey returns the public key for this  seed.
func GenPkBySeed(seed []byte) []byte {
	if(Dilithium3CryptoSeedBytes != len(seed)){ 
		fmt.Println(errors.New("err: Dilithium3CryptoSeedBytes is not match"))
		return nil
	}
	var pk []C.uchar = make([]C.uchar, Dilithium3CryptoPkBytes)
	seedptr := C.CBytes(seed)
	// const int genpk = crypto_generate_pk( pk,seed );
	err := int32(C.crypto_generate_pk(&pk[0], (*C.uchar)(seedptr)))
	if err != 0 {
		fmt.Println(errors.New("err: crypto_generate_pk"))
		return nil
	}
	C.free(seedptr)

	return C.GoBytes(unsafe.Pointer(&pk[0]), C.int(Dilithium3CryptoPkBytes))
}

// Sign signs the given message, which must be 32 bytes long.
func Sign(sk, msg []byte) ([]byte, error) {
	if(Dilithium3CryptoSkBytes != len(sk)){ 
		return nil, errors.New("err: Dilithium3CryptoSkBytes is not match")
	}

	var sg []C.uchar = make([]C.uchar, Dilithium3CryptoSignBytes)

	msgptr := C.CBytes(msg)
	skptr := C.CBytes(sk)
	// crypto_sign_signature( sign_msg, text, text_length, sk, 0 );
	err := int32(C.crypto_sign_signature(&sg[0], (*C.uchar)(msgptr), C.ulonglong(len(msg)), (*C.uchar)(skptr), 0))
	if err != 0 {
		fmt.Println(errors.New("err: crypto_sign_signature"))
		return nil, errors.New("err: crypto_sign_signature")
	}
	C.free(msgptr)
	C.free(skptr)

	//return Signature, nil
	return C.GoBytes(unsafe.Pointer(&sg[0]), C.int(Dilithium3CryptoSignBytes)), nil		
}

// Sign signs the given message, which must be 32 bytes long.
func SignBySeed(seed, msg []byte) ([]byte, error) {
	if(Dilithium3CryptoSeedBytes != len(seed)){ 
		return nil, errors.New("err: Dilithium3CryptoSeedBytes is not match")
	}

	var sk []C.uchar = make([]C.uchar, Dilithium3CryptoSkBytes)
	//gen privkey rom seed
	seedptr := C.CBytes(seed)
	// const int gensk = crypto_generate_sk( sk, seed );
	err := int32(C.crypto_generate_sk(&sk[0], (*C.uchar)(seedptr)))
	if err != 0 {
		fmt.Println(err)
		return nil, errors.New("err: crypto_generate_sk")
	}

	var sg []C.uchar = make([]C.uchar, Dilithium3CryptoSignBytes)
	msgptr := C.CBytes(msg)
	//signs msg  
	// const int gensign = crypto_sign_signature( sign_msg, text, text_length, sk, 0 );
	serr := int32(C.crypto_sign_signature(&sg[0], (*C.uchar)(msgptr), C.ulonglong(len(msg)), &sk[0], 0))
	if serr != 0 {
		fmt.Println(serr)
		return nil, errors.New("err: crypto_sign_signature")
	}
	C.free(msgptr)
	C.free(seedptr)

	//return Signature, nil
	return C.GoBytes(unsafe.Pointer(&sg[0]), C.int(Dilithium3CryptoSignBytes)), nil		
}

// Equals compares two private key for equality and returns true if they are the same.
func Equals(sk, other []byte) bool {
	return bytes.Equal(sk, other)
}

// Verify checks the given signature and returns true if it is valid.
func Verify(pk, msg, signature []byte) bool {
	if(Dilithium3CryptoSignBytes != len(signature)){ 
		fmt.Println(errors.New("err: Dilithium3CryptoSignBytes is not match"))
		return false
	}

	sgptr := C.CBytes(signature)
	msgptr := C.CBytes(msg)
	pkptr := C.CBytes(pk)
	// const int verify_ans = crypto_sign_verify( sign, text, text_length, pk );
	v := int32(C.crypto_sign_verify((*C.uchar)(sgptr), (*C.uchar)(msgptr),  C.ulonglong(len(msg)), (*C.uchar)(pkptr)))
	if v != 0 {
		fmt.Println(errors.New("err: crypto_sign_verify"))
		return false
	}

	C.free(pkptr)
	C.free(msgptr)
	C.free(sgptr)
	return true
}

// GenerateKeyFromSeed generates a new key from the given reader.
func GenerateKeyFromSeed(seed []byte) ([]byte, error) {
	if(Dilithium3CryptoSeedBytes != len(seed)){ 
		return nil, errors.New("err: Dilithium3CryptoSeedBytes is not match")
	}

	var pk []C.uchar = make([]C.uchar, Dilithium3CryptoPkBytes)
	var sk []C.uchar = make([]C.uchar, Dilithium3CryptoSkBytes)

	fmt.Println("sk:",sk,"\n")

	seedptr := C.CBytes(seed)
	// const int genkey = crypto_sign_keypair_by_seed( pk, sk, seed );
	err := int32(C.crypto_sign_keypair_by_seed(&pk[0], &sk[0], (*C.uchar)(seedptr)));
	if err != 0 {
		fmt.Println(err)
		return nil, errors.New("err: crypto_sign_keypair_by_seed")
	}
	C.free(seedptr)

	//return privkey, nil
	return C.GoBytes(unsafe.Pointer(&sk[0]), C.int(Dilithium3CryptoSkBytes)), nil
}

// GenerateKey creates a new key using secure randomness from crypto.rand.
func GenerateKey() ([]byte, []byte, []byte, error) {
	var pk []C.uchar = make([]C.uchar, Dilithium3CryptoPkBytes)
	var sk []C.uchar = make([]C.uchar, Dilithium3CryptoSkBytes)
	var seed []C.uchar = make([]C.uchar, Dilithium3CryptoSeedBytes)

	// const int genkey = crypto_sign_keypair( pk, sk, seed );
	err := int32(C.crypto_sign_keypair(&pk[0], &sk[0], &seed[0]))
	if err != 0 {
		fmt.Println(err)
		return nil, nil, nil, errors.New("err: crypto_sign_keypair ")
	}
	seedbytes := C.GoBytes(unsafe.Pointer(&seed[0]), C.int(Dilithium3CryptoSeedBytes))
	skbytes := C.GoBytes(unsafe.Pointer(&sk[0]), C.int(Dilithium3CryptoSkBytes))
	pkbytes := C.GoBytes(unsafe.Pointer(&pk[0]), C.int(Dilithium3CryptoPkBytes))
	//return seed , privkey, publickey nil
	return seedbytes, skbytes, pkbytes, nil
}


