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
    malloc(8 * 1024 * 1024);
}

void test_malloc()
{
    f();
}