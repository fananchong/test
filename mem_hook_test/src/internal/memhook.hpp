#pragma once

#include <stdio.h>
#include <stdlib.h>

// extern "C"
// {
//     void *__real_malloc(size_t size);
//     void __real_free(void *ptr);
//     void *__wrap_malloc(size_t size);
//     void __wrap_free(void *ptr);
// }

void init_mem_hook();
