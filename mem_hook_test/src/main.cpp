#include "internal/memhook.hpp"
#include <string>

int main()
{
    init_mem_hook();

    int *p1 = (int *)malloc(1024);
    free(p1);
    printf("\n");

    int *p2 = (int *)malloc(1023);
    free(p2);
    printf("\n");

    int *p3 = new int[2000];
    delete[] (p3);
    printf("\n");

    return 0;
}
