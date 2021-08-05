#include <iostream>

int main(int argc, char **argv)
{
    auto ptr = malloc(1024);
    free(ptr);
    free(ptr);
    printf("done.\n");
    return 0;
}
