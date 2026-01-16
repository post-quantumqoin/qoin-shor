#ifndef COMMON_H
#define COMMON_H

#include <stdio.h>
#include <stdlib.h>
#include <cstdlib>
#include <time.h>
#include <string.h>
#ifdef __cplusplus
extern "C" {
#endif
#define SMALLOC(type, num) \
    ((type*) safe_malloc((num) * sizeof(type)))

#define CHECK(call)                                                            \
{                                                                              \
    const cudaError_t error = call;                                            \
    if (error != cudaSuccess)                                                  \
    {                                                                          \
        fprintf(stderr, "Error: %s:%d, ", __FILE__, __LINE__);                 \
        fprintf(stderr, "code: %d, reason: %s\n", error,                       \
                cudaGetErrorString(error));                                    \
        exit(1);                                                               \
    }                                                                          \
}

unsigned long long dtime_msec(unsigned long long start);
void* safe_malloc(size_t size);
int hex2int(char ch);
char *hex2bin(char *s);
#ifdef __cplusplus
}
#endif
#endif /* COMMON_H */