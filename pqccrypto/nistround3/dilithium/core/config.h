#ifndef CONFIG_H
#define CONFIG_H

#include <stdint.h>
#include <stdlib.h>
#include <string.h>

#ifdef __cplusplus
extern "C" {
#endif

#ifndef DILITHIUM_PREFIX
#define DILITHIUM_PREFIX   pqcrystals_dilithium
#pragma message("#ifndef DILITHIUM_PREFIX")
#endif



#define DILITHIUM_NAMESPACE(name)             DILITHIUM_NAMESPACE_(DILITHIUM_PREFIX, name)
#define DILITHIUM_NAMESPACE_(prefix, name)    DILITHIUM_NAMESPACE__(prefix, name)
#define DILITHIUM_NAMESPACE__(prefix, name)   prefix ## _ ## name  

#ifdef __cplusplus
}
#endif

#endif
