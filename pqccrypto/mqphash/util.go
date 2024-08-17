package mqphash

import (
	"strconv"
	"strings"
)

type BufferBitModeOpt struct {
	UnwantedBit        int
	RemoveSpace        bool
	DisplayUnwantedBit bool
}

func getBufferBit(buf []byte, bitIndex int) uint8 {
	byteIndex := bitIndex >> 3
	relativeIndex := bitIndex % 8
	return ((buf[byteIndex] >> (7 - relativeIndex)) & 1)
}

func setBufferBit(buf []byte, bitIndex int, value bool) bool {
	byteIndex := bitIndex >> 3
	relativeIndex := bitIndex % 8

	if byteIndex > len(buf) {
		return false
	}

	if value {
		buf[byteIndex] |= (1 << (7 - relativeIndex))
	} else {
		buf[byteIndex] &= ^(1 << (7 - relativeIndex))
	}

	return true
}

func setBufferShiftBit(buf []byte, shiftBit int, value []byte) bool {
	byteIndex := shiftBit >> 3
	relativeIndex := shiftBit % 8

	if byteIndex > len(buf) {
		return false
	}

	if relativeIndex == 0 {
		buf[byteIndex] = value[0]
	} else {
		bufTemp := make([]byte, 2)
		mask := make([]byte, 2)

		bufTemp[0] = value[0] >> relativeIndex
		bufTemp[1] = value[0] << (8 - relativeIndex)

		mask[1] = (0xff >> relativeIndex)
		mask[0] = ^mask[1]

		buf[byteIndex] = buf[byteIndex]&(bufTemp[0]|mask[0]) | bufTemp[0]

		if buf[byteIndex+1] > 0 {
			buf[byteIndex+1] = buf[byteIndex+1]&(bufTemp[1]|mask[1]) | bufTemp[1]
		}
	}

	return true

}

func setZeroBufferShiftBit(buf []byte, shiftBit int, value []byte) bool {
	byteIndex := shiftBit >> 3
	relativeIndex := shiftBit % 8

	if byteIndex > len(buf) {
		return false
	}
	// fmt.Println(fmt.Sprintf("setZeroBufferShiftBit byteIndex:%d relativeIndex:%d buf:%02x len(buf):%d value:%s len(value):%d", byteIndex, relativeIndex, buf, len(buf), hex.EncodeToString(value), len(value)))
	if relativeIndex == 0 {
		j := 0
		for i := byteIndex; i < len(buf) && j < len(value); i++ {
			buf[i] = value[j]
			j++
		}
		// for (let i = byteIndex, j = 0; i < buf.length && j < value.length; i++, j++) {
		// 	buf[i] = value[j];
		// }

	} else {
		bufTemp := make([]byte, 2)
		j := 0
		// byteIndex: 80      	len(buf):84		len(value):5
		for i := byteIndex; i < len(buf) && j < len(value); i++ {
			bufTemp[0] = value[j] >> relativeIndex //relativeIndex:5
			bufTemp[1] = value[j] << (8 - relativeIndex)
			buf[i] |= bufTemp[0]
			if (i + 1) < len(buf) {
				buf[i+1] |= bufTemp[1]
				// fmt.Println("setZeroBufferShiftBit buf[i+1]:", buf[i+1])
			}
			// fmt.Println("setZeroBufferShiftBit bufTemp[0]: bufTemp[1]: buf[i]:", bufTemp[0], bufTemp[1], buf[i])
			j++
		}

	}
	// fmt.Println(fmt.Sprintf("setZeroBufferShiftBit byteIndex:%d relativeIndex:%d buf:%02x len(buf):%d value:%s len(value):%d", byteIndex, relativeIndex, buf, len(buf), hex.EncodeToString(value), len(value)))
	return true
}

func setBufferFillZero(buf []byte, shiftBit int) bool {
	byteIndex := shiftBit >> 3
	relativeIndex := shiftBit % 8

	if byteIndex > len(buf) {
		return false
	}

	buf[byteIndex] = buf[byteIndex] >> relativeIndex
	buf[byteIndex] = buf[byteIndex] << relativeIndex

	for i := byteIndex; i < len(buf); i++ {
		buf[i] = 0x00
	}
	return true
}

func bufferShiftOneBit(buf []byte) {
	var lastBit uint8 = 0
	for i := 0; i < len(buf); i++ {
		if lastBit > 0 {
			lastBit = buf[i] & 0x01
			buf[i] >>= 1
			buf[i] |= 0x80
		} else {
			lastBit = buf[i] & 0x01
			buf[i] >>= 1
		}
	}
}

func bufferUnshiftOneBit(buf []byte) {
	var lastBit uint8 = 0
	for i := len(buf) - 1; i >= 0; i-- {
		if lastBit > 0 {
			lastBit = buf[i] & 0x80
			buf[i] <<= 1
			buf[i] |= 0x01
		} else {
			lastBit = buf[i] & 0x80
			buf[i] <<= 1
		}
	}
}

func bufferAnd(a, b, c []byte) {
	for i := 0; i < len(a); i++ {
		c[i] = a[i] & b[i]
	}
}

func bufferXor(a, b, c []byte) {
	for i := 0; i < len(a); i++ {
		c[i] = a[i] ^ b[i]
	}
}

func bufferXorInside(buf []byte) uint8 {
	var temp uint8 = 0x00

	for i := 0; i < len(buf); i++ {
		temp ^= buf[i]
	}

	temp = temp ^ (temp >> 4)
	temp = temp ^ (temp >> 2)
	temp = temp ^ (temp >> 1)

	return (temp & 1)
}

func padStart(str string, length int, fillChar string) string {
	for {
		if len(str) >= length {
			break
		}
		str = fillChar + str
	}
	return str
}

func BufferBitModeString(buf []byte, opt BufferBitModeOpt) string {
	// fmt.Println("BufferBitModeString buf: opt:", buf, opt)
	var tempArr []string
	for i := 0; i < len(buf); i++ {
		tempArr = append(tempArr, padStart(strconv.FormatInt(int64(buf[i]), 2), 8, "0"))
	}
	// fmt.Println("BufferBitModeString tempArr: ", tempArr)
	var str string
	// reg := regexp.MustCompile(`.{1,8}`)
	// tempArr := reg.FindAllString(string(temp), -1)
	if opt.RemoveSpace {
		str = strings.Join(tempArr, "")
	} else {
		str = strings.Join(tempArr, " ")
	}
	// fmt.Println("BufferBitModeString str: ", str)
	var rt string
	if opt.UnwantedBit > 0 {
		stringIndex := ((len(buf) << 3) - opt.UnwantedBit)
		if !opt.RemoveSpace {
			stringIndex += (stringIndex >> 3)
		}
		if opt.DisplayUnwantedBit {
			rt = str[0:stringIndex]
		} else {
			rt = str[0:stringIndex] + "<unwantedBit>:" + str[stringIndex:]
			// temp = `${temp.slice(0, stringIndex)} <unwantedBit>: ${temp.slice(stringIndex)}`
		}
	} else {
		return str
	}

	return rt
}

func checkBufferIsZero(buf []byte) bool {
	for i := 0; i < len(buf); i++ {
		if buf[i] != 0 {
			return false
		}
	}
	return true
}
