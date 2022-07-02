#include <fcntl.h>
#include <sched.h>
#include <unistd.h>

/**
[fananchong@qa4-haidao test1]$ ls -l /proc/$$/ns/net
lrwxrwxrwx. 1 fananchong fananchong 0 7月   2 18:45 /proc/3085834/ns/net -> net:[4026531968]
[fananchong@qa4-haidao test1]$ ./a.out 4026531968 /bin/bash
 */

int main(int argc, char **argv)
{
    int fd = open(argv[1], O_RDONLY); // 获取 namespace 的文件描述符
    setns(fd, 0);                     // 加入 namespace
    execvp(argv[2], &argv[2]);        // 执行用户自定义程序
    return 0;
}