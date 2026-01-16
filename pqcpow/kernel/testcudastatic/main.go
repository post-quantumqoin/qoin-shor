package main

/*
#cgo LDFLAGS: -L${SRCDIR} -lsomewrapper -lstdc++
void vectorAdd(float* A, float* B, float* C, int N);
 int getDeviceCount();
 unsigned long long dtime_msec(unsigned long long start);
*/
import "C"
import "fmt"

func main() {
	// #cgo LDFLAGS:-L. -lstdc++  -L./ -lvectorAdd -lcommon -L/usr/local/cuda-12.4/lib64 -lcudart
	// N := 10
	// A := make([]float32, N)
	// B := make([]float32, N)
	// C := make([]float32, N)
	// for i := range A {
	// 	A[i] = 1.0
	// 	B[i] = 2.0
	// }

	// // Convert slices to pointers for C calls.
	// c_A := (*C.float)(&A[0])
	// c_B := (*C.float)(&B[0])
	// c_C := (*C.float)(&C[0])

	// C.vectorAdd(c_A, c_B, c_C, C.int(N))

	fmt.Println("dtime_msec", C.dtime_msec(C.ulonglong(0)))
	// fmt.Println("getDeviceCount", C.getDeviceCount())
	// C is now [3,3,3,3,3,3,3,3,3,3]
}
