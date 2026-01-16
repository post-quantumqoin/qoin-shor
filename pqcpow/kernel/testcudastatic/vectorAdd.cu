#include <cuda_runtime.h>
#include <iostream>
__global__ void vectorAddKernel(float* A, float* B, float* C, int N) {
    int i = threadIdx.x + blockIdx.x * blockDim.x;
    if (i < N) {
        C[i] = A[i] + B[i];
    }
}
 
extern "C" {
    int getDeviceCount()
    {
        int deviceCount = 0;
        cudaGetDeviceCount(&deviceCount);
        printf("Device count:%d", deviceCount);
        return deviceCount;
    }
    void vectorAdd(float* A, float* B, float* C, int N) {
        float* devA;
        float* devB;
        float* devC;
        cudaMalloc((void**)&devA, N * sizeof(float));
        cudaMalloc((void**)&devB, N * sizeof(float));
        cudaMalloc((void**)&devC, N * sizeof(float));
        cudaMemcpy(devA, A, N * sizeof(float), cudaMemcpyHostToDevice);
        cudaMemcpy(devB, B, N * sizeof(float), cudaMemcpyHostToDevice);
 
        int blockSize = 256;
        int numBlocks = (N + blockSize - 1) / blockSize;
        vectorAddKernel<<<numBlocks, blockSize>>>(devA, devB, devC, N);
 
        cudaMemcpy(C, devC, N * sizeof(float), cudaMemcpyDeviceToHost);
        cudaFree(devA);
        cudaFree(devB);
        cudaFree(devC);
    }
}