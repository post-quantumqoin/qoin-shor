package cgo

/*
#include <stdio.h>
#include<stdlib.h>
#include "sha3.h"
*/
import "C"
import (
	"unsafe"
)

// shake256(data, data.length, temp, outputLength);
func Shake256XOF(input []byte, outputLength int) []byte {
	var context C.sha3_ctx_t
	var output []C.uchar = make([]C.uchar, outputLength)
	inputptr := C.CBytes(input)

	C.shake256_init(&context)
	C.shake_update(&context, inputptr, C.int(len(input)))

	C.shake_xof(&context)
	C.shake_out(&context, &output[0], C.int(len(outputLength)))

	return C.GoBytes(unsafe.Pointer(&output[0]), C.int(outputLength))
}
