package shared

/*
#cgo CFLAGS: -I${SRCDIR}/../dilithium/core  -I${SRCDIR}/../../randombytes
#cgo CFLAGS: -Wall -Wextra -Wpedantic -Wredundant-decls -Wshadow -Wvla -Wpointer-arith -O3 -march=native -mtune=native -w 
#cgo LDFLAGS: -lm

#include "fips202.c"
#include "randombytes.c"
*/
import "C"