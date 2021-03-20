## pprof_cgo

实验 cgo 下 dump heap ，用来分析内存泄漏

## gperftools 编译

```shell
./get_gperftools.sh
```

编译好后（删除了些不必要文件，方便查看）：

```shell
[fananchong@vm-centos7 pprof_cgo]$ tree lib
lib
├── bin
│   ├── pprof
│   └── pprof-symbolize
├── include
│   ├── google
│   │   ├── heap-checker.h
│   │   ├── heap-profiler.h
│   │   ├── malloc_extension_c.h
│   │   ├── malloc_extension.h
│   │   ├── malloc_hook_c.h
│   │   ├── malloc_hook.h
│   │   ├── profiler.h
│   │   ├── stacktrace.h
│   │   └── tcmalloc.h
│   └── gperftools
│       ├── heap-checker.h
│       ├── heap-profiler.h
│       ├── malloc_extension_c.h
│       ├── malloc_extension.h
│       ├── malloc_hook_c.h
│       ├── malloc_hook.h
│       ├── nallocx.h
│       ├── profiler.h
│       ├── stacktrace.h
│       └── tcmalloc.h
└── lib
    ├── libprofiler.a
    ├── libtcmalloc.a
    ├── libtcmalloc_and_profiler.a
    ├── libtcmalloc_debug.a
    ├── libtcmalloc_minimal.a
    └── libtcmalloc_minimal_debug.a
```

- lib/bin/pprof 工具，用它查看 .heap 文件
- lib/lib/libtcmalloc.a 查看内存泄漏，link 它，或者用它的 so


## 编译测试程序

```shell
make clean && make
```

## 运行测试

类似如下输出：

```shell
[fananchong@vm-centos7 pprof_cgo]$ ./main 
Starting tracking the heap
Dumping heap profile to 2573.0001.heap (dump)
[fananchong@vm-centos7 pprof_cgo]$ lib/bin/pprof ./main 2573.0001.heap 
Using local file ./main.
Using local file 2573.0001.heap.
Welcome to pprof!  For help, type 'help'.
(pprof) top10
Total: 17.0 MB
     9.0  52.9%  52.9%      9.0  52.9% f
     8.0  47.1% 100.0%      8.0  47.1% _cgo_fc6d911b7b49_Cfunc__Cmalloc
     0.0   0.0% 100.0%      8.0  47.1% 0x00007ffc8cf2b477
     0.0   0.0% 100.0%      8.0  47.1% 0x00007ffc8cf2b55f
     0.0   0.0% 100.0%      9.0  52.9% runtime.asmcgocall
     0.0   0.0% 100.0%      9.0  52.9% test_malloc
(pprof) quit
```

要输出 pdf 文件，使用命令：

```shell
lib/bin/pprof --pdf ./main 2573.0001.heap > 2573.0001.pdf
```

要比较 2 个快照，查看内存泄漏，使用命令：

```shell
lib/bin/pprof --pdf ./main --base=11990.0001.heap 11990.0002.heap > diff.pdf
```

## 其他

- 本例子演示的是使用 libtcmalloc.a 静态库的方式，**也可以用动态库**（自行百度）
- gcc 加编译选项 **-O3** 的话， **f() 函数的内存泄漏会检查不到**。原因未知
  - 因此最好区分 debug/release 版本。 debug 版请不要添加 -O3 选项，可用来做内存泄漏
