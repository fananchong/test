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