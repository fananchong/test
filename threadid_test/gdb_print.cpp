

#include "gdb_print.hpp"
#include <signal.h>
#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <sys/types.h>
#ifndef _MSC_VER
#include <unistd.h>
#endif
#include <limits.h>
#include <vector>
#include <string>

std::vector<std::function<void(void)>> g_fns;
void register_crash_exit_fn(const std::function<void(void)> &fn)
{
    g_fns.push_back(fn);
}

#ifndef _MSC_VER
size_t get_executable_path(char *processdir, char *processname, size_t len)
{
    char *path_end;
    if (readlink("/proc/self/exe", processdir, len) <= 0)
        return -1;
    path_end = strrchr(processdir, '/');
    if (path_end == NULL)
        return -1;
    ++path_end;
    strcpy(processname, path_end);
    *path_end = '\0';
    return (size_t)(path_end - processdir);
}
void print_core(int signum, siginfo_t *info, void *secret, struct sigaction *oldact)
{
    // ERR("crash signum:{} si_code:{}", signum, info->si_code);
    char cmd[128] = {0};
    sprintf(cmd, "gcore %u", getpid());
    system((const char *)cmd);
    for (auto &fn : g_fns)
    {
        fn();
    }
    sprintf(cmd, "./gdb_print.sh ./core.%u %s", getpid(), "");
    // ERR("cmd={}", cmd);
    system((const char *)cmd);
    oldact->sa_sigaction(signum, info, secret);
}

#endif

void sigsetup()
{
#ifndef _MSC_VER
    struct sigaction act;
    memset(&act, 0, sizeof(act));
    act.sa_flags = SA_ONSTACK | SA_SIGINFO;

#define SIGACTION(SIG)                                               \
    static struct sigaction old##SIG;                                \
    act.sa_sigaction = [](int signum, siginfo_t *info, void *secret) \
    {                                                                \
        print_core(signum, info, secret, &old##SIG);                 \
    };                                                               \
    sigaction(SIG, &act, &old##SIG);

    SIGACTION(SIGSEGV);
    SIGACTION(SIGABRT);
    SIGACTION(SIGFPE);
#endif
}
