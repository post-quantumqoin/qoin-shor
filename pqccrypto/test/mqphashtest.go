package main

import (
	// "bytes"
	// "io"
	"fmt"

	"pqccrypto/shake3"

	"pqccrypto/mqphash"
)

func main() {
	msg := []byte("hello")
	sh := shake3.Shake256XOF(msg, 72)
	fmt.Printf("Hello hash:%x\n", sh)
	fmt.Println("Hello len:", len(sh))

	m := 16
	n := m + 5
	mh := mqphash.CreateMQP(sh, m, n)

	// Seed                   []byte
	// Equations              [][]byte
	// Variables              int
	// VariablesByte          int
	// UnwantedVariablesBit   int
	// Coefficient            int
	// CoefficientByte        int
	// UnwantedCoefficientBit int
	// HashBit                int
	// HashByte               int
	// UnwantedHashBit        int

	fmt.Printf("mqphash seed:%x\n", mh.Seed)

	for i, eq := range mh.Equations {
		fmt.Printf("mqphash Equations:%x len:%d  index:%d \n", eq, len(eq), i)
	}

	fmt.Println("mqphash Variables:", mh.Variables)
	fmt.Println("mqphash variablesByte:", mh.VariablesByte)
	fmt.Println("mqphash UnwantedVariablesBit:", mh.UnwantedVariablesBit)
	fmt.Println("mqphash Coefficient:", mh.Coefficient)
	fmt.Println("mqphash CoefficientByte:", mh.CoefficientByte)
	fmt.Println("mqphash UnwantedCoefficientBit:", mh.UnwantedCoefficientBit)
	fmt.Println("mqphash HashBit:", mh.HashBit)
	fmt.Println("mqphash HashByte:", mh.HashByte)
	fmt.Println("mqphash UnwantedHashBit:", mh.UnwantedHashBit)

	x2 := []byte("123")
	x2h := mh.Update(x2)
	fmt.Println("mqphash update x2h:", x2h)
}
