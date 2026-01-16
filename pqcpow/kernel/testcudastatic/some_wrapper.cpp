#include <stdio.h>
#include <stdlib.h>
#include <cstdlib>
#include <time.h>
#include <string.h>
#include <iostream>
#ifdef __cplusplus
extern "C" {
#endif

unsigned long long dtime_msec(unsigned long long start)
{
    clock_t t = clock();
    return t - start;
}

void *safe_malloc(size_t size)
{
    void *ptr = NULL;
    if (NULL == (ptr = malloc(size)))
    {
        printf("err: insufficient memory\n");
    }
    return ptr;
}

int hex2int(char ch)
{
    if (ch >= '0' && ch <= '9')
        return ch - '0';
    if (ch >= 'A' && ch <= 'F')
        return ch - 'A' + 10;
    if (ch >= 'a' && ch <= 'f')
        return ch - 'a' + 10;
    return -1;
}

char *hex2bin(char *s)
{
    unsigned int i, e, l = 0, L = strlen(s);
    char *bin = (char *)malloc(sizeof(char) * L * 4 + 1);
    for (i = 0; i < L; i += 1)
    {
        sscanf(s + i, "%01x", &e);
        for (int j = 8; j > 0; j/=2)
        {
            bin[l++] = e / j == 1 ? '1' : '0';
            e %= j;
        }
    }
    bin[l] = 0;
    return bin;
}
#ifdef __cplusplus
}
#endif