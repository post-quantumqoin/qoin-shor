package dilithium5

/*
#cgo CFLAGS: -I${SRCDIR}/../core -I${SRCDIR}/../../../randombytes 
#cgo CFLAGS: -DDILITHIUM_PREFIX=pqcrystals_dilithium5 -DDILITHIUM_MODE=5
#cgo CFLAGS: -Wall -Wextra -Wpedantic -Wno-pragma-messages -Wredundant-decls -Wshadow -Wvla -Wpointer-arith -O3 -march=native -mtune=native
#cgo LDFLAGS: -lm

#include "../core/sign.c"
#include "../core/packing.c"
#include "../core/polyvec.c"
#include "../core/poly.c"
#include "../core/ntt.c"
#include "../core/reduce.c"
#include "../core/rounding.c"
#include "../core/symmetric-shake.c"
#include "../core/fips202.h"
#include "../../../randombytes/randombytes.h"
*/
import "C"