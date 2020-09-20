## go vet

go vet 可以静态检查一些可能的错误

如果使用 golang 默认安装它提示的工具集时，会自动安装，并在 IDE 中用`橙色下波浪线`提示


## main.go 中的告警提示

```vim
fananchong@ali-ubuntu:~/test/go_vet_test$ go vet
# github.com/fananchong/test/go_vet_test
./main.go:63:2: self-assignment of i to i
./main.go:31:14: suspect or: i != 0 || i != 1
./main.go:32:14: suspect and: i == 0 && i == 1
./main.go:33:14: redundant and: i == 0 && i == 0
./main.go:55:8: using res before checking for errors
./main.go:41:32: loop variable word captured by func literal
./main.go:67:14: comparison of function test7 == nil is always false
./main.go:12:2: Printf format %d has arg str of wrong type string
./main.go:17:2: Printf format %s has arg &str of wrong type *string
./main.go:26:2: customLogf format %s has arg i of wrong type int
./main.go:72:24: i (64 bits) too small for shift of 64
```
