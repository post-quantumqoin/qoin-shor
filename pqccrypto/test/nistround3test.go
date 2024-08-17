package main

import "C"
import (
	// "bytes"
	// "io"
	"fmt"
	"pqccrypto/nistround3/dilithium/dilithium2"
	"pqccrypto/nistround3/dilithium/dilithium3"

	"pqccrypto/nistround3/falcon/falcon1024"
	"pqccrypto/nistround3/falcon/falcon512"
)

func TestDilithium2() {
	fmt.Println("start: TestDilithium2 GenerateKey test start\n")
	_, skbytes, pkbytes, _ := dilithium2.GenerateKey()

	msg := []byte("hello")
	signature, _ := dilithium2.Sign(skbytes, msg)

	// msg2 := []byte("hello2")
	err := dilithium2.Verify(pkbytes, msg, signature)
	if !err {
		fmt.Println("MessageErrrdilithium2Sign")
		return
	}
	fmt.Println("end: TestDilithium2 GenerateKey is ok")
}
func TestDilithium3() {
	fmt.Println("start: TestDilithium3 GenerateKey test start\n")
	_, skbytes, pkbytes, _ := dilithium3.GenerateKey()

	msg := []byte("hello")
	signature, _ := dilithium3.Sign(skbytes, msg)

	// msg2 := []byte("hello2")
	err := dilithium3.Verify(pkbytes, msg, signature)
	if !err {
		fmt.Println("MessageErrrdilithium3Sign")
		return
	}
	fmt.Println("end: TestDilithium3 GenerateKey is ok")
}
func TestF512() {
	fmt.Println("start: F512 GenerateKey test start")
	_, skbytes, pkbytes, _ := falcon512.GenerateKey()

	msg := []byte("hello")
	signature, _ := falcon512.Sign(skbytes, msg)

	// msg2 := []byte("hello2")
	err := falcon512.Verify(pkbytes, msg, signature)
	if !err {
		fmt.Println("MessageErrrFalconSign")
		return
	}
	fmt.Println("end: F512 GenerateKey is ok")
}
func TestF1024() {
	fmt.Println("start: F1024 GenerateKey test start")
	_, skbytes, pkbytes, _ := falcon1024.GenerateKey()

	msg := []byte("hello")
	signature, _ := falcon1024.Sign(skbytes, msg)

	// msg2 := []byte("hello2")
	err := falcon1024.Verify(pkbytes, msg, signature)
	if !err {
		fmt.Println("MessageErrrFalconSign")
		return
	}
	fmt.Println("end: F1024 GenerateKey is ok")
}

func TestF512GenBySeed() {
	fmt.Println("start: F512 GenBySeed test ")

	seedbytes, sk1, pk1, _ := falcon512.GenerateKey()
	sk2 := falcon512.GenSkBySeed(seedbytes)
	pk2 := falcon512.GenPkBySeed(seedbytes)
	skr := falcon512.Equals(sk1, sk2)
	if !skr {
		fmt.Println("err: sk1 skr is not Equals ")
	}
	pkr := falcon512.Equals(pk1, pk2)
	if !pkr {
		fmt.Println("err: pk1 pk2 is not Equals ")
	}

	msg := []byte("hello")
	signature1, _ := falcon512.Sign(sk1, msg)

	signature2, _ := falcon512.SignBySeed(seedbytes, msg)

	err := falcon512.Verify(pk1, msg, signature2)
	if !err {
		fmt.Println("MessageErrrFalconSign")
		return
	}
	// msg2 := []byte("hello2")
	err2 := falcon512.Verify(pk2, msg, signature1)
	if !err2 {
		fmt.Println("MessageErrrFalconSign")
		return
	}
	fmt.Println("end: F512GenBySeed is ok")

}

func TestF1024GenBySeed() {
	fmt.Println("start: F1024 GenBySeed test ")

	seedbytes, sk1, pk1, _ := falcon1024.GenerateKey()
	sk2 := falcon1024.GenSkBySeed(seedbytes)
	pk2 := falcon1024.GenPkBySeed(seedbytes)
	skr := falcon1024.Equals(sk1, sk2)
	if !skr {
		fmt.Println("err: sk1 skr is not Equals ")
	}
	pkr := falcon1024.Equals(pk1, pk2)
	if !pkr {
		fmt.Println("err: pk1 pk2 is not Equals ")
	}

	msg := []byte("hello")
	signature1, _ := falcon1024.Sign(sk1, msg)

	signature2, _ := falcon1024.SignBySeed(seedbytes, msg)

	err := falcon1024.Verify(pk1, msg, signature2)
	if !err {
		fmt.Println("MessageErrrFalconSign")
		return
	}
	// msg2 := []byte("hello2")
	err2 := falcon1024.Verify(pk2, msg, signature1)
	if !err2 {
		fmt.Println("MessageErrrFalconSign")
		return
	}
	fmt.Println("end: F1024GenBySeed is ok")

}

func main() {
	TestDilithium2()

	TestDilithium3()

	TestF512()
	TestF512GenBySeed()

	TestF1024()
	TestF1024GenBySeed()
}
