package pqcpow

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/post-quantumqoin/qoin-shor/pqccrypto/mqphash"
)

type fix struct {
	mqphash *mqphash.MQPHash
	m       int
	n       int
	newN    int

	newCoe            int
	newCoeByte        int
	unwantedNewCoeBit int
}

func NewFix(mh *mqphash.MQPHash, fixLenght int) *fix {
	f := &fix{}
	f.mqphash = mh
	f.m = mh.HashBit
	f.n = mh.Variables
	f.newN = f.n - fixLenght

	f.newCoe = ((f.newN * (f.newN + 1)) >> 1) + 1
	if f.newCoe%8 > 0 {
		f.newCoeByte = (f.newCoe >> 3) + 1
		f.unwantedNewCoeBit = (8 - f.newCoe%8)
	} else {
		f.newCoeByte = f.newCoe >> 3
		f.unwantedNewCoeBit = 0
	}
	// fmt.Println("NewFix fixLenght: f.m: f.n: f.newN f.newCoe: f.newCoeByte: f.unwantedNewCoeBit:", fixLenght, f.m, f.n, f.newN, f.newCoe, f.newCoeByte, f.unwantedNewCoeBit)
	return f
}

func (f *fix) FixOneEquation(fixString string, coefficientByByte string, unwantedBit int) ([]byte, int, int, int) {
	newN := f.n - len(fixString)
	// fmt.Println("FixOneEquation fixString:", fixString)
	matched, _ := regexp.MatchString("[0-1]", fixString)
	if !matched {
		return nil, 0, 0, 0
	}
	var fix []int
	for _, f := range strings.Split(fixString, "") {
		v, _ := strconv.Atoi(f)
		fix = append(fix, v)
	}
	// fmt.Println("FixOneEquation fixString: fix:", fixString, fix)
	mode := mqphash.BufferBitModeOpt{
		UnwantedBit:        unwantedBit,
		RemoveSpace:        true,
		DisplayUnwantedBit: true,
	}
	data, _ := hex.DecodeString(coefficientByByte)
	// fmt.Println("FixOneEquation fixString:", fixString)
	// fmt.Println("FixOneEquation coefficientByByte data:", data)
	// bufferBitModeString displayUnwantedBit aLine: 0101111001100010101011101000010100111010110111101011....  <unwantedBit>: 00000
	aLine := mqphash.BufferBitModeString(data, mode)
	// fmt.Println("FixOneEquation BufferBitModeString aLine:", fixString, len(aLine), aLine)
	if len(aLine) != ((f.n*(f.n+1))/2)+1 {
		fmt.Println("Err: fixString: len(aLine): aLine: f.n:", fixString, len(aLine), aLine, f.n)
		return nil, 0, 0, 0
	}

	sq := make([]int, 0)
	lin := make([]int, newN)
	var theConst int

	for i := 0; i < f.n; i++ {
		for j := i; j < f.n; j++ {
			val, _ := strconv.Atoi(string(aLine[0]))
			// num, _ := strconv.ParseInt(string(aLine[0]), 2, 8)
			// val := int(num)
			// val := int(aLine[0]) //aLine: 01011110011000101010....
			aLine = aLine[1:]

			if i >= newN {
				if j >= newN {
					theConst ^= val & fix[i-newN] & fix[j-newN]
				} else {
					lin[j] ^= val & fix[i-newN]
				}
			} else {
				if j >= newN {
					lin[i] ^= val & fix[j-newN]
				} else {
					sq = append(sq, val)
				}
			}

		}
	}

	var index int = 0
	for i := 0; i < newN; i++ {
		for j := i; j < newN; j++ {
			if i == j {
				sq[index] ^= lin[i]
			}
			index++
		}
	}
	ta, _ := strconv.Atoi(string(aLine[0]))
	theConst ^= ta
	sq = append(sq, theConst)
	// fmt.Println("FixOneEquation sq:", fixString, len(sq), sq)
	var tmparr []string
	for _, v := range sq { //sq: int 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 0,..
		tmparr = append(tmparr, strconv.Itoa(v))
	}
	// tmparr:"0","0","1"...
	tmp := strings.Join(tmparr, "") //tmp :001...00

	newCoe := ((newN * (newN + 1)) >> 1) + 1
	var newCoeByte int
	var unwantedNewCoeBit int
	if newCoe%8 > 0 {
		newCoeByte = (newCoe >> 3) + 1
		unwantedNewCoeBit = (8 - newCoe%8)
	} else {
		newCoeByte = newCoe >> 3
		unwantedNewCoeBit = 0
	}
	// var tmpd []string
	if newCoe%8 > 0 {
		for i := 0; i < 8-(newCoe%8); i++ {
			// tmpd = append(tmpd, "0")
			tmp += "0"
		}
	}
	// fmt.Println("FixOneEquation tmp:", fixString, len(tmp), tmp)
	// tmparr:"0","0","1"... tmpd:"0","0"
	// tmp := strings.Join(tmparr, "") + strings.Join(tmpd, "") //tmp :001...00

	var tmpHex []string
	for i := 0; i < len(tmp); i += 4 {
		end := i + 4
		t := tmp[i:end]                     //"0010"
		num, _ := strconv.ParseInt(t, 2, 8) // num:2

		tmpHex = append(tmpHex, strconv.FormatInt(num, 16)) // sring:"02"
	}
	// fmt.Println("FixOneEquation tmpHex:", fixString, len(tmpHex), strings.Join(tmpHex, ""))
	newCoeBuf, _ := hex.DecodeString(strings.Join(tmpHex, ""))
	return newCoeBuf, newCoe, newCoeByte, unwantedNewCoeBit

}

func (f *fix) fixBack(x64 string, fixStr string) []byte {
	if len(x64) != 64 {
		return nil
	}
	x := x64[len(x64)-f.newN : len(x64)]
	x = x + fixStr
	// var ed []string
	for index := 0; index < f.mqphash.UnwantedVariablesBit; index++ {
		// ed = append(ed, "0")
		x += "0"
	}
	// x := x64[len(x64)-f.newN:len(x64)] + fixStr + strings.Join(ed, "")

	xBuf := make([]byte, 32)
	index := 0
	for i := 0; i < len(x); i += 8 {
		v, _ := strconv.ParseInt(x[i:i+8], 2, 32)
		xBuf[index] = uint8(v)
		index++
	}

	return xBuf
}
