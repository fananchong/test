## set -e

脚本中，执行某命令返回值`非0`，中断后续脚本命令执行，退出

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

管道中命令执行失败，返回`非0`
不加这行命令，只返回管道中最后一个命令的执行结果值



比如：

```vim
#!/bin/bash
set -o pipefail
ls /xxxxx | echo "hello"
echo $?
```

输出为：

```vim
hello
ls: cannot access '/xxxxx': No such file or directory
2
```

若不加这行：

```vim
#!/bin/bash
#set -o pipefail
ls /xxxxx | echo "hello"
echo $?
```

输出为：
```vim
hello
ls: cannot access '/xxxxx': No such file or directory
0
```

## trap

收到某信号量执行某函数


比如：

```vim
on_exit() {
  echo "xxxxxxxx"
}
trap on_exit EXIT
```

输出为：

```vim
fananchong@ali-ubuntu:~/test/shell_test$ ./trap.sh 
xxxxxxxx
```

## 自定义函数 not

确保一定返回`非0`的命令，返回`0值`

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

## 自定义函数 fail_on_output

有错误到标准错误，则返回`非0`

比如 git.sh 中：

```vim
# Check to make sure it's safe to modify the user's git repo.
git status --porcelain | fail_on_output
echo "aaa"
```

如果没有任何修改或没有未纳入版本库的文件，则输出：

```vim
fananchong@ali-ubuntu:~/test/shell_test$ ./git.sh 
aaa
```

如果有修改或有未纳入版本库的文件，则输出：

```vim
fananchong@ali-ubuntu:~/test/shell_test$ ./git.sh 
 M shell_test/README.md
 D shell_test/set_e.sh
?? shell_test/xxxxxxxxx
```