package shake3

/*
#cgo CFLAGS: -I./libs
#cgo LDFLAGS: -L./libs

#include<stdlib.h>
#include "sha3.h"
#include "sha3.c"
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

	// C.shake256_init(&context) #define shake256_init(c) sha3_init(c, 32)
	C.sha3_init(&context, C.int(32))
	// int sha3_init(sha3_ctx_t *c, int mdlen);
	C.sha3_update(&context, inputptr, C.size_t(len(input)))
	// int sha3_update(sha3_ctx_t *c, const void *data, size_t len);
	C.shake_xof(&context)
	C.shake_out(&context, unsafe.Pointer(&output[0]), C.size_t(outputLength))
	// void shake_out(sha3_ctx_t *c, void *out, size_t len);
	return C.GoBytes(unsafe.Pointer(&output[0]), C.int(outputLength))
}
