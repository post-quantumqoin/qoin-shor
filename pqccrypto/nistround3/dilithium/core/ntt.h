#ifndef NTT_H
#define NTT_H

#include <stdint.h>
#include "params.h"

#ifdef __cplusplus
extern "C" {
#endif

#define ntt DILITHIUM_NAMESPACE(_ntt)
void ntt(int32_t a[N]);

#define invntt_tomont DILITHIUM_NAMESPACE(_invntt_tomont)
void invntt_tomont(int32_t a[N]);

#ifdef __cplusplus
}
#endif

#endif
