#include <iostream>
#include <bitset>
#include "./include/common.h"

#define UINT64 unsigned long long
#define UINT32 unsigned int
#define UINT8 unsigned char

#define GET_DEVICE_COUNT 0
#define GET_NUM_OF_EXECUTION 1
#define GET_X 2
#define GET_INFO 3

__device__ UINT64 unity = 1; 
__device__ UINT64 unityL;
__device__ int unuseZiro;


/**
 * Get the number of devices.
 */
int getDeviceCount()
{
    int deviceCount = 0;
    cudaGetDeviceCount(&deviceCount);
    // printf("Device count:%d", deviceCount);
    return deviceCount;
}

/**
 * Get the number of equations.
 */
UINT64 getNumOfExecution(int n, int num_block, int num_Thread)
{
    UINT64 unity = 1;
    UINT64 num_sm = num_block * num_Thread;
    UINT64 tableCount = (UINT64)((((unity << n) - 1) / num_sm) + 1);
    // printf("Total number of executions:%llu", tableCount);
    return tableCount;
}

/**
 * Generation coefficient.
 */
void generateCoef(int n, int m, UINT64 *c, UINT64 *c0, char **EquCharArray)
{
    for (int i = 0; i < m; ++i)
    {
        char *EquChar = EquCharArray[i];
        UINT64 *EQU = (UINT64 *)malloc((n * (n + 1) / 2) * sizeof(UINT64));

        for (int j = 0; j < n * (n + 1) / 2; ++j)
        {
            EQU[j] = EquChar[j] - 48;
        }

        c0[i] = EquChar[(n * (n + 1) / 2)] - 48;
        int index = 0;
        for (int j = 0; j < n; ++j)
        {
            UINT64 temp2 = 0;
            for (int k = j; k < n; ++k)
            {
                temp2 ^= (EQU[index++] << (n - 1 - k));
            }
            c[i * n + j] = temp2;
        }
        free(EQU);
    }
}

/**
 * Generate a table of even numbers from 0 - 255.
 */
void generateByteParity(UINT64 *byte_parity)
{
    for (int i = 0; i < 256; ++i)
    {
        unsigned char temp = i;
        temp ^= temp >> 4;
        temp ^= temp >> 2;
        temp ^= temp >> 1;
        byte_parity[i] = temp & 1;
    }
}

/**
 * Receive equation coefficients in hexadecimal code from external sources.
 * Stop receiving when the 'END' string is received.
 */
void getEquations(int m, char **equations, int coefficientBit, char **eqs)
{   
    //  printf("getEquations:");
    if(eqs == NULL){perror("Error getEquations eqs is null");}
    int numOfBytes = coefficientBit / 8 + (coefficientBit % 8 == 0 ? 0 : 1);
    numOfBytes *= 2;

    char **cinBuf = (char**)malloc(sizeof(char*) * m);
    for (int i = 0; i < m; i++)
    {
        // printf("getEquations numOfBytes:%d\n",numOfBytes);
        // printf("getEquations eqs length:%d\n",strlen(eqs[i]));
        if (numOfBytes != strlen(eqs[i])){perror("Error numOfBytes and eqs is unequal in length");}
        cinBuf[i] = (char*)malloc(sizeof(char) * (numOfBytes + 1));
        // std::cin >> cinBuf[i];
        // strncpy(cinBuf[i], eqs[i], (numOfBytes + 1));
        strcpy(cinBuf[i], eqs[i]);
        // printf("getEquations cinBuf length:%s\n",cinBuf[i]);
        
    }
   
    // char endBuf[32];
    // std::cin >> endBuf;
    // if (strncmp(endBuf, "end", 3))
    // {
    //     fprintf(stderr, "Equations amount wrong!");
    //     exit(EXIT_FAILURE);
    // }
    // for (int i = 0; i < m; i++){
    //      printf("getEquations cinBuf:%x\n",cinBuf[0][i]);
    // }
   
    for (int i = 0; i < m; i++)
    {
        char* bs = hex2bin(cinBuf[i]);
        equations[i] = bs;
        
        free(cinBuf[i]);
    }
    
    free(cinBuf);
}

/**
 * Generate cuda index tables.
 */
void generateWhichXTable(UINT64 num_sm, UINT64 *whichXTable, UINT64 tableCount, int whichXWidth, UINT64 smCount = 0)
{
    whichXTable[0] = smCount * whichXWidth;
    for (UINT64 i = 1; i < num_sm; i++)
    {
        whichXTable[i] = whichXTable[i - 1] + tableCount;
    }
}

/**
 * Initializing cuda devices.
 */
bool Init(cudaDeviceProp *deviceProp, int dev = 0)
{
    // set up device
    CHECK(cudaGetDeviceProperties(deviceProp, dev));
    CHECK(cudaSetDevice(dev));
    return true;
}

/**
 * Start Parallelizing Operations.
 */
__global__ void parallel(int n,
                         int m,
                         UINT64 num_sm,
                         int whichXWidth,
                         UINT64 *whichXTable,
                         UINT64 *__restrict__ c0,
                         UINT64 *__restrict__ c,
                         bool *resultFlg,
                         UINT64 *result)
{
    UINT32 ix = blockIdx.x * blockDim.x + threadIdx.x;
    UINT64 x = whichXTable[ix];
    UINT64 start , end = 0;
    // start = dtime_msec(0);
    start = clock();
    if (*resultFlg)
        return;
    for (int w = 0; w < whichXWidth; w++)
    {
        UINT64 *p;
        p = c;
        int i;
        UINT64 tempX = x << unuseZiro;
       
        // end = clock() - start;
        // if((end / (float)CLOCKS_PER_SEC)>15){
        //      *resultFlg = false;
        //      return;
        // }

        for (i = 0; i < m; ++i)
        {
            UINT64 y = 0;
            UINT64 tmptmpX = tempX;

            end = clock() - start;
            if((end / (float)CLOCKS_PER_SEC)>15){
                    *resultFlg = false;
                    return;
            }

            while (tmptmpX)
            {
                int pos1 = __clzll(tmptmpX);
                tmptmpX ^= unityL >> pos1;
                y ^= (p[pos1] & x);
            }
            if ((__popcll(y) & 1) != c0[i])
                break;
            p += n;
        }
        if (i == m)
        {
            *resultFlg = true;
            *result = x;
            return;
        }
        x++;
    }
    whichXTable[ix] = x;
}

/**
 * Get x and print out json format.
 */

char* getX(int m, int n, int num_block, int num_Thread, int whichXWidth, UINT64 startSMCount, int coefficientBit, char **eqs)
{
    UINT64 num_sm = num_block * num_Thread;

    UINT64 unity = 1;
    UINT64 unityLhost = unity << 63;
    int unuseZiroHost = 64 - n;
    UINT64 *c0, *c, *byte_parity, *result;
    char **equations;
    bool *resultFlg;
    UINT64 *whichXTable;
    UINT64 tableCount = (UINT64)((((unity << n) - 1) / num_sm) + 1);
    char *xstr = (char *)malloc(200 * sizeof(char *));//golang external release
    if(eqs == NULL){perror("Error eqs is null");  return NULL;}

    if (tableCount < whichXWidth)
    {
        whichXWidth = 1;
    }
    // printf("getX:tableCount < whichXWidth)");
    CHECK(cudaMallocManaged(&c, m * n * sizeof(UINT64)));
    CHECK(cudaMallocManaged(&c0, m * sizeof(UINT64)));   
    CHECK(cudaMallocManaged(&byte_parity, 256 * sizeof(UINT64)));
    CHECK(cudaMallocManaged(&whichXTable, num_sm * sizeof(UINT64)));
    CHECK(cudaMallocManaged(&result, sizeof(UINT64)));
    CHECK(cudaMallocManaged(&resultFlg, sizeof(bool)));
    // printf("getX:cudaMallocManaged");
    cudaMemcpyToSymbol(unityL, &unityLhost, sizeof(UINT64), 0, cudaMemcpyHostToDevice);
    cudaMemcpyToSymbol(unuseZiro, &unuseZiroHost, sizeof(int), 0, cudaMemcpyHostToDevice);
   
    *resultFlg = false;
    equations = (char **)malloc(m * sizeof(char *));

    UINT64 start, lastTime = 0, end = 0;
    start = dtime_msec(0);
    // printf("getX:getEquations");
    getEquations(m, equations, coefficientBit, eqs);
    generateCoef(n, m, c, c0, equations);
    generateWhichXTable(num_sm, whichXTable, tableCount, whichXWidth, startSMCount);
    generateByteParity(byte_parity);

    UINT64 solution = 0;
    UINT64 smCount = 0, maxSmCount = ceil(tableCount / (double)whichXWidth);
    start = dtime_msec(0);

    for (smCount = startSMCount; smCount < maxSmCount; smCount++)
    {
        parallel<<<num_block, num_Thread>>>(n, m, num_sm, whichXWidth, whichXTable, c0, c, resultFlg, result);
        CHECK(cudaDeviceSynchronize());
        if (*resultFlg)
        {
            solution = *result;
            break;
        }
        lastTime = end;
        end = dtime_msec(start);
        if((end / (float)CLOCKS_PER_SEC)>15){
             *resultFlg = false;
             break;
        }
    }
    end = dtime_msec(start);

    if (*resultFlg)
    {
        double rate = 0.0;
        UINT64 smUse = ((smCount - startSMCount + 1) * whichXWidth * num_sm);
        if (smCount - startSMCount != 0 && lastTime != 0)
        {
            rate = ((smCount - startSMCount) * whichXWidth * num_sm) / ((float)lastTime / CLOCKS_PER_SEC);
        }

        std::bitset<64> solBitset(solution);
        // strcpy(xstr, solBitset.to_string().c_str());
        sprintf(xstr,"x found:{\"x\":\"%s\", \"gpuTime\":\"%f\", \"rate\":\"%f\", \"smCount\":\"%llu\", \"smUse\":\"%llu\"}\n",
               solBitset.to_string().c_str(),
               end / (float)CLOCKS_PER_SEC,
               rate,
               smCount,
               smUse);
        // printf("x found:{\"x\":\"%s\", \"gpuTime\":\"%f\", \"rate\":\"%f\", \"smCount\":\"%llu\", \"smUse\":\"%llu\"}\n",
        //        solBitset.to_string().c_str(),
        //        end / (float)CLOCKS_PER_SEC,
        //        rate,
        //        smCount,
        //        smUse);
        
		// fflush(stdout);
    }
    else
    {
        printf("x not found\n");
    }
	
    cudaFree(c);
    cudaFree(c0);
    cudaFree(byte_parity);
    cudaFree(whichXTable);
    cudaFree(result);
    cudaFree(resultFlg);
    for (int i = 0; i < m; i++)
		free(equations[i]);
	free(equations);

    return xstr;
}

// int main(int argc, char **argv)
// {
//     int opt = 0, deviceID = 0, m = 1, n = 1, whichXWidth = 1;
//     int coefficientBit = 0;
//     UINT64 startSMCount = 0;
//     for (int i = 1; i < argc; i++)
//     {
//         switch (i)
//         {
//         case 1:
//             opt = atoi(argv[i]); //2 
//             break;
//         case 2:
//             deviceID = atoi(argv[i]);//0 
//             break;
//         case 3:
//             m = atoi(argv[i]);//16 
//             break;
//         case 4:
//             n = atoi(argv[i]);//21 
//             break;
//         case 5:
//             whichXWidth = atoi(argv[i]);//1000 
//             break;
//         case 6:
//             startSMCount = strtoull(argv[i], NULL, 10);//0 
//             break;
//         case 7:
//             coefficientBit = atoi(argv[i]);//232
//             break;
//         default:
//             break;
//         }
//     }

//     if (argc < 2)
//     {
//         printf("Insufficient parameters!\n");
//         return EXIT_FAILURE;
//     }
//     cudaDeviceProp deviceProp;
//     Init(&deviceProp, deviceID);

//     int num_block = deviceProp.multiProcessorCount * 128;
//     int num_Thread = 1024;

//     switch (opt)
//     {
//     case GET_DEVICE_COUNT: 
//         getDeviceCount();
//         break;
//     case GET_NUM_OF_EXECUTION: 
//         getNumOfExecution(n, num_block, num_Thread);
//         break;
//     case GET_X:
//         getX(m, n, num_block, num_Thread, whichXWidth, startSMCount, coefficientBit);
//         break;
//     default:
//         printf("unknow opt!\n");
//         return EXIT_FAILURE;
//     }

//     return EXIT_SUCCESS;
// }

extern "C" {


char* cudaGetX(int deviceID, int m, int n, int whichXWidth, UINT64 startSMCount, int coefficientBit, char **eqs){
    // int deviceID = 0, m = 1, n = 1, whichXWidth = 1;
    // int coefficientBit = 0;
    // UINT64 startSMCount = 0;

    cudaDeviceProp deviceProp;
    Init(&deviceProp, deviceID);
    int num_block = deviceProp.multiProcessorCount * 128;
    int num_Thread = 1024;
    
    // printf("cudaGetX:deviceID:%d,m:%d,n:%d,whichXWidth:%d,startSMCount:%d,coefficientBit:%d,num_block:%d,num_Thread:%d\n",deviceID, m, n, whichXWidth, startSMCount, coefficientBit, num_block, num_Thread);
    
    return getX(m, n, num_block, num_Thread, whichXWidth, startSMCount, coefficientBit, eqs);

}

int cudaGetDevCount() {
    cudaDeviceProp deviceProp;
    Init(&deviceProp, 0);
    return getDeviceCount();
}
// let opts = ['1', '0', `${this.m}`, `${this.n}`];
// getNumOfExecution opt: 1  deviceID: 0 m:this.m n:this.n
UINT64 cudaGetNumOfExecution(int n, int m) {
    cudaDeviceProp deviceProp;
    Init(&deviceProp, 0);

    int num_block = deviceProp.multiProcessorCount * 128;
    int num_Thread = 1024;

    return getNumOfExecution(n, num_block, num_Thread);
}


}