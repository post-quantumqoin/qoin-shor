{
  'targets': [
    {
      'target_name': 'Falcon512NistRound3',
      'defines': [ 'NODE_ADDON_API_ENABLE_MAYBE', 'NAPI_DISABLE_CPP_EXCEPTIONS' ],
      'cflags!': [ '-fno-exceptions'],
      'cflags_cc!': [ '-fno-exceptions'],
      'cflags': [
        '-W',
        '-O2'
      ],
      'configuration': { 'Release': { 'msvs_settings': { 'VCCLCompilerTool': { 'Optimization': 0 } } } },
      'sources': [
        'falcon512/falcon512Napi.cc',
        '../../randombytes/randombytes.c',
        'falcon512/falconCore/fpr.c',
        'falcon512/falconCore/rng.c',
        'falcon512/falconCore/codec.c',
        'falcon512/falconCore/common.c',
        'falcon512/falconCore/fft.c',
        'falcon512/falconCore/keygen.c',
        'falcon512/falconCore/api.c',
        'falcon512/falconCore/shake.c',
        'falcon512/falconCore/sign.c',
        'falcon512/falconCore/vrfy.c'
      ],
      "include_dirs": [ "<!(node -p \"require('node-addon-api').include_dir\")" ]
    },
    {
      'target_name': 'Falcon1024NistRound3',
      'defines': [ 'NODE_ADDON_API_ENABLE_MAYBE', 'NAPI_DISABLE_CPP_EXCEPTIONS' ],
      'cflags!': [ '-fno-exceptions' ],
      'cflags_cc!': [ '-fno-exceptions' ],
      'configuration': { 'Release': { 'msvs_settings': { 'VCCLCompilerTool': { 'Optimization': 0 } } } },
      'cflags': [
        '-W',
        '-O2'
      ],
      'sources': [
        'falcon1024/falcon1024Napi.cc',
        '../../randombytes/randombytes.c',
        'falcon1024/falconCore/fpr.c',
        'falcon1024/falconCore/rng.c',
        'falcon1024/falconCore/codec.c',
        'falcon1024/falconCore/common.c',
        'falcon1024/falconCore/fft.c',
        'falcon1024/falconCore/keygen.c',
        'falcon1024/falconCore/api.c',
        'falcon1024/falconCore/shake.c',
        'falcon1024/falconCore/sign.c',
        'falcon1024/falconCore/vrfy.c'
      ],
      "include_dirs": [ "<!(node -p \"require('node-addon-api').include_dir\")" ]
    },
  ],
}
