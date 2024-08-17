package main
 

import "C"
import (
	// "bytes"
	// "io"
	"fmt"

	"pqccrypto/nistRound3/dilithium/dilithium2"
	// "pqccrypto/nistRound3/dilithium/dilithium3"
	// "pqccrypto/nistRound3/dilithium/dilithium5"
	// "pqccrypto/nistRound3/falcon/falcon1024"
	
)
func TestDilithium2()  {
	fmt.Println("start: TestDilithium2 GenerateKey test start\n")
	_, skbytes, pkbytes, _:= dilithium2.GenerateKey()

	msg := []byte("hello")
	signature, _:= dilithium2.Sign(skbytes, msg)

	// msg2 := []byte("hello2")
	err := dilithium2.Verify(pkbytes, msg, signature)
	if !err {
		fmt.Println("MessageErrrdilithium2Sign")
		return
	}
	fmt.Println("end: TestDilithium2 GenerateKey is ok")
}

// func TestDilithium3()  {
// 	fmt.Println("start: TestDilithium3 GenerateKey test start\n")
// 	_, skbytes, pkbytes, _:= dilithium3.GenerateKey()

// 	msg := []byte("hello")
// 	signature, _:= dilithium3.Sign(skbytes, msg)

// 	// msg2 := []byte("hello2")
// 	err := dilithium3.Verify(pkbytes, msg, signature)
// 	if !err {
// 		fmt.Println("MessageErrrdilithium3Sign")
// 		return
// 	}
// 	fmt.Println("end: TestDilithium3 GenerateKey is ok")
// }

// func TestDilithium5()  {
// 	fmt.Println("start: TestDilithium5 GenerateKey test start\n")
// 	_, skbytes, pkbytes, _:= dilithium5.GenerateKey()

// 	msg := []byte("hello")
// 	signature, _:= dilithium5.Sign(skbytes, msg)

// 	// msg2 := []byte("hello2")
// 	err := dilithium5.Verify(pkbytes, msg, signature)
// 	if !err {
// 		fmt.Println("MessageErrrdilithium3Sign")
// 		return
// 	}
// 	fmt.Println("end: TestDilithium5 GenerateKey is ok")
// }

func TestDilithium2GenBySeed() {
	fmt.Println("start: TestDilithium2GenBySeed GenBySeed test \n")

	seedbytes, sk1, pk1, _:= dilithium2.GenerateKey()
	sk2 := dilithium2.GenSkBySeed(seedbytes)
	pk2 := dilithium2.GenPkBySeed(seedbytes)
	skr := dilithium2.Equals(sk1,sk2)
	if !skr {
		fmt.Println("err: sk1 skr is not Equals ")
	}
	pkr := dilithium2.Equals(pk1,pk2)
	if !pkr {
		fmt.Println("err: pk1 pk2 is not Equals ")
	}
	
	msg := []byte("hello")
	signature1, _:= dilithium2.Sign(sk1, msg)

	signature2, _:= dilithium2.SignBySeed(seedbytes, msg)

	err := dilithium2.Verify(pk1, msg, signature2)
	if !err {
		fmt.Println("MessageErrrFalconSign")
		return
	}
	// msg2 := []byte("hello2")
	err2 := dilithium2.Verify(pk2, msg, signature1)
	if !err2 {
		fmt.Println("MessageErrrDilithium2Sign")
		return
	}
	fmt.Println("end: TestDilithium2GenBySeed is ok")

}

// func TestDilithium3GenBySeed() {
// 	fmt.Println("start: TestDilithium3GenBySeed GenBySeed test \n")

// 	seedbytes, sk1, pk1, _:= dilithium3.GenerateKey()
// 	sk2 := dilithium3.GenSkBySeed(seedbytes)
// 	pk2 := dilithium3.GenPkBySeed(seedbytes)
// 	skr := dilithium3.Equals(sk1,sk2)
// 	if !skr {
// 		fmt.Println("err: sk1 skr is not Equals ")
// 	}
// 	pkr := dilithium3.Equals(pk1,pk2)
// 	if !pkr {
// 		fmt.Println("err: pk1 pk2 is not Equals ")
// 	}
	
// 	msg := []byte("hello")
// 	signature1, _:= dilithium3.Sign(sk1, msg)

// 	signature2, _:= dilithium3.SignBySeed(seedbytes, msg)

// 	err := dilithium3.Verify(pk1, msg, signature2)
// 	if !err {
// 		fmt.Println("MessageErrrFalconSign")
// 		return
// 	}
// 	// msg2 := []byte("hello2")
// 	err2 := dilithium3.Verify(pk2, msg, signature1)
// 	if !err2 {
// 		fmt.Println("MessageErrrDilithium2Sign")
// 		return
// 	}
// 	fmt.Println("end: TestDilithium3GenBySeed is ok")

// }

// func TestDilithium5GenBySeed() {
// 	fmt.Println("start: TestDilithium5GenBySeed GenBySeed test \n")

// 	seedbytes, sk1, pk1, _:= dilithium5.GenerateKey()
// 	sk2 := dilithium5.GenSkBySeed(seedbytes)
// 	pk2 := dilithium5.GenPkBySeed(seedbytes)
// 	skr := dilithium5.Equals(sk1,sk2)
// 	if !skr {
// 		fmt.Println("err: sk1 skr is not Equals ")
// 	}
// 	pkr := dilithium5.Equals(pk1,pk2)
// 	if !pkr {
// 		fmt.Println("err: pk1 pk2 is not Equals ")
// 	}
	
// 	msg := []byte("hello")
// 	signature1, _:= dilithium5.Sign(sk1, msg)

// 	signature2, _:= dilithium5.SignBySeed(seedbytes, msg)

// 	err := dilithium5.Verify(pk1, msg, signature2)
// 	if !err {
// 		fmt.Println("MessageErrrFalconSign")
// 		return
// 	}
// 	// msg2 := []byte("hello2")
// 	err2 := dilithium5.Verify(pk2, msg, signature1)
// 	if !err2 {
// 		fmt.Println("MessageErrrDilithium5Sign")
// 		return
// 	}
// 	fmt.Println("end: TestDilithium5GenBySeed is ok")

// }

// }
 
func main() {
 
	TestDilithium2() 
	TestDilithium2GenBySeed() 

	// TestDilithium3()
	// TestDilithium3GenBySeed()

	// TestDilithium5()  
	// TestDilithium5GenBySeed()

}