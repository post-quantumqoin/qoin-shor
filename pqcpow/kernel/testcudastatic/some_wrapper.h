#ifndef SOME_WRAPPER_H
#define SOME_WRAPPER_H

	#ifdef __cplusplus
	extern "C" {
	#endif
	    struct some_item {
	        char version[8];
	        char value[32];
	    };
	
	    struct some_result {
	        struct some_item data[1];
	        int size;
	    };
	
	    int generate(struct some_result* result, char* id);
	
	#ifdef __cplusplus
	}
	#endif

#endif // SOME_WRAPPER_H
