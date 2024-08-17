/*
 * Wrapper for implementing the NIST API for the PQC standardization
 * process.
 */

#include <stddef.h>
#include <string.h>
#include <stdio.h>

#include "api.h"
#include "inner.h"
#include "../../../../randombytes/randombytes.h"

/**
 * Generate secret key from seed.
 */
int falcon_seed_to_sk( unsigned char *sk, const unsigned char *seed)
{
	union {
		uint8_t b[FALCON_KEYGEN_TEMP_10];
		uint64_t dummy_u64;
		fpr dummy_fpr;
	} tmp;
	int8_t f[1024], g[1024], F[1024];
	uint16_t h[1024];
	inner_shake256_context rng;
	size_t u, v;
	
	/*
	 * Generate key pair.
	 */
	inner_shake256_init(&rng);
	inner_shake256_inject(&rng, seed, SEEDLEN);
	inner_shake256_flip(&rng);
	Zf(keygen)(&rng, f, g, F, NULL, h, 10, tmp.b);

	/*
	 * Encode private key.
	 */
	sk[0] = 0x50 + 10;
	u = 1;
	v = Zf(trim_i8_encode)(sk + u, CRYPTO_SECRETKEYBYTES - u,
		f, 10, Zf(max_fg_bits)[10]);
	if (v == 0) {
		return -1;
	}
	u += v;
	v = Zf(trim_i8_encode)(sk + u, CRYPTO_SECRETKEYBYTES - u,
		g, 10, Zf(max_fg_bits)[10]);
	if (v == 0) {
		return -1;
	}
	u += v;
	v = Zf(trim_i8_encode)(sk + u, CRYPTO_SECRETKEYBYTES - u,
		F, 10, Zf(max_FG_bits)[10]);
	if (v == 0) {
		return -1;
	}
	u += v;
	if (u != CRYPTO_SECRETKEYBYTES) {
		return -1;
	}

	return 0;
}

/**
 * Generate public key from seed.
 */
int falcon_seed_to_pk(unsigned char *pk, const unsigned char *seed)
{
	union {
		uint8_t b[FALCON_KEYGEN_TEMP_10];
		uint64_t dummy_u64;
		fpr dummy_fpr;
	} tmp;
	int8_t f[1024], g[1024], F[1024];
	uint16_t h[1024];
	inner_shake256_context rng;
	size_t v;

	/*
	 * Generate key pair.
	 */
	inner_shake256_init(&rng);
	inner_shake256_inject(&rng, seed, SEEDLEN);
	inner_shake256_flip(&rng);
	Zf(keygen)(&rng, f, g, F, NULL, h, 10, tmp.b);

	/*
	 * Encode public key.
	 */
	pk[0] = 0x00 + 10;
	v = Zf(modq_encode)(pk + 1, CRYPTO_PUBLICKEYBYTES - 1, h, 10);
	if (v != CRYPTO_PUBLICKEYBYTES - 1) {
		return -1;
	}

	return 0;
}

/**
 * Generate keypair's core.
 */
int
falcon_genkey_core(unsigned char *pk, unsigned char *sk, unsigned char *seed)
{
	union {
		uint8_t b[FALCON_KEYGEN_TEMP_10];
		uint64_t dummy_u64;
		fpr dummy_fpr;
	} tmp;
	int8_t f[1024], g[1024], F[1024];
	uint16_t h[1024];
	inner_shake256_context rng;
	size_t u, v;


	/*
	 * Generate key pair.
	 */
	// randombytes(seed, SEEDLEN);
	inner_shake256_init(&rng);
	inner_shake256_inject(&rng, seed, SEEDLEN);
	inner_shake256_flip(&rng);
	Zf(keygen)(&rng, f, g, F, NULL, h, 10, tmp.b);


	/*
	 * Encode private key.
	 */
	sk[0] = 0x50 + 10;
	u = 1;
	v = Zf(trim_i8_encode)(sk + u, CRYPTO_SECRETKEYBYTES - u,
		f, 10, Zf(max_fg_bits)[10]);
	if (v == 0) {
		return -1;
	}
	u += v;
	v = Zf(trim_i8_encode)(sk + u, CRYPTO_SECRETKEYBYTES - u,
		g, 10, Zf(max_fg_bits)[10]);
	if (v == 0) {
		return -1;
	}
	u += v;
	v = Zf(trim_i8_encode)(sk + u, CRYPTO_SECRETKEYBYTES - u,
		F, 10, Zf(max_FG_bits)[10]);
	if (v == 0) {
		return -1;
	}
	u += v;
	if (u != CRYPTO_SECRETKEYBYTES) {
		return -1;
	}

	/*
	 * Encode public key.
	 */
	pk[0] = 0x00 + 10;
	v = Zf(modq_encode)(pk + 1, CRYPTO_PUBLICKEYBYTES - 1, h, 10);
	if (v != CRYPTO_PUBLICKEYBYTES - 1) {
		return -1;
	}

	return 0;
}

/**
 * Generate keypair by random seed.
 */
int falcon_genkey(unsigned char *pk, unsigned char *sk, unsigned char *seed)
{
	int r = randombytes(seed, SEEDLEN);
	if (r != 0) {
		return -1;
	}
	printf("falcon_genkey_core1024 CRYPTO_SECRETKEYBYTES:%d\n",CRYPTO_PUBLICKEYBYTES);
	fflush(stdout);
	return falcon_genkey_core(pk, sk, seed);
}

/**
 * Generate keypair by seed.
 */
int falcon_genkey_by_seed(unsigned char *pk, unsigned char *sk, unsigned char *seed)
{
	return falcon_genkey_core(pk, sk, seed);
}

/**
 * Signature core
 */
int
falcon_sign_core(unsigned char *sm, const unsigned char *m, 
	unsigned long long mlen, const unsigned char *sk)
{
	union {
		uint8_t b[72 * 1024];
		uint64_t dummy_u64;
		fpr dummy_fpr;
	} tmp;
	int8_t f[1024], g[1024], F[1024], G[1024];
	union {
		int16_t sig[1024];
		uint16_t hm[1024];
	} r;
	unsigned char seed[SEEDLEN];
	// unsigned char esig[CRYPTO_BYTES - 2 - SEEDLEN];
	inner_shake256_context sc;
	size_t u, v, sig_len;

	/*
	 * Decode the private key.
	 */
	if (sk[0] != 0x50 + 10) {
		return -1;
	}
	u = 1;
	v = Zf(trim_i8_decode)(f, 10, Zf(max_fg_bits)[10],
		sk + u, CRYPTO_SECRETKEYBYTES - u);
	if (v == 0) {
		return -1;
	}
	u += v;
	v = Zf(trim_i8_decode)(g, 10, Zf(max_fg_bits)[10],
		sk + u, CRYPTO_SECRETKEYBYTES - u);
	if (v == 0) {
		return -1;
	}
	u += v;
	v = Zf(trim_i8_decode)(F, 10, Zf(max_FG_bits)[10],
		sk + u, CRYPTO_SECRETKEYBYTES - u);
	if (v == 0) {
		return -1;
	}
	u += v;
	if (u != CRYPTO_SECRETKEYBYTES) {
		return -1;
	}
	if (!Zf(complete_private)(G, f, g, F, 10, tmp.b)) {
		return -1;
	}

	/*
	 * Hash message nonce + message into a vector.
	 */
	inner_shake256_init(&sc);
	inner_shake256_inject(&sc, sm + 2, NONCELEN);
	inner_shake256_inject(&sc, m, mlen);
	inner_shake256_flip(&sc);
	Zf(hash_to_point_vartime)(&sc, r.hm, 10);

	/*
	 * Initialize a RNG.
	 */
	int r_status = randombytes(seed, SEEDLEN);
	if (r_status != 0) {
		return -1;
	}
	inner_shake256_init(&sc);
	inner_shake256_inject(&sc, seed, sizeof seed);
	inner_shake256_flip(&sc);

	/*
	 * Compute the signature.
	 */
	Zf(sign_dyn)(r.sig, &sc, f, g, F, G, r.hm, 10, tmp.b);


	/*
	 * Encode the signature and bundle it with the message. Format is:
	 *   signature length     2 bytes
	 *   nonce                40 bytes
	 *   signature            slen bytes
	 */
	sm[2 + NONCELEN] = 0x20 + 10;
	sig_len = Zf(comp_encode)(sm + 3 + NONCELEN, CRYPTO_BYTES - NONCELEN - 3, r.sig, 10);
	if (sig_len == 0) {
		return -1;
	}
	sig_len ++;
	sm[0] = (unsigned char)(sig_len >> 8);
	sm[1] = (unsigned char)sig_len;
	
	return 0;
}

/**
 * Signature
 * Create nonce form randombytes and set to signature
 * message.
 */
int
falcon_sign(unsigned char *sm, const unsigned char *m, 
	unsigned long long mlen, const unsigned char *sk)
{
	int r = randombytes(sm + 2, NONCELEN);
	if (r != 0) {
		return -1;
	}
	return falcon_sign_core(sm,m,mlen,sk);
}

/**
 * Signature
 * Signature message comes with nonce.
 */
int
falcon_sign_custom_nonce(unsigned char *sm, const unsigned char *m, 
	unsigned long long mlen, const unsigned char *sk)
{
	return falcon_sign_core(sm,m,mlen,sk);
}

/**
 * Verify signature with pk
 */
int
verify_sign(const unsigned char *m, unsigned long long mlen,
	const unsigned char *sm, const unsigned char *pk)
{
	union {
		uint8_t b[2 * 1024];
		uint64_t dummy_u64;
		fpr dummy_fpr;
	} tmp;
	const unsigned char *esig;
	uint16_t h[1024], hm[1024];
	int16_t sig[1024];
	inner_shake256_context sc;
	size_t sig_len;

	/*
	 * Decode public key.
	 */
	if (pk[0] != 0x00 + 10) {
		return -1;
	}
	if (Zf(modq_decode)(h, 10, pk + 1, CRYPTO_PUBLICKEYBYTES - 1)
		!= CRYPTO_PUBLICKEYBYTES - 1)
	{
		return -1;
	}
	Zf(to_ntt_monty)(h, 10);

	/*
	 * Find nonce, signature, message length.
	 */
	sig_len = ((size_t)sm[0] << 8) | (size_t)sm[1];

	/*
	 * Decode signature.
	 */
	esig = sm + NONCELEN + 2;
	if (sig_len < 1 || esig[0] != 0x20 + 10) {
		return -1;
	}
	if (Zf(comp_decode)(sig, 10,
		esig + 1, sig_len - 1) != sig_len - 1)
	{
		return -1;
	}

	/*
	 * Hash nonce + message into a vector.
	 */
	inner_shake256_init(&sc);
	inner_shake256_inject(&sc, sm + 2, NONCELEN);
	inner_shake256_inject(&sc, m, mlen);
	inner_shake256_flip(&sc);
	Zf(hash_to_point_vartime)(&sc, hm, 10);

	/*
	 * Verify signature.
	 */
	if (!Zf(verify_raw)(hm, sig, h, 10, tmp.b)) {
		return -1;
	}

	return 0;
}