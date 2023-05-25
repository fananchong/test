package main

import (
	"flag"
	"fmt"
	"go_analysis_mutex_test/test1"
)

var path string

func test() {
	test1.F1()
	test1.F2()
}

func main() {
	flag.StringVar(&path, "path", ".", "package path")
	flag.Parse()

	test()

	cg, prog, err := doCallgraph("vta", false, []string{fmt.Sprintf("%s/...", path)})
	if err != nil {
		panic(err)
	}

	analyzer1 := NewVarAnalyzer(path, cg, prog)
	analyzer1.Analysis()
	analyzer1.Print()

	analyzer2 := NewStructFieldAnalyzer(path, cg, prog)
	analyzer2.Analysis()
	analyzer2.Print()
}
