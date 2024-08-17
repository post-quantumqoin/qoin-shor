#include <napi.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>

#include "core/api.h"
#include "core/params.h"

using namespace Napi;

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
    !info[1].IsBuffer() ){
        return Napi::Boolean::New( env, false );
    }
    
    unsigned char *pk = (unsigned char*) info[0].As<Napi::Buffer<unsigned char>>().Data();
    unsigned char *sk = (unsigned char*) info[1].As<Napi::Buffer<unsigned char>>().Data();
    unsigned char *seed = (unsigned char*) info[2].As<Napi::Buffer<unsigned char>>().Data();

    const int genkey = crypto_sign_keypair( pk, sk, seed );
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
    !info[1].IsBuffer() ){
        return Napi::Boolean::New( env, false );
    }
    
    unsigned char *pk = (unsigned char*) info[0].As<Napi::Buffer<unsigned char>>().Data();
    unsigned char *sk = (unsigned char*) info[1].As<Napi::Buffer<unsigned char>>().Data();
    unsigned char *seed = (unsigned char*) info[2].As<Napi::Buffer<unsigned char>>().Data();

    const int genkey = crypto_sign_keypair_by_seed( pk, sk, seed );
    if ( genkey != 0 ){
        return Napi::Boolean::New( env, false );
    } 
    else {
        return Napi::Boolean::New( env, true );
    }
}

// /**
//  * info[0]: buffer of sk
//  * info[1]: buffer of seed
//  */
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
    const int gensk = crypto_generate_sk( sk, seed );
    if ( gensk != 0 ){
        return Napi::Boolean::New( env, false );
    } 
    else {
        return Napi::Boolean::New( env, true );
    }
}

// /**
//  * info[0]: buffer of pk
//  * info[1]: buffer of seed
//  */
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
    const int genpk = crypto_generate_pk( pk,seed );
    if ( genpk != 0 ){
        return Napi::Boolean::New( env, false );
    } 
    else {
        return Napi::Boolean::New( env, true );
    }
}

// /**
//  * info[0]: buffer of signature
//  * info[1]: buffer of message
//  * info[2]: buffer of sk
//  */
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

    const int gensign = crypto_sign_signature( sign_msg, text, text_length, sk, 0 );
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
 * info[3]: flag of random sign
 */
Napi::Boolean NapiSignBySeed( const Napi::CallbackInfo& info )
{
    Napi::Env env = info.Env();
    if ( info.Length() != 3 ) {
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
    unsigned char sk[CRYPTO_SECRETKEYBYTES];
    const int gensk = crypto_generate_sk( sk, seed );
    if( gensk != 0 ){
        return Napi::Boolean::New( env, false );
    }

    const int gensign = crypto_sign_signature( sign_msg, text, text_length, sk, 0 );
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
    if ( info.Length() != 3 ) {
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

    const int verify_ans = crypto_sign_verify( sign, text, text_length, pk );
    if( verify_ans != 0 ){
        return Napi::Boolean::New( env, false );
    } else {
        return Napi::Boolean::New( env, true );
    }
}

Napi::Number getSkLength( const Napi::CallbackInfo& info ) 
{
    return Napi::Number::New(info.Env(), CRYPTO_SECRETKEYBYTES);
}

Napi::Number getPkLength( const Napi::CallbackInfo& info ) 
{
    return Napi::Number::New(info.Env(), CRYPTO_PUBLICKEYBYTES);
}

Napi::Number getSignLength( const Napi::CallbackInfo& info ) 
{
    return Napi::Number::New(info.Env(), CRYPTO_BYTES);
}

Napi::Number getNonceLength( const Napi::CallbackInfo& info ) 
{
    return Napi::Number::New(info.Env(), CRHBYTES);
}

Napi::Number getSeedLength( const Napi::CallbackInfo& info ) 
{
    return Napi::Number::New(info.Env(), SEEDBYTES);
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