#include "backtrace.hpp"
#include <stdio.h>
#include <execinfo.h>
#include <unistd.h>
#include <stdlib.h>
#define BACKTRACE_SIZE 100

void print_backtrace()
{
    void *buffer[BACKTRACE_SIZE] = {0};
    int pointer_num = backtrace(buffer, BACKTRACE_SIZE);
    backtrace_symbols_fd(buffer, pointer_num, fileno(stderr));
    return;
}
