#pragma once

#if __cplusplus
extern "C"
{
#endif

    extern void setup_gperftools();
    extern void dump_heap();
    extern void test_malloc();

#if __cplusplus
}
#endif
