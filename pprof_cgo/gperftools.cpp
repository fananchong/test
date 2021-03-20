#include "gperftools.h"
#include <gperftools/tcmalloc.h>
#include <gperftools/malloc_extension.h>
#include <gperftools/heap-profiler.h>
#include <gperftools/profiler.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

void setup_gperftools()
{
    char prefix[128];
    sprintf(prefix, "%d", getpid());
    HeapProfilerStart(prefix);
}

void dump_heap()
{
    HeapProfilerDump("dump");
}

void f()
{
    char *ptr1 = (char *)malloc(10 * 1024 * 1024);
    ptr1[1000] = 1;
    free(ptr1);
    char *ptr2 = (char *)malloc(30 * 1024 * 1024);
    ptr2[1000] = 1;
}

void test_malloc()
{
    f();
}