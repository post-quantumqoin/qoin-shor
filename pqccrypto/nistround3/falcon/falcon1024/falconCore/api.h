#define CRYPTO_SECRETKEYBYTES   2305
#define CRYPTO_PUBLICKEYBYTES   1793
#define CRYPTO_BYTES            1332
#define CRYPTO_ALGNAME          "Falcon-1024"
#define NONCELEN                40
#define SEEDLEN                 48

#include "inner.h"
#ifdef __cplusplus
extern "C" {
#endif

#define falcon_seed_to_sk Zf(_seed_to_sk)
int falcon_seed_to_sk(unsigned char *sk, const unsigned char *seed);

#define falcon_seed_to_pk Zf(_seed_to_pk)
int falcon_seed_to_pk(unsigned char *pk, const unsigned char *seed);

#define falcon_genkey Zf(_falcon_genkey)
int falcon_genkey(unsigned char *pk, unsigned char *sk, 
	unsigned char *seed);

#define falcon_genkey_by_seed Zf(_genkey_by_seed)
int falcon_genkey_by_seed(unsigned char *pk, unsigned char *sk, unsigned char *seed);

#define falcon_sign Zf(_falcon_sign)
int falcon_sign(unsigned char *sm,
	const unsigned char *m, unsigned long long mlen,
	const unsigned char *sk);

#define falcon_sign_custom_nonce Zf(_custom_nonce)
int falcon_sign_custom_nonce(unsigned char *sm, const unsigned char *m, 
	unsigned long long mlen, const unsigned char *sk);

#define verify_sign Zf(_sign)
int verify_sign(const unsigned char *m, unsigned long long mlen,
	const unsigned char *sm, const unsigned char *pk);

#define falcon_sign_core Zf(_sign_core)
int falcon_sign_core(unsigned char *sm, const unsigned char *m, 
	unsigned long long mlen, const unsigned char *sk);
	
#define falcon_genkey_core Zf(_genkey_core)
int falcon_genkey_core(unsigned char *pk, unsigned char *sk, unsigned char *seed);
#ifdef __cplusplus
}
#endif