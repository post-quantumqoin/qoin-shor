{
  'targets': [
    {
      'target_name': 'napiSha3',
      'cflags': [
        '-O2', 
      ],
      'configuration': { 'Release': { 'msvs_settings': { 'VCCLCompilerTool': {
          'Optimization': 0
      }}}},
      'sources': [ 
        'shake3.cc',
        'libs/sha3.c'
      ],
      'include_dirs': ["<!(node -p \"require('node-addon-api').include_dir\")"],
      'defines': [ 'NAPI_DISABLE_CPP_EXCEPTIONS' ],
    },
  ],
}