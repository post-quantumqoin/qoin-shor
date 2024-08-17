#include <napi.h>
#include <stdint.h>
#include <stdlib.h>
#include "libs/sha3.h"

using namespace Napi;

/**
 * info[0]: buffer of input
 * info[1]: len of input
 * info[2]: buffer of output
 * info[3]: len of output
 */
Napi::Boolean NapiShake256(const Napi::CallbackInfo& info)
{
	Napi::Env env = info.Env();
	if ( info.Length() != 4 ) {
        return Napi::Boolean::New( env, false );
    }
	else if ( !info[0].IsBuffer() || !info[1].IsNumber() || !info[2].IsBuffer() || !info[3].IsNumber() ){
        return Napi::Boolean::New( env, false );
    }

	uint8_t* input = info[0].As<Napi::Buffer<uint8_t>>().Data();
	size_t inputLen = info[1].As<Napi::Number>().Uint32Value();
	uint8_t* output = info[2].As<Napi::Buffer<unsigned char>>().Data();
	size_t outputLen = info[3].As<Napi::Number>().Uint32Value();

	sha3_ctx_t context;
	shake256_init(&context);
	shake_update(&context, input, inputLen);
	shake_xof(&context);
	shake_out(&context, output, outputLen);

	// EVP_MD_CTX* context = EVP_MD_CTX_create();
    // EVP_DigestInit_ex(context, EVP_shake256(), 0);
    // EVP_DigestUpdate(context, input, inputLen);
	// EVP_DigestFinalXOF(context, output, outputLen);
    // EVP_MD_CTX_destroy(context);
	return Napi::Boolean::New( env, true );
}

Napi::Object Init(Napi::Env env, Napi::Object exports) 
{
    exports["shake256"] = Function::New(env, NapiShake256);
    return exports;
}

NODE_API_MODULE( NODE_GYP_MODULE_NAME, Init );