package main

import (
	"fmt"
	"go/types"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

type A struct{}

func (a *A) F(p1 int, p2 string) {
	fmt.Println("A")
}

func main() {
	// 设置构建参数
	conf := &packages.Config{
		Mode:  packages.LoadAllSyntax,
		Tests: false,
	}

	// 加载包
	pkgs, err := packages.Load(conf, "./...")
	if err != nil {
		fmt.Println("Error loading packages:", err)
		return
	}

	// 创建Program
	prog, _ := ssautil.AllPackages(pkgs, ssa.SanityCheckFunctions)

	// 查找结构体类型
	var structType *types.Named
	for _, pkg := range prog.AllPackages() {
		obj := pkg.Pkg.Scope().Lookup("A")
		if obj != nil {
			structType, _ = obj.Type().(*types.Named)
			break
		}
	}

	if structType == nil {
		fmt.Println("Struct type A not found.")
		return
	}

	// 遍历结构体的方法
	for i := 0; i < structType.NumMethods(); i++ {
		method := structType.Method(i)
		// 获取方法名
		methodName := method.Name()

		// 获取方法签名
		sig := method.Type().(*types.Signature)

		// 获取参数名
		params := make([]string, sig.Params().Len())
		for i := 0; i < sig.Params().Len(); i++ {
			params[i] = sig.Params().At(i).Name()
		}

		// 打印方法名和参数名
		fmt.Printf("%s: %v\n", methodName, params)
	}
}
