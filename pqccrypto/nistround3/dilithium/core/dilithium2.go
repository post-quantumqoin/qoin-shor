package dilithium2

/*
#cgo CFLAGS: -I./

#define DILITHIUM_MODE=2
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

// Dilithium2CryptoNonceBytes is the size of a serialized nonce key.
const Dilithium2CryptoNonceBytes = 48 

// Dilithium2CryptoSeedBytes is the size of a serialized seed key.
const Dilithium2CryptoSeedBytes = 32

// Dilithium2CryptoSkBytes is the size of a serialized private key.
const Dilithium2CryptoSkBytes = 2544

// Dilithium2CryptoPkBytes is the size of a serialized public key.
const Dilithium2CryptoPkBytes = 1312

// Dilithium2CryptoSignBytes is the size of a serialized sign key.
const Dilithium2CryptoSignBytes = 2420

// Private returns the public key for this  seed.
func GenSkBySeed(seed []byte) []byte {
	if(Dilithium2CryptoSeedBytes != len(seed)){ 
		fmt.Println(errors.New("err:Dilithium2CryptoSeedBytes is not match"))
		return nil
	}
	var sk []C.uchar = make([]C.uchar, Dilithium2CryptoSkBytes)
	seedptr := C.CBytes(seed)
	// const int gensk = crypto_generate_sk( sk, seed );
	err := int32(C.crypto_generate_sk(&sk[0], (*C.uchar)(seedptr)))
	if err != 0 {
		fmt.Println(errors.New("err: crypto_generate_sk"))
		return nil
	}
	C.free(seedptr)

	return C.GoBytes(unsafe.Pointer(&sk[0]), C.int(Dilithium2CryptoSkBytes))
}

// PublicKey returns the public key for this  seed.
func GenPkBySeed(seed []byte) []byte {
	if(Dilithium2CryptoSeedBytes != len(seed)){ 
		fmt.Println(errors.New("err: Dilithium2CryptoSeedBytes is not match"))
		return nil
	}
	var pk []C.uchar = make([]C.uchar, Dilithium2CryptoPkBytes)
	seedptr := C.CBytes(seed)
	// const int genpk = crypto_generate_pk( pk,seed );
	err := int32(C.crypto_generate_pk(&pk[0], (*C.uchar)(seedptr)))
	if err != 0 {
		fmt.Println(errors.New("err: crypto_generate_pk"))
		return nil
	}
	C.free(seedptr)

	return C.GoBytes(unsafe.Pointer(&pk[0]), C.int(Dilithium2CryptoPkBytes))
}

// Sign signs the given message, which must be 32 bytes long.
func Sign(sk, msg []byte) ([]byte, error) {
	if(Dilithium2CryptoSkBytes != len(sk)){ 
		return nil, errors.New("err: Dilithium2CryptoSkBytes is not match")
	}

	var sg []C.uchar = make([]C.uchar, Dilithium2CryptoSignBytes)

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
	return C.GoBytes(unsafe.Pointer(&sg[0]), C.int(Dilithium2CryptoSignBytes)), nil		
}

// Sign signs the given message, which must be 32 bytes long.
func SignBySeed(seed, msg []byte) ([]byte, error) {
	if(Dilithium2CryptoSeedBytes != len(seed)){ 
		return nil, errors.New("err: Dilithium2CryptoSeedBytes is not match")
	}

	var sk []C.uchar = make([]C.uchar, Dilithium2CryptoSkBytes)
	//gen privkey rom seed
	seedptr := C.CBytes(seed)
	// const int gensk = crypto_generate_sk( sk, seed );
	err := int32(C.crypto_generate_sk(&sk[0], (*C.uchar)(seedptr)))
	if err != 0 {
		fmt.Println(err)
		return nil, errors.New("err: crypto_generate_sk")
	}

	var sg []C.uchar = make([]C.uchar, Dilithium2CryptoSignBytes)
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
	return C.GoBytes(unsafe.Pointer(&sg[0]), C.int(Dilithium2CryptoSignBytes)), nil		
}

// Equals compares two private key for equality and returns true if they are the same.
func Equals(sk, other []byte) bool {
	return bytes.Equal(sk, other)
}

// Verify checks the given signature and returns true if it is valid.
func Verify(pk, msg, signature []byte) bool {
	if(Dilithium2CryptoSignBytes != len(signature)){ 
		fmt.Println(errors.New("err: Dilithium2CryptoSignBytes is not match"))
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
	if(Dilithium2CryptoSeedBytes != len(seed)){ 
		return nil, errors.New("err: Dilithium2CryptoSeedBytes is not match")
	}

	var pk []C.uchar = make([]C.uchar, Dilithium2CryptoPkBytes)
	var sk []C.uchar = make([]C.uchar, Dilithium2CryptoSkBytes)

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
	return C.GoBytes(unsafe.Pointer(&sk[0]), C.int(Dilithium2CryptoSkBytes)), nil
}

// GenerateKey creates a new key using secure randomness from crypto.rand.
func GenerateKey() ([]byte, []byte, []byte, error) {
	var pk []C.uchar = make([]C.uchar, Dilithium2CryptoPkBytes)
	var sk []C.uchar = make([]C.uchar, Dilithium2CryptoSkBytes)
	var seed []C.uchar = make([]C.uchar, Dilithium2CryptoSeedBytes)

	// const int genkey = crypto_sign_keypair( pk, sk, seed );
	err := int32(C.crypto_sign_keypair(&pk[0], &sk[0], &seed[0]))
	if err != 0 {
		fmt.Println(err)
		return nil, nil, nil, errors.New("err: crypto_sign_keypair ")
	}
	seedbytes := C.GoBytes(unsafe.Pointer(&seed[0]), C.int(Dilithium2CryptoSeedBytes))
	skbytes := C.GoBytes(unsafe.Pointer(&sk[0]), C.int(Dilithium2CryptoSkBytes))
	pkbytes := C.GoBytes(unsafe.Pointer(&pk[0]), C.int(Dilithium2CryptoPkBytes))
	//return seed , privkey, publickey nil
	return seedbytes, skbytes, pkbytes, nil
}


