package mqphash

import (
	"fmt"

	// "github.com/filecoin-project/lotus/pqccrypto/shake3"
	"github.com/filecoin-project/lotus/pqccrypto/shake3"
)

type MQPHash struct {
	Seed                   []byte
	Equations              [][]byte
	Variables              int
	VariablesByte          int
	UnwantedVariablesBit   int
	Coefficient            int
	CoefficientByte        int
	UnwantedCoefficientBit int
	HashBit                int
	HashByte               int
	UnwantedHashBit        int
}

func CreateMQP(seed []byte, equationsN int, variablesN int) *MQPHash {
	var variablesByte, unwantedVariablesBit int

	if variablesN%8 > 0 {
		variablesByte = (variablesN >> 3) + 1
		unwantedVariablesBit = 8 - variablesN%8
	} else {
		variablesByte = variablesN >> 3
		unwantedVariablesBit = 0
	}

	equations := make([][]byte, equationsN)
	coefficient := ((variablesN * (variablesN + 1)) >> 1) + 1
	var coefficientByte, unwantedCoefficientBit int

	if coefficient%8 > 0 {
		coefficientByte = (coefficient >> 3) + 1
		unwantedCoefficientBit = (8 - coefficient%8)
	} else {
		coefficientByte = coefficient >> 3
		unwantedCoefficientBit = 0
	}
	var hashByte, unwantedHashBit int
	// unwantedHashBit := (equationsN % 8) ? 8 - (equationsN % 8) : 0;
	if equationsN%8 > 0 {
		hashByte = (equationsN >> 3) + 1
		unwantedHashBit = (8 - equationsN%8)
	} else {
		hashByte = equationsN >> 3
		unwantedHashBit = 0
	}

	allCoefficient := shake3.Shake256XOF(seed, coefficientByte*equationsN)
	for i := 0; i < equationsN; i++ {
		byteStart := i * coefficientByte
		equations[i] = allCoefficient[byteStart : byteStart+coefficientByte]
		// Discard extra bits
		equations[i][coefficientByte-1] >>= unwantedCoefficientBit
		equations[i][coefficientByte-1] <<= unwantedCoefficientBit
	}
	// return &crypto.Signature{
	mh := &MQPHash{
		Seed:                   seed,
		Equations:              equations,
		Variables:              variablesN,
		VariablesByte:          variablesByte,
		UnwantedVariablesBit:   unwantedVariablesBit,
		Coefficient:            coefficient,
		CoefficientByte:        coefficientByte,
		UnwantedCoefficientBit: unwantedCoefficientBit,
		HashBit:                equationsN,
		HashByte:               hashByte,
		UnwantedHashBit:        unwantedHashBit,
	}
	return mh

}

func (m *MQPHash) xToXx(x []byte) ([]byte, error) {
	// fmt.Println("xToXx start x:", x)
	if len(x) != m.VariablesByte {
		return nil, fmt.Errorf("error input buffer length: ${x.length}, need: ${this._MQP.variablesByte}")
	}

	if x[len(x)-1]&(^(0xff << m.UnwantedVariablesBit)) > 0 {
		return nil, fmt.Errorf("error input bit, last %d bits is not zero", m.UnwantedVariablesBit)
	}
	// fmt.Println("x[len(x)-1]&:", ^(0xff << m.UnwantedVariablesBit))
	xixj := make([]byte, m.CoefficientByte)
	xTemp := make([]byte, m.VariablesByte)
	setIndex := 0

	copy(xTemp[:], x[:])
	// fmt.Println("xToXx CoefficientByte: VariablesByte: xTemp:", m.CoefficientByte, m.VariablesByte, xTemp)
	for i := 0; i < m.Variables; i++ {
		xi := getBufferBit(x, i)
		// fmt.Println(fmt.Sprintf("xToXx getBufferBit xi:%d x:%x i:%d", xi, x, i))
		if xi > 0 {
			setZeroBufferShiftBit(xixj, setIndex, xTemp)
			// fmt.Println(fmt.Sprintf("xToXx getBufferBit xixj:%x setIndex:%d xTemp:%x", xixj, setIndex, xTemp))
		}

		bufferUnshiftOneBit(xTemp)
		setIndex += (m.Variables - i)
		// fmt.Println(fmt.Sprintf("xToXx getBufferBit setIndex:%d xTemp:%x", setIndex, xTemp))
	}

	setBufferBit(xixj, m.Coefficient-1, true) //Constant
	// fmt.Println("xToXx end xixj:", xixj, len(xixj))
	return xixj, nil

}

func (m *MQPHash) Update(x []byte) []byte {
	xixj, err := m.xToXx(x)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	resultL := m.HashByte

	result := make([]byte, resultL) //equations length = hash length

	tempBuf := make([]byte, m.CoefficientByte)
	for m, eq := range m.Equations {
		bufferAnd(eq, xixj, tempBuf)
		sum := bufferXorInside(tempBuf)
		setBufferBit(result, m, sum > 0)
	}
	return result
}

func (m *MQPHash) CheckIsSolution(x []byte) bool {
	xixj, err := m.xToXx(x)
	if err != nil {
		return false
	}

	tempBuf := make([]byte, m.CoefficientByte)

	for _, eq := range m.Equations {
		bufferAnd(eq, xixj, tempBuf)
		sum := bufferXorInside(tempBuf)
		if sum != 0 {
			fmt.Println("mqphash CheckIsSolution is fail:", false)
			return false
		}
	}
	// fmt.Println("mqphash CheckIsSolution is ok:", true)
	return true

}
