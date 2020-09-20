## set -e

脚本 [set_e.sh](set_e.sh) 中，执行某命令返回值`非0`，中断后续脚本命令执行，退出

比如：

```vim
set -e
echo "aaa"
ls /xxxxxx
echo "bbb"
```

输出为：

```vim
aaa
ls: cannot access '/xxxxxx': No such file or directory
```

## set -o pipefail

脚本 [set_pipefail.sh](set_pipefail.sh) 中，管道中命令执行失败，返回`非0`
不加这行命令，只返回管道中最后一个命令的执行结果值



比如：

```vim
#!/bin/bash
set -o pipefail
ls /xxxxxx | echo "hello"
echo $?
```

输出为：

```vim
hello
ls: cannot access '/xxxxxx': No such file or directory
2
```

若不加这行：

```vim
#!/bin/bash
#set -o pipefail
ls /xxxxxx | echo "hello"
echo $?
```

输出为：
```vim
hello
ls: cannot access '/xxxxxx': No such file or directory
0
```

## trap

脚本 [trap.sh](trap.sh) 中，收到某信号量执行某函数

比如：

```vim
on_exit() {
  echo "xxxxxx"
}
trap on_exit EXIT
```

输出为：

```vim
fananchong@ali-ubuntu:~/test/shell_test$ ./trap.sh 
xxxxxx
```

## 自定义函数 not

脚本 [not.sh](not.sh) 中，确保一定返回`非0`的命令，返回`0值`

比如：

```vim
echo "aaa"
echo $?
not echo "aaa"
echo $?
```

输出为：

```vim
aaa
0
aaa
1
```

## read && tee

read 接收标准输入，有返回 0 ；没有返回 1
tee 把标准输入，到标准输出；同时到其他设备

比如：

```vim
#!/bin/bash

ls . | read                         # cmd return value: 0
echo "cmd return value: "$?
ls . > ~/aa.log | read              # cmd return value: 1
echo "cmd return value: "$?
ls . | tee /dev/stderr |  read      # cmd return value: 0
echo "cmd return value: "$?
````

输出为：

```vim
fananchong@ali-ubuntu:~/test/shell_test$ ./read_and_tee.sh 
cmd return value: 0
cmd return value: 1
git.sh
not.sh
read_and_tee.sh
README.md
set_e.sh
set_pipefail.sh
trap.sh
cmd return value: 0
```

## pushd && popd


切换目录的命令

比如：

```vim
pwd
pushd ../
pwd
dirs -v
popd
pwd
```

输出为：

```vim
fananchong@ali-ubuntu:~/test/shell_test$ ./pushd_and_popd.sh 
/home/fananchong/test/shell_test
~/test ~/test/shell_test
/home/fananchong/test
 0  ~/test
 1  ~/test/shell_test
~/test/shell_test
/home/fananchong/test/shell_test
```
