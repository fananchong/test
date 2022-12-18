#include "internal/memhook.hpp"
#include <string>

void test1(size_t size)
{
    int *p = (int *)malloc(size);
    free(p);
    printf("\n");
}

void test2(size_t size)
{
    int *p = new int[2000];
    delete[] p;
    printf("\n");
}

int main()
{
    init_mem_hook();

    test1(100);
    test1(1024);
    test1(1023);
    test2(2000);

    return 0;
}
