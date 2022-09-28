#include <iostream>
#include <thread>
#include "gdb_print.hpp"
#include <unistd.h>
#include <sys/syscall.h>

pid_t gettid()
{
    return syscall(SYS_gettid);
}

void f(pid_t tids[2], int i)
{
    tids[i] = gettid();
    std::cout << "[f" << i << "] thread_id = " << tids[i] << std::endl;
    std::this_thread::sleep_for(std::chrono::seconds(5));
    int *p = 0;
    *p = 0;
    std::cout << "[f" << i << "] thread_id = " << tids[i] << std::endl;
}

int main()
{
    sigsetup();
    pid_t tids[2];
    std::thread t0(f, tids, 0);
    std::thread t1(f, tids, 1);
    char cmd[128] = {0};
    sprintf(cmd, "echo %u %u > tids.txt", tids[0], tids[1]);
    system((const char *)cmd);
    t0.join();
    t1.join();
}
