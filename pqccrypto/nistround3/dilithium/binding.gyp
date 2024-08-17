{
  'targets': [
    {
      'target_name': 'Dilithium2RefNistRound3',
      'cflags!': [ '-fno-exceptions' ],
      'cflags_cc!': [ '-fno-exceptions' ],
      'cflags': [
        '-Wextra', 
        '-Wpedantic', 
        '-Wmissing-prototypes', 
        '-Wredundant-decls',
        '-Wshadow',
        '-Wvla',
        '-Wpointer-arith',
        '-O2', 
        '-fomit-frame-pointer'
      ],
      'configuration': { 'Release': { 'msvs_settings': { 'VCCLCompilerTool': {
          'Optimization': 0
      }}}},
      'sources': [ 
        'dilithiumNapi.cc',
        'core/sign.c',
        'core/packing.c',
        'core/polyvec.c',
        'core/poly.c',
        'core/ntt.c',
        'core/reduce.c',
        'core/rounding.c',
        '../../randombytes/randombytes.c',
        'core/fips202.c',
        'core/symmetric-shake.c',
      ],
      'defines': [ 'NAPI_DISABLE_CPP_EXCEPTIONS', 'DILITHIUM_MODE=2' ],
      "include_dirs": [ "<!(node -p \"require('node-addon-api').include_dir\")" ]
    },
    {
      'target_name': 'Dilithium3RefNistRound3',
      'cflags!': [ '-fno-exceptions' ],
      'cflags_cc!': [ '-fno-exceptions' ],
      'cflags': [
        '-Wextra', 
        '-Wpedantic', 
        '-Wmissing-prototypes', 
        '-Wredundant-decls',
        '-Wshadow',
        '-Wvla',
        '-Wpointer-arith',
        '-O2', 
        '-fomit-frame-pointer'
      ],
      'configuration': { 'Release': { 'msvs_settings': { 'VCCLCompilerTool': {
          'Optimization': 0
      }}}},
      'sources': [ 
        'dilithiumNapi.cc',
        'core/sign.c',
        'core/packing.c',
        'core/polyvec.c',
        'core/poly.c',
        'core/ntt.c',
        'core/reduce.c',
        'core/rounding.c',
        '../../randombytes/randombytes.c',
        'core/fips202.c',
        'core/symmetric-shake.c',
      ],
      'defines': [ 'NAPI_DISABLE_CPP_EXCEPTIONS', 'DILITHIUM_MODE=3' ],
      "include_dirs": [ "<!(node -p \"require('node-addon-api').include_dir\")" ]
    },
    {
      'target_name': 'Dilithium5RefNistRound3',
      'cflags!': [ '-fno-exceptions' ],
      'cflags_cc!': [ '-fno-exceptions' ],
      'cflags': [
        '-Wall', 
        '-Wextra', 
        '-Wpedantic', 
        '-Wmissing-prototypes', 
        '-Wredundant-decls',
        '-Wshadow',
        '-Wvla',
        '-Wpointer-arith',
        '-O2', 
        '-fomit-frame-pointer'
      ],
      'configuration': { 'Release': { 'msvs_settings': { 'VCCLCompilerTool': {
          'Optimization': 0
      }}}},
      'sources': [ 
        'dilithiumNapi.cc',
        'core/sign.c',
        'core/packing.c',
        'core/polyvec.c',
        'core/poly.c',
        'core/ntt.c',
        'core/reduce.c',
        'core/rounding.c',
        '../../randombytes/randombytes.c',
        'core/fips202.c',
        'core/symmetric-shake.c',
      ],
      'defines': [ 'NAPI_DISABLE_CPP_EXCEPTIONS', 'DILITHIUM_MODE=5' ],
      "include_dirs": [ "<!(node -p \"require('node-addon-api').include_dir\")" ]
    },
  ],
}