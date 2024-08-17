CC ?= /usr/bin/cc
AR = ar
CFLAGS += -Wall -Wextra -Wpedantic -Wmissing-prototypes -Wredundant-decls \
  -Wshadow -Wvla -Wpointer-arith -O3 -march=native -mtune=native
NISTFLAGS += -Wno-unused-result -O3
SOURCES = sign.c packing.c polyvec.c poly.c ntt.c reduce.c rounding.c

HEAD1 = config.h params.h api.h sign.h packing.h polyvec.h poly.h ntt.h \
  reduce.h rounding.h symmetric.h ../../../randombytes/randombytes.h fips202.h

HEADERS = config.h params.h api.h sign.h packing.h polyvec.h poly.h ntt.h \
  reduce.h rounding.h symmetric.h ../../../randombytes/randombytes.h 
KECCAK_SOURCES = $(SOURCES) fips202.c symmetric-shake.c
KECCAK_HEADERS = $(HEADERS) fips202.h
AES_SOURCES = $(SOURCES) fips202.c aes256ctr.c symmetric-aes.c
AES_HEADERS = $(HEADERS) fips202.h aes256ctr.h

OBJ2 = build2/sign.o build2/packing.o build2/polyvec.o build2/poly.o build2/ntt.o \
	build2/reduce.o  build2/rounding.o build2/randombytes.o build2/symmetric-shake.o build2/fips202.o

OBJ3 = build3/sign.o build3/packing.o build3/polyvec.o build3/poly.o build3/ntt.o \
	build3/reduce.o  build3/rounding.o build3/randombytes.o build3/symmetric-shake.o build3/fips202.o

OBJ5 = build5/sign.o build5/packing.o build5/polyvec.o build5/poly.o build5/ntt.o \
	build5/reduce.o  build5/rounding.o build5/randombytes.o build5/symmetric-shake.o build5/fips202.o

.PHONY:  all speed shared  clean


all: PQCgenKAT_sign

speed: \
  test/test_mul \
  test/test_speed2 \
  test/test_speed3 \
  test/test_speed5 \
  test/test_speed2aes \
  test/test_speed3aes \
  test/test_speed5aes

shared: \
  build/libpqcrystals_dilithium2_ref.so \
  build/libpqcrystals_dilithium3_ref.so \
  build/libpqcrystals_dilithium5_ref.so \
  build/libpqcrystals_fips202_ref.so \
  #libpqcrystals_dilithium2aes_ref.so \
  #libpqcrystals_dilithium3aes_ref.so \
  #libpqcrystals_dilithium5aes_ref.so \

  #libpqcrystals_aes256ctr_ref.so

static: \
  build2/libpqcrystals_dilithium2_ref.a \
  build3/libpqcrystals_dilithium3_ref.a \
  build5/libpqcrystals_dilithium5_ref.a \

build/libpqcrystals_fips202_ref.so: fips202.c fips202.h
	$(CC) -shared -fPIC $(CFLAGS) -o $@ $<

libpqcrystals_aes256ctr_ref.so: aes256ctr.c aes256ctr.h
	$(CC) -shared -fPIC $(CFLAGS) -o $@ $<

build/libpqcrystals_dilithium2_ref.so: $(SOURCES) $(HEADERS) symmetric-shake.c
	$(CC) -shared -fPIC $(CFLAGS) -DDILITHIUM_MODE=2 \
	  -o $@ $(SOURCES) symmetric-shake.c ../../../randombytes/randombytes.c

build/libpqcrystals_dilithium3_ref.so: $(SOURCES) $(HEADERS) symmetric-shake.c
	$(CC) -shared -fPIC $(CFLAGS) -DDILITHIUM_MODE=3 \
	  -o $@ $(SOURCES) symmetric-shake.c ../../../randombytes/randombytes.c

build/libpqcrystals_dilithium5_ref.so: $(SOURCES) $(HEADERS) symmetric-shake.c
	$(CC) -shared -fPIC $(CFLAGS) -DDILITHIUM_MODE=5 \
	  -o $@ $(SOURCES) symmetric-shake.c ../../../randombytes/randombytes.c

libpqcrystals_dilithium2aes_ref.so: $(SOURCES) $(HEADERS) symmetric-aes.c
	$(CC) -shared -fPIC $(CFLAGS) -DDILITHIUM_MODE=2 -DDILITHIUM_USE_AES \
	   -o $@ $(SOURCES) symmetric-aes.c

libpqcrystals_dilithium3aes_ref.so: $(SOURCES) $(HEADERS) symmetric-aes.c
	$(CC) -shared -fPIC $(CFLAGS) -DDILITHIUM_MODE=3 -DDILITHIUM_USE_AES \
	   -o $@ $(SOURCES) symmetric-aes.c

libpqcrystals_dilithium5aes_ref.so: $(SOURCES) $(HEADERS) symmetric-aes.c
	$(CC) -shared -fPIC $(CFLAGS) -DDILITHIUM_MODE=5 -DDILITHIUM_USE_AES \
	   -o $@ $(SOURCES) symmetric-aes.c

build2/libpqcrystals_dilithium2_ref.a: $(OBJ2) $(HEADERS) 
	$(AR) rcs -o $@ $(OBJ2) 

build3/libpqcrystals_dilithium3_ref.a: $(OBJ3) $(HEADERS) 
	$(AR) rcs -o $@ $(OBJ3) 

build5/libpqcrystals_dilithium5_ref.a:  $(OBJ5) $(HEADERS) 
	$(AR) rcs -o $@ $(OBJ5) 

#-----------------------DDILITHIUM_MODE=2------------------------|
build2/sign.o: sign.c $(HEAD1)
	$(CC) $(CFLAGS) -c -DDILITHIUM_MODE=2 -DDILITHIUM_PREFIX=pqcrystals_dilithium2 -o $@ sign.c

build2/packing.o: packing.c $(HEAD1)
	$(CC) $(CFLAGS) -c -DDILITHIUM_MODE=2 -DDILITHIUM_PREFIX=pqcrystals_dilithium2 -o $@ packing.c

build2/polyvec.o: polyvec.c $(HEAD1)
	$(CC) $(CFLAGS) -c -DDILITHIUM_MODE=2 -DDILITHIUM_PREFIX=pqcrystals_dilithium2 -o $@ polyvec.c

build2/poly.o: poly.c $(HEAD1)
	$(CC) $(CFLAGS) -c -DDILITHIUM_MODE=2 -DDILITHIUM_PREFIX=pqcrystals_dilithium2 -o $@ poly.c

build2/ntt.o: ntt.c $(HEAD1)
	$(CC) $(CFLAGS) -c -DDILITHIUM_MODE=2 -DDILITHIUM_PREFIX=pqcrystals_dilithium2 -o $@ ntt.c

build2/reduce.o: reduce.c $(HEAD1)
	$(CC) $(CFLAGS) -c -DDILITHIUM_MODE=2 -DDILITHIUM_PREFIX=pqcrystals_dilithium2 -o $@ reduce.c

build2/rounding.o: rounding.c $(HEAD1)
	$(CC) $(CFLAGS) -c -DDILITHIUM_MODE=2 -DDILITHIUM_PREFIX=pqcrystals_dilithium2 -o $@ rounding.c

build2/fips202.o: fips202.c $(HEAD1)
	$(CC) $(CFLAGS) -c -DDILITHIUM_MODE=2 -DDILITHIUM_PREFIX=pqcrystals_dilithium2 -o $@ fips202.c

build2/symmetric-shake.o: symmetric-shake.c $(HEAD1)
	$(CC) $(CFLAGS) -c -DDILITHIUM_MODE=2 -DDILITHIUM_PREFIX=pqcrystals_dilithium2 -o $@ symmetric-shake.c

build2/randombytes.o: ../../../randombytes/randombytes.c $(HEAD1)
	$(CC) $(CFLAGS) -c -DDILITHIUM_MODE=2 -DDILITHIUM_PREFIX=pqcrystals_dilithium2 -o $@ ../../../randombytes/randombytes.c

#-----------------------DDILITHIUM_MODE=3------------------------|
build3/sign.o: sign.c $(HEAD1)
	$(CC) $(CFLAGS) -c -DDILITHIUM_MODE=3 -DDILITHIUM_PREFIX=pqcrystals_dilithium3 -o $@ sign.c

build3/packing.o: packing.c $(HEAD1)
	$(CC) $(CFLAGS) -c -DDILITHIUM_MODE=3 -DDILITHIUM_PREFIX=pqcrystals_dilithium3 -o $@ packing.c

build3/polyvec.o: polyvec.c $(HEAD1)
	$(CC) $(CFLAGS) -c -DDILITHIUM_MODE=3 -DDILITHIUM_PREFIX=pqcrystals_dilithium3 -o $@ polyvec.c

build3/poly.o: poly.c $(HEAD1)
	$(CC) $(CFLAGS) -c -DDILITHIUM_MODE=3 -DDILITHIUM_PREFIX=pqcrystals_dilithium3 -o $@ poly.c

build3/ntt.o: ntt.c $(HEAD1)
	$(CC) $(CFLAGS) -c -DDILITHIUM_MODE=3 -DDILITHIUM_PREFIX=pqcrystals_dilithium3 -o $@ ntt.c

build3/reduce.o: reduce.c $(HEAD1)
	$(CC) $(CFLAGS) -c -DDILITHIUM_MODE=3 -DDILITHIUM_PREFIX=pqcrystals_dilithium3 -o $@ reduce.c

build3/rounding.o: rounding.c $(HEAD1)
	$(CC) $(CFLAGS) -c -DDILITHIUM_MODE=3 -DDILITHIUM_PREFIX=pqcrystals_dilithium3 -o $@ rounding.c

build3/fips202.o: fips202.c $(HEAD1)
	$(CC) $(CFLAGS) -c -DDILITHIUM_MODE=3 -DDILITHIUM_PREFIX=pqcrystals_dilithium3 -o $@ fips202.c

build3/symmetric-shake.o: symmetric-shake.c $(HEAD1)
	$(CC) $(CFLAGS) -c -DDILITHIUM_MODE=3 -DDILITHIUM_PREFIX=pqcrystals_dilithium3 -o $@ symmetric-shake.c

build3/randombytes.o: ../../../randombytes/randombytes.c $(HEAD1)
	$(CC) $(CFLAGS) -c -DDILITHIUM_MODE=3 -DDILITHIUM_PREFIX=pqcrystals_dilithium3 -o $@ ../../../randombytes/randombytes.c
#-----------------------DDILITHIUM_MODE=5------------------------|
build5/sign.o: sign.c $(HEAD1)
	$(CC) $(CFLAGS) -c -DDILITHIUM_MODE=5 -DDILITHIUM_PREFIX=pqcrystals_dilithium5 -o $@ sign.c

build5/packing.o: packing.c $(HEAD1)
	$(CC) $(CFLAGS) -c -DDILITHIUM_MODE=5 -DDILITHIUM_PREFIX=pqcrystals_dilithium5 -o $@ packing.c

build5/polyvec.o: polyvec.c $(HEAD1)
	$(CC) $(CFLAGS) -c -DDILITHIUM_MODE=5 -DDILITHIUM_PREFIX=pqcrystals_dilithium5 -o $@ polyvec.c

build5/poly.o: poly.c $(HEAD1)
	$(CC) $(CFLAGS) -c -DDILITHIUM_MODE=5 -DDILITHIUM_PREFIX=pqcrystals_dilithium5 -o $@ poly.c

build5/ntt.o: ntt.c $(HEAD1)
	$(CC) $(CFLAGS) -c -DDILITHIUM_MODE=5 -DDILITHIUM_PREFIX=pqcrystals_dilithium5 -o $@ ntt.c

build5/reduce.o: reduce.c $(HEAD1)
	$(CC) $(CFLAGS) -c -DDILITHIUM_MODE=5 -DDILITHIUM_PREFIX=pqcrystals_dilithium5 -o $@ reduce.c

build5/rounding.o: rounding.c $(HEAD1)
	$(CC) $(CFLAGS) -c -DDILITHIUM_MODE=5 -DDILITHIUM_PREFIX=pqcrystals_dilithium5 -o $@ rounding.c

build5/fips202.o: fips202.c $(HEAD1)
	$(CC) $(CFLAGS) -c -DDILITHIUM_MODE=5 -DDILITHIUM_PREFIX=pqcrystals_dilithium5 -o $@ fips202.c

build5/symmetric-shake.o: symmetric-shake.c $(HEAD1)
	$(CC) $(CFLAGS) -c -DDILITHIUM_MODE=5 -DDILITHIUM_PREFIX=pqcrystals_dilithium5 -o $@ symmetric-shake.c

build5/randombytes.o: ../../../randombytes/randombytes.c $(HEAD1)
	$(CC) $(CFLAGS) -c -DDILITHIUM_MODE=5 -DDILITHIUM_PREFIX=pqcrystals_dilithium5 -o $@ ../../../randombytes/randombytes.c
	
test/test_dilithium2: test/test_dilithium.c randombytes.c $(KECCAK_SOURCES) \
  $(KECCAK_HEADERS)
	$(CC) $(CFLAGS) -DDILITHIUM_MODE=2 \
	  -o $@ $< randombytes.c $(KECCAK_SOURCES)

test/test_dilithium3: test/test_dilithium.c randombytes.c $(KECCAK_SOURCES) \
  $(KECCAK_HEADERS)
	$(CC) $(CFLAGS) -DDILITHIUM_MODE=3 \
	  -o $@ $< randombytes.c $(KECCAK_SOURCES)

test/test_dilithium5: test/test_dilithium.c randombytes.c $(KECCAK_SOURCES) \
  $(KECCAK_HEADERS)
	$(CC) $(CFLAGS) -DDILITHIUM_MODE=5 \
	  -o $@ $< randombytes.c $(KECCAK_SOURCES)

test/test_dilithium2aes: test/test_dilithium.c randombytes.c $(AES_SOURCES) \
  $(AES_HEADERS)
	$(CC) $(CFLAGS) -DDILITHIUM_MODE=2 -DDILITHIUM_USE_AES \
	  -o $@ $< randombytes.c $(AES_SOURCES)

test/test_dilithium3aes: test/test_dilithium.c randombytes.c $(AES_SOURCES) \
  $(AES_HEADERS)
	$(CC) $(CFLAGS) -DDILITHIUM_MODE=3 -DDILITHIUM_USE_AES \
	  -o $@ $< randombytes.c $(AES_SOURCES)

test/test_dilithium5aes: test/test_dilithium.c randombytes.c $(AES_SOURCES) \
  $(AES_HEADERS)
	$(CC) $(CFLAGS) -DDILITHIUM_MODE=5 -DDILITHIUM_USE_AES \
	  -o $@ $< randombytes.c $(AES_SOURCES)

test/test_vectors2: test/test_vectors.c $(KECCAK_SOURCES) \
  $(KECCAK_HEADERS)
	$(CC) $(CFLAGS) -DDILITHIUM_MODE=2 \
	  -o $@ $< $(KECCAK_SOURCES)

test/test_vectors3: test/test_vectors.c $(KECCAK_SOURCES) $(KECCAK_HEADERS)
	$(CC) $(CFLAGS) -DDILITHIUM_MODE=3 \
	  -o $@ $< $(KECCAK_SOURCES)

test/test_vectors5: test/test_vectors.c $(KECCAK_SOURCES) \
  $(KECCAK_HEADERS)
	$(CC) $(CFLAGS) -DDILITHIUM_MODE=5 \
	  -o $@ $< $(KECCAK_SOURCES)

test/test_vectors2aes: test/test_vectors.c $(AES_SOURCES) \
  $(AES_HEADERS)
	$(CC) $(CFLAGS) -DDILITHIUM_MODE=2 -DDILITHIUM_USE_AES \
	  -o $@ $< $(AES_SOURCES)

test/test_vectors3aes: test/test_vectors.c $(AES_SOURCES) \
  $(AES_HEADERS)
	$(CC) $(CFLAGS) -DDILITHIUM_MODE=3 -DDILITHIUM_USE_AES \
	  -o $@ $< $(AES_SOURCES)

test/test_vectors5aes: test/test_vectors.c $(AES_SOURCES) \
  $(AES_HEADERS)
	$(CC) $(NISTFLAGS) -DDILITHIUM_MODE=5 -DDILITHIUM_USE_AES \
	  -o $@ $< $(AES_SOURCES)

test/test_speed2: test/test_speed.c test/speed_print.c test/speed_print.h \
  test/cpucycles.c test/cpucycles.h randombytes.c $(KECCAK_SOURCES) \
  $(KECCAK_HEADERS)
	$(CC) $(CFLAGS) -DDILITHIUM_MODE=2 \
	  -o $@ $< test/speed_print.c test/cpucycles.c randombytes.c \
	  $(KECCAK_SOURCES)

test/test_speed3: test/test_speed.c test/speed_print.c test/speed_print.h \
  test/cpucycles.c test/cpucycles.h randombytes.c $(KECCAK_SOURCES) \
  $(KECCAK_HEADERS)
	$(CC) $(CFLAGS) -DDILITHIUM_MODE=3 \
	  -o $@ $< test/speed_print.c test/cpucycles.c randombytes.c \
	  $(KECCAK_SOURCES)

test/test_speed5: test/test_speed.c test/speed_print.c test/speed_print.h \
  test/cpucycles.c test/cpucycles.h randombytes.c $(KECCAK_SOURCES) \
  $(KECCAK_HEADERS)
	$(CC) $(CFLAGS) -DDILITHIUM_MODE=5 \
	  -o $@ $< test/speed_print.c test/cpucycles.c randombytes.c \
	  $(KECCAK_SOURCES)

test/test_speed2aes: test/test_speed.c test/speed_print.c test/speed_print.h \
  test/cpucycles.c test/cpucycles.h randombytes.c $(AES_SOURCES) $(AES_HEADERS)
	$(CC) $(CFLAGS) -DDILITHIUM_MODE=2 -DDILITHIUM_USE_AES \
	  -o $@ $< test/speed_print.c test/cpucycles.c randombytes.c \
	  $(AES_SOURCES)

test/test_speed3aes: test/test_speed.c test/speed_print.c test/speed_print.h \
  test/cpucycles.c test/cpucycles.h randombytes.c $(AES_SOURCES) $(AES_HEADERS)
	$(CC) $(CFLAGS) -DDILITHIUM_MODE=3 -DDILITHIUM_USE_AES \
	  -o $@ $< test/speed_print.c test/cpucycles.c randombytes.c \
	  $(AES_SOURCES)

test/test_speed5aes: test/test_speed.c test/speed_print.c test/speed_print.h \
  test/cpucycles.c test/cpucycles.h randombytes.c $(AES_SOURCES) $(AES_HEADERS)
	$(CC) $(CFLAGS) -DDILITHIUM_MODE=5 -DDILITHIUM_USE_AES \
	  -o $@ $< test/speed_print.c test/cpucycles.c randombytes.c \
	  $(AES_SOURCES)

test/test_mul: test/test_mul.c randombytes.c $(KECCAK_SOURCES) $(KECCAK_HEADERS)
	$(CC) $(CFLAGS) -UDBENCH -o $@ $< randombytes.c $(KECCAK_SOURCES)

PQCgenKAT_sign: PQCgenKAT_sign.c rng.c rng.h $(KECCAK_SOURCES) \
  $(KECCAK_HEADERS)
	$(CC) $(NISTFLAGS) -o $@ $< rng.c $(KECCAK_SOURCES) $(LDFLAGS) -lcrypto

PQCgenKAT_sign2: PQCgenKAT_sign.c rng.c rng.h $(KECCAK_SOURCES) \
  $(KECCAK_HEADERS)
	$(CC) $(NISTFLAGS) -DDILITHIUM_MODE=2 \
	  -o $@ $< rng.c $(KECCAK_SOURCES) $(LDFLAGS) -lcrypto

PQCgenKAT_sign3: PQCgenKAT_sign.c rng.c rng.h $(KECCAK_SOURCES) \
  $(KECCAK_HEADERS)
	$(CC) $(NISTFLAGS) -DDILITHIUM_MODE=3 \
	  -o $@ $< rng.c $(KECCAK_SOURCES) $(LDFLAGS) -lcrypto

PQCgenKAT_sign5: PQCgenKAT_sign.c rng.c rng.h $(KECCAK_SOURCES) \
  $(KECCAK_HEADERS)
	$(CC) $(NISTFLAGS) -DDILITHIUM_MODE=5 \
	  -o $@ $< rng.c $(KECCAK_SOURCES) $(LDFLAGS) -lcrypto

PQCgenKAT_sign2aes: PQCgenKAT_sign.c rng.c rng.h $(AES_SOURCES) \
  $(AES_HEADERS)
	$(CC) $(NISTFLAGS) -DDILITHIUM_MODE=2 -DDILITHIUM_USE_AES \
	  -o $@ $< rng.c $(AES_SOURCES) $(LDFLAGS) -lcrypto

PQCgenKAT_sign3aes: PQCgenKAT_sign.c rng.c rng.h $(AES_SOURCES) \
  $(AES_HEADERS)
	$(CC) $(NISTFLAGS) -DDILITHIUM_MODE=3 -DDILITHIUM_USE_AES \
	  -o $@ $< rng.c $(AES_SOURCES) $(LDFLAGS) -lcrypto

PQCgenKAT_sign5aes: PQCgenKAT_sign.c rng.c rng.h $(AES_SOURCES) \
  $(AES_HEADERS)
	$(CC) $(NISTFLAGS) -DDILITHIUM_MODE=5 -DDILITHIUM_USE_AES \
	  -o $@ $< rng.c $(AES_SOURCES) $(LDFLAGS) -lcrypto

clean:
	rm -f *~ test/*~ *.gcno *.gcda *.lcov
	rm -f libpqcrystals_dilithium2_ref.so
	rm -f libpqcrystals_dilithium3_ref.so
	rm -f libpqcrystals_dilithium5_ref.so
	rm -f libpqcrystals_dilithium2aes_ref.so
	rm -f libpqcrystals_dilithium3aes_ref.so
	rm -f libpqcrystals_dilithium5aes_ref.so
	rm -f libpqcrystals_fips202_ref.so
	rm -f libpqcrystals_aes256ctr_ref.so
	rm -f test/test_dilithium2
	rm -f test/test_dilithium3
	rm -f test/test_dilithium5
	rm -f test/test_dilithium2aes
	rm -f test/test_dilithium3aes
	rm -f test/test_dilithium5aes
	rm -f test/test_vectors2
	rm -f test/test_vectors3
	rm -f test/test_vectors5
	rm -f test/test_vectors2aes
	rm -f test/test_vectors3aes
	rm -f test/test_vectors5aes
	rm -f test/test_speed2
	rm -f test/test_speed3
	rm -f test/test_speed5
	rm -f test/test_speed2aes
	rm -f test/test_speed3aes
	rm -f test/test_speed5aes
	rm -f test/test_mul
	rm -f PQCgenKAT_sign2
	rm -f PQCgenKAT_sign3
	rm -f PQCgenKAT_sign5
	rm -f PQCgenKAT_sign2aes
	rm -f PQCgenKAT_sign3aes
	rm -f PQCgenKAT_sign5aes
	rm -f build2/libpqcrystals_dilithium2_ref.a
	rm -f build2/*.o
	rm -f build3/libpqcrystals_dilithium3_ref.a
	rm -f build3/*.o
	rm -f build5/libpqcrystals_dilithium5_ref.a
	rm -f build5/*.o
