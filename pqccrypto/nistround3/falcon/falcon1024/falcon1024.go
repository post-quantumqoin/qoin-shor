package falcon1024

/*
#cgo CFLAGS: -I./falconCore/ -DFALCON_PREFIX=falcon_inner1024
#cgo LDFLAGS: -L${SRCDIR}/falconCore/build -lkat1024int


#include <stdio.h>
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

// Falcon1024CryptoNonceBytes is the size of a serialized nonce key.
const Falcon1024CryptoNonceBytes = 40 

// Falcon1024CryptoSeedBytes is the size of a serialized seed key.
const Falcon1024CryptoSeedBytes = 48

// Falcon1024CryptoSkBytes is the size of a serialized private key.
const Falcon1024CryptoSkBytes = 2305

// Falcon1024CryptoPkBytes is the size of a serialized public key.
const Falcon1024CryptoPkBytes = 1793

// Falcon1024CryptoSignBytes is the size of a serialized sign key.
const Falcon1024CryptoSignBytes = 1332

// Private returns the public key for this  seed.
func GenSkBySeed(seed []byte) []byte {
	if(Falcon1024CryptoSeedBytes != len(seed)){ 
		fmt.Println(errors.New("err:Falcon1024CryptoSeedBytes is not match"))
		return nil
	}
	var sk []C.uchar = make([]C.uchar, Falcon1024CryptoSkBytes)
	seedptr := C.CBytes(seed)

	err := int32(C.falcon_seed_to_sk(&sk[0], (*C.uchar)(seedptr)))
	if err != 0 {
		fmt.Println(errors.New("err: falcon_seed_to_sk"))
		return nil
	}
	C.free(seedptr)

	return C.GoBytes(unsafe.Pointer(&sk[0]), C.int(Falcon1024CryptoSkBytes))
}

// PublicKey returns the public key for this  seed.
func GenPkBySeed(seed []byte) []byte {
	if(Falcon1024CryptoSeedBytes != len(seed)){ 
		fmt.Println(errors.New("err:Falcon1024CryptoSeedBytes is not match"))
		return nil
	}
	var pk []C.uchar = make([]C.uchar, Falcon1024CryptoPkBytes)
	seedptr := C.CBytes(seed)

	err := int32(C.falcon_seed_to_pk(&pk[0], (*C.uchar)(seedptr)))
	if err != 0 {
		fmt.Println(errors.New("err: falcon_seed_to_sk"))
		return nil
	}
	C.free(seedptr)

	return C.GoBytes(unsafe.Pointer(&pk[0]), C.int(Falcon1024CryptoPkBytes))
}

// Sign signs the given message, which must be 32 bytes long.
func Sign(sk, msg []byte) ([]byte, error) {
	if(Falcon1024CryptoSkBytes != len(sk)){ 
		return nil, errors.New("err:Falcon1024CryptoSkBytes is not match")
	}

	var sg []C.uchar = make([]C.uchar, Falcon1024CryptoSignBytes)

	msgptr := C.CBytes(msg)
	skptr := C.CBytes(sk)

	err := int32(C.falcon_sign(&sg[0], (*C.uchar)(msgptr), C.ulonglong(len(msg)), (*C.uchar)(skptr)))
	if err != 0 {
		fmt.Println(err)
		return nil, errors.New("err: falcon_sign")
	}
	C.free(msgptr)
	C.free(skptr)

	//return Signature, nil
	return C.GoBytes(unsafe.Pointer(&sg[0]), C.int(Falcon1024CryptoSignBytes)), nil		
}

// Sign signs the given message, which must be 32 bytes long.
func SignBySeed(seed, msg []byte) ([]byte, error) {
	if(Falcon1024CryptoSeedBytes != len(seed)){ 
		return nil, errors.New("err:Falcon1024CryptoSeedBytes is not match")
	}

	var sk []C.uchar = make([]C.uchar, Falcon1024CryptoSkBytes)
	//gen privkey rom seed
	seedptr := C.CBytes(seed)
    
	err := int32(C.falcon_seed_to_sk(&sk[0], (*C.uchar)(seedptr)))
	if err != 0 {
		fmt.Println(err)
		return nil, errors.New("err: falcon_seed_to_sk")
	}

	var sg []C.uchar = make([]C.uchar, Falcon1024CryptoSignBytes)
	msgptr := C.CBytes(msg)
	//signs msg  
	serr := int32(C.falcon_sign(&sg[0], (*C.uchar)(msgptr), C.ulonglong(len(msg)), &sk[0]))
	if serr != 0 {
		fmt.Println(serr)
		return nil, errors.New("err: falcon_sign")
	}
	C.free(msgptr)
	C.free(seedptr)

	//return Signature, nil
	return C.GoBytes(unsafe.Pointer(&sg[0]), C.int(Falcon1024CryptoSignBytes)), nil		
}

// Equals compares two private key for equality and returns true if they are the same.
func Equals(sk, other []byte) bool {
	return bytes.Equal(sk, other)
}

// Verify checks the given signature and returns true if it is valid.
func Verify(pk, msg, signature []byte) bool {
	if(Falcon1024CryptoSignBytes != len(signature)){ 
		fmt.Println(errors.New("err: Falcon1024CryptoSignBytes is not match"))
		return false
	}

	sgptr := C.CBytes(signature)
	msgptr := C.CBytes(msg)
	pkptr := C.CBytes(pk)
	
	v := int32(C.verify_sign((*C.uchar)(msgptr), C.ulonglong(len(msg)), (*C.uchar)(sgptr), (*C.uchar)(pkptr)))
	if v != 0 {
		fmt.Println(errors.New("err: verify_sign"))
		return false
	}

	C.free(pkptr)
	C.free(msgptr)
	C.free(sgptr)
	return true
}

// GenerateKeyFromSeed generates a new key from the given reader.
func GenerateKeyFromSeed(seed []byte) ([]byte, error) {
	if(Falcon1024CryptoSeedBytes != len(seed)){ 
		return nil, errors.New("err:Falcon1024CryptoSeedBytes is not match")
	}

	var pk []C.uchar = make([]C.uchar, Falcon1024CryptoPkBytes)
	var sk []C.uchar = make([]C.uchar, Falcon1024CryptoSkBytes)

	fmt.Println("sk:",sk,"\n")

	seedptr := C.CBytes(seed)
	err := int32(C.falcon_genkey_by_seed(&pk[0], &sk[0], (*C.uchar)(seedptr)));
	if err != 0 {
		fmt.Println(err)
		return nil, errors.New("err: falcon_genkey_by_seed")
	}
	C.free(seedptr)

	//return privkey, nil
	return C.GoBytes(unsafe.Pointer(&sk[0]), C.int(Falcon1024CryptoSkBytes)), nil
}

// GenerateKey creates a new key using secure randomness from crypto.rand.
func GenerateKey() ([]byte, []byte, []byte, error) {
	var pk []C.uchar = make([]C.uchar, Falcon1024CryptoPkBytes)
	var sk []C.uchar = make([]C.uchar, Falcon1024CryptoSkBytes)
	var seed []C.uchar = make([]C.uchar, Falcon1024CryptoSeedBytes)

	err := int32(C.falcon_genkey(&pk[0], &sk[0], &seed[0]))
	if err != 0 {
		fmt.Println(err)
		return nil, nil, nil, errors.New("err: falcon genkey ")
	}
	seedbytes := C.GoBytes(unsafe.Pointer(&seed[0]), C.int(Falcon1024CryptoSeedBytes))
	skbytes := C.GoBytes(unsafe.Pointer(&sk[0]), C.int(Falcon1024CryptoSkBytes))
	pkbytes := C.GoBytes(unsafe.Pointer(&pk[0]), C.int(Falcon1024CryptoPkBytes))
	//return seed , privkey, publickey nil
	return seedbytes, skbytes, pkbytes, nil
}


