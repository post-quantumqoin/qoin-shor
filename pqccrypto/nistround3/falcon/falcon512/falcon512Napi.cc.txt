#include <napi.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>
#include "falconCore/api.h"
using namespace Napi;

#define FALCON512_CRYPTO_SK_SIZE           CRYPTO_SECRETKEYBYTES
#define FALCON512_CRYPTO_PK_SIZE           CRYPTO_PUBLICKEYBYTES
#define FALCON512_CRYPTO_SIGN_SIZE         CRYPTO_BYTES
#define FALCON512_CRYPTO_NONCE_LEN         40
#define FALCON512_CRYPTO_SEED_LEN          48
/**
 * info[0]: buffer of pk
 * info[1]: buffer of sk
 * info[2]: buffer of seed
 */
Napi::Boolean NapiGenkey( const Napi::CallbackInfo& info )
{
    Napi::Env env = info.Env();
    if ( info.Length() != 3 ) {
        return Napi::Boolean::New( env, false );
    }
    else if ( !info[0].IsBuffer() || !info[1].IsBuffer() ||
        !info[2].IsBuffer() ){
        return Napi::Boolean::New( env, false );
    }
    
    unsigned char *pk = (unsigned char*) info[0].As<Napi::Buffer<unsigned char>>().Data();
    unsigned char *sk = (unsigned char*) info[1].As<Napi::Buffer<unsigned char>>().Data();
    unsigned char *seed = (unsigned char*) info[2].As<Napi::Buffer<unsigned char>>().Data();

    const int genkey = falcon_genkey( pk, sk, seed );
    if ( genkey != 0 ){
        return Napi::Boolean::New( env, false );
    } 
    else {
        return Napi::Boolean::New( env, true );
    }
}

/**
 * info[0]: buffer of pk
 * info[1]: buffer of sk
 * info[2]: buffer of seed
 */
Napi::Boolean NapiGenkeyBySeed( const Napi::CallbackInfo& info )
{
    Napi::Env env = info.Env();
    if ( info.Length() != 3 ) {
        return Napi::Boolean::New( env, false );
    }
    else if ( !info[0].IsBuffer() || !info[1].IsBuffer() ||
        !info[2].IsBuffer() ){
        return Napi::Boolean::New( env, false );
    }
    
    unsigned char *pk = (unsigned char*) info[0].As<Napi::Buffer<unsigned char>>().Data();
    unsigned char *sk = (unsigned char*) info[1].As<Napi::Buffer<unsigned char>>().Data();
    unsigned char *seed = (unsigned char*) info[2].As<Napi::Buffer<unsigned char>>().Data();

    const int genkey = falcon_genkey_by_seed( pk, sk, seed );
    if ( genkey != 0 ){
        return Napi::Boolean::New( env, false );
    } 
    else {
        return Napi::Boolean::New( env, true );
    }
}

/**
 * info[0]: buffer of sk
 * info[1]: buffer of seed
 */
Napi::Boolean NapiGenSkBySeed( const Napi::CallbackInfo& info )
{
    Napi::Env env = info.Env();
    if ( info.Length() != 2 ) {
        return Napi::Boolean::New(env, false);
    }
    else if ( !info[0].IsBuffer() || !info[1].IsBuffer() ){
        return Napi::Boolean::New(env, false);
    }

    unsigned char *sk = (unsigned char*) info[0].As<Napi::Buffer<unsigned char>>().Data();
    unsigned char *seed = (unsigned char*) info[1].As<Napi::Buffer<unsigned char>>().Data();
    const int gensk = falcon_seed_to_sk( sk, seed );
    if ( gensk != 0 ){
        return Napi::Boolean::New( env, false );
    } 
    else {
        return Napi::Boolean::New( env, true );
    }
}

/**
 * info[0]: buffer of pk
 * info[1]: buffer of seed
 */
Napi::Boolean NapiGenPkBySeed( const Napi::CallbackInfo& info )
{
    Napi::Env env = info.Env();
    if ( info.Length() != 2 ) {
        return Napi::Boolean::New( env, false );
    }
    else if ( !info[0].IsBuffer() || !info[1].IsBuffer() ){
        return Napi::Boolean::New( env, false );
    }

    unsigned char *pk = (unsigned char*) info[0].As<Napi::Buffer<unsigned char>>().Data();
    unsigned char *seed = (unsigned char*) info[1].As<Napi::Buffer<unsigned char>>().Data();
    const int genpk = falcon_seed_to_pk( pk,seed );
    if ( genpk != 0 ){
        return Napi::Boolean::New( env, false );
    } 
    else {
        return Napi::Boolean::New( env, true );
    }
}

/**
 * info[0]: buffer of signature
 * info[1]: buffer of message
 * info[2]: buffer of sk
 */
Napi::Boolean NapiSign( const Napi::CallbackInfo& info )
{
    Napi::Env env = info.Env();
    if (info.Length() != 3) {
        return Napi::Boolean::New( env, false );
    }
    else if( !info[0].IsBuffer() || !info[1].IsBuffer()
        || !info[2].IsBuffer() ){
        return Napi::Boolean::New( env, false );
    }

    unsigned char* sign_msg = info[0].As<Napi::Buffer<unsigned char >>().Data();
    unsigned char* text = info[1].As<Napi::Buffer<unsigned char >>().Data();
    unsigned long long text_length = info[1].As<Napi::Buffer<unsigned char >>().Length();
    unsigned char* sk = info[2].As<Napi::Buffer<unsigned char>>().Data();

    const int gensign = falcon_sign( sign_msg, text, text_length, sk );
    if ( gensign != 0 ){
        return Napi::Boolean::New( env, false );
    } 
    else {
        return Napi::Boolean::New( env, true );
    }
}

/**
 * info[0]: buffer of signature
 * info[1]: buffer of message
 * info[2]: buffer of seed
 */
Napi::Boolean NapiSignBySeed( const Napi::CallbackInfo& info )
{
    Napi::Env env = info.Env();
    if (info.Length() != 3) {
        return Napi::Boolean::New( env, false );
    }
    else if( !info[0].IsBuffer() || !info[1].IsBuffer() || 
    !info[2].IsBuffer() ){
        return Napi::Boolean::New( env, false );
    }

    unsigned char* sign_msg = info[0].As<Napi::Buffer<unsigned char >>().Data();
    unsigned char* text = info[1].As<Napi::Buffer<unsigned char >>().Data();
    unsigned long long text_length = info[1].As<Napi::Buffer<unsigned char >>().Length();
    unsigned char* seed = info[2].As<Napi::Buffer<unsigned char>>().Data();
    unsigned char sk[FALCON512_CRYPTO_SK_SIZE];
    const int gensk = falcon_seed_to_sk( sk, seed );
    if( gensk != 0 ){
        return Napi::Boolean::New( env, false );
    }

    const int gensign = falcon_sign( sign_msg, text, text_length, sk );
    if( gensign != 0 ){
        return Napi::Boolean::New( env, false );
    } else {
        return Napi::Boolean::New( env, true );
    }
}

/**
 * info[0]: buffer of signature
 * info[1]: buffer of message
 * info[2]: buffer of pk
 */
Napi::Boolean NapiVerifySign( const Napi::CallbackInfo& info )
{
    Napi::Env env = info.Env();
    if (info.Length() != 3) {
        return Napi::Boolean::New( env, false );
    }
    else if( !info[0].IsBuffer() || !info[1].IsBuffer() || 
    !info[2].IsBuffer() ){
        return Napi::Boolean::New( env, false );
    }

    unsigned char* sign = info[0].As<Napi::Buffer<unsigned char >>().Data();
    unsigned char* text = info[1].As<Napi::Buffer<unsigned char >>().Data();
    unsigned long long text_length = info[1].As<Napi::Buffer<unsigned char >>().Length();
    unsigned char* pk = info[2].As<Napi::Buffer<unsigned char>>().Data();

    const int verify_ans = verify_sign( text, text_length, sign, pk );
    if( verify_ans != 0 ){
        return Napi::Boolean::New( env, false );
    } else {
        return Napi::Boolean::New( env, true );
    }
}

Napi::Number getSkLength( const Napi::CallbackInfo& info ) 
{
    return Napi::Number::New(info.Env(), FALCON512_CRYPTO_SK_SIZE);
}

Napi::Number getPkLength( const Napi::CallbackInfo& info ) 
{
    return Napi::Number::New(info.Env(), FALCON512_CRYPTO_PK_SIZE);
}

Napi::Number getSignLength( const Napi::CallbackInfo& info ) 
{
    return Napi::Number::New(info.Env(), FALCON512_CRYPTO_SIGN_SIZE);
}

Napi::Number getNonceLength( const Napi::CallbackInfo& info ) 
{
    return Napi::Number::New(info.Env(), FALCON512_CRYPTO_NONCE_LEN);
}

Napi::Number getSeedLength( const Napi::CallbackInfo& info ) 
{
    return Napi::Number::New(info.Env(), FALCON512_CRYPTO_SEED_LEN);
}

Napi::Object Init(Napi::Env env, Napi::Object exports) 
{
    exports["genkey"] = Function::New(env, NapiGenkey);
    exports["genkeyBySeed"] = Function::New(env, NapiGenkeyBySeed);
    exports["genSkBySeed"] = Function::New(env, NapiGenSkBySeed);
    exports["genPkBySeed"] = Function::New(env, NapiGenPkBySeed);
    exports["sign"] = Function::New(env, NapiSign);
    exports["signBySeed"] = Function::New(env, NapiSignBySeed);
    exports["verifySign"] = Function::New(env, NapiVerifySign);

    exports["getSkLength"] = Function::New(env, getSkLength);
    exports["getPkLength"] = Function::New(env, getPkLength);
    exports["getSignLength"] = Function::New(env, getSignLength);
    exports["getNonceLength"] = Function::New(env, getNonceLength);
    exports["getSeedLength"] = Function::New(env, getSeedLength);

    return exports;
}

NODE_API_MODULE( NODE_GYP_MODULE_NAME, Init );
