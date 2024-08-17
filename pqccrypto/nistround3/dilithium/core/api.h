#ifndef API_H
#define API_H

#include "params.h"
#ifdef __cplusplus
extern "C" {
#endif

#define crypto_sign_keypair DILITHIUM_NAMESPACE(_keypair)
int crypto_sign_keypair(unsigned char *pk, unsigned char *sk, unsigned char *seed);

#define crypto_sign_keypair_by_seed DILITHIUM_NAMESPACE(_keypair_by_seed)
int crypto_sign_keypair_by_seed(unsigned char *pk, unsigned char *sk, unsigned char *seed);

#define crypto_generate_pk DILITHIUM_NAMESPACE(_genpk)
int crypto_generate_pk(unsigned char *pk, unsigned char *seed);

#define crypto_generate_sk DILITHIUM_NAMESPACE(_gensk)
int crypto_generate_sk(unsigned char *sk, unsigned char *seed);

#define crypto_sign_signature DILITHIUM_NAMESPACE(_signature)
int crypto_sign_signature(unsigned char *sm, const unsigned char *msg, 
                          unsigned long long len, const unsigned char *sk,
                          unsigned char random);

#define crypto_sign_verify DILITHIUM_NAMESPACE(_verify)
int crypto_sign_verify(unsigned char *m, const unsigned char *sm, 
                       unsigned long long smlen, const unsigned char *pk);

#ifdef __cplusplus
}
#endif

#endif
