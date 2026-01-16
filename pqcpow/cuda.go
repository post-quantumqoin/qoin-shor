package pqcpow

/*
#include <stdlib.h>
#include <stdio.h>

char* cudaGetX(int deviceID, int m, int n, int whichXWidth, unsigned long long int startSMCount, int coefficientBit, char **eqs);
int cudaGetDevCount();
unsigned long long int cudaGetNumOfExecution(int n, int m);
#cgo LDFLAGS:-L${SRCDIR}/kernel -lgpuworker -lcommon  -lstdc++  -L/usr/local/cuda/lib64 -lcudart -lnvidia-ml
*/
import "C"
import (
	"unsafe"
)

// var lk sync.Mutex

func GetDeviceCount() int32 {
	// lk.Lock()
	// defer lk.Unlock()
	return int32(C.cudaGetDevCount())
}

func GetNumOfExecution(n int32, m int32) uint64 {
	// lk.Lock()
	// defer lk.Unlock()
	return uint64(C.cudaGetNumOfExecution(C.int(n), C.int(m)))
}

func CudaGetX(deviceID int, m int, n int, whichXWidth int, startSMCount uint64, coefficientBit int, xIn []string) string {
	// lk.Lock()
	// defer lk.Unlock()
	x := make([]*C.char, 0)

	for _, value := range xIn {
		x = append(x, C.CString(value))
		// fmt.Printf("CudaGetX: index:%d,len:%d,ptr:%x\n", index, len(value), x[index])
	}

	// select {
	// case <-build.Clock.After(4 * time.Second):
	// }
	xstr := C.cudaGetX(C.int(deviceID), C.int(m), C.int(n), C.int(whichXWidth),
		C.ulonglong(startSMCount), C.int(coefficientBit), &x[0])

	gxstr := C.GoString(xstr)

	C.free(unsafe.Pointer(xstr)) //xstr inside cuda malloc, golang external free xstr

	//free C.CBytes malloc
	for _, ptr := range x {
		// fmt.Printf("free C.CBytes malloc:%d,,ptr:%x\n", i, ptr)
		C.free(unsafe.Pointer(ptr))
	}

	return gxstr
	// return ""
}
