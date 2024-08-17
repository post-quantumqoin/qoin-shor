#ifndef SIGN_H
#define SIGN_H

#include <stddef.h>
#include <stdint.h>
#include "params.h"
#include "polyvec.h"
#include "poly.h"

#ifdef __cplusplus
extern "C" {
#endif

#define challenge DILITHIUM_NAMESPACE(_challenge)
void challenge(poly *c, const uint8_t seed[SEEDBYTES]);

#define crypto_sign_keypair DILITHIUM_NAMESPACE(_keypair)
int crypto_sign_keypair(uint8_t *pk, uint8_t *sk, uint8_t *seed);

#define crypto_sign_keypair_by_seed DILITHIUM_NAMESPACE(_keypair_by_seed)
int crypto_sign_keypair_by_seed(uint8_t *pk, uint8_t *sk, uint8_t *seed);

#define crypto_generate_pk DILITHIUM_NAMESPACE(_genpk)
int crypto_generate_pk(unsigned char *pk, unsigned char *seed);

#define crypto_generate_sk DILITHIUM_NAMESPACE(_gensk)
int crypto_generate_sk(unsigned char *sk, unsigned char *seed);

#define crypto_sign_signature DILITHIUM_NAMESPACE(_signature)
int crypto_sign_signature(uint8_t *sig, const uint8_t *m, 
                          size_t mlen, const uint8_t *sk,
                          unsigned char random);

#define crypto_sign_verify DILITHIUM_NAMESPACE(_verify)
int crypto_sign_verify(const uint8_t *sig, const uint8_t *m, 
                       size_t mlen, const uint8_t *pk);

#define crypto_sign_keypair_core DILITHIUM_NAMESPACE(_core)
int crypto_sign_keypair_core(uint8_t *pk, uint8_t *sk, uint8_t *seed, uint8_t gensk);

#ifdef __cplusplus
}
#endif

#endif
