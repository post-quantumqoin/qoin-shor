extend=./common.cpp
# detect if running under unix by finding 'rm' in $PATH :
ifeq ($(wildcard $(addsuffix /rm,$(subst :, ,$(PATH)))),)
WINMODE=1
else
WINMODE=0
endif

ifeq ($(WINMODE),1)
# native windows setup :
LIBFILE=libgpuworker.dll
OUTPUTFILE=mineWin32.exe
else
# cross-compile setup :
LIBFILE=libgpuworker.so
OUTPUTFILE=mine
LDFLAGS := -lnvidia-ml -L /usr/lib/x86_64-linux-gnu 
endif

AR = ar
ARFLAGS = rcs

NVCFLAGS = -O2 --std=c++14 --resource-usage -lineinfo #-g
CUDA_ARCH =  -arch=sm_50 -gencode=arch=compute_50,code=sm_50 -gencode=arch=compute_52,code=sm_52 -gencode=arch=compute_60,code=sm_60 \
-gencode=arch=compute_61,code=sm_61 -gencode=arch=compute_70,code=sm_70 \
-gencode=arch=compute_75,code=sm_75 -gencode=arch=compute_75,code=compute_75
NVCFLAGS += $(CUDA_ARCH)
LIBCFLAGS = --ptxas-options=-v --compiler-options '-fPIC' --shared
LIBCFLAGS += $(NVCFLAGS)

OBJ1 =  common.o
HEAD1 = ./common.h

build: gpu_worker.cu
	nvcc $(NVCFLAGS) $< $(extend) $(LDFLAGS) -o ./$(OUTPUTFILE)

runCU: build
	cat ./test_case_10_16_21_232.txt | ./$(OUTPUTFILE) 2 0 16 21 1000 0 232

buildlib:gpu_worker.cu
	nvcc $(LIBCFLAGS) $< $(extend) $(LDFLAGS) -o ./$(LIBFILE)

common.o: common.cpp  
	g++ -O2 --std=c++14  -lineinfo -c common.cpp -o common.o 

common.a: $(OBJ1) 
	ar rcs libcommon.a common.o

gpu_worker.a: gpu_worker.cu
	nvcc $(NVCFLAGS) --lib  gpu_worker.cu  -o libgpuworker.a
# $(AR) $(ARFLAGS) libgpuworkera.a $(OBJ1) 

