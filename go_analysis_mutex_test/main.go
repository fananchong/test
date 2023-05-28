package main

import (
	"flag"
	"fmt"
)

var path string

func main() {
	flag.StringVar(&path, "path", ".", "package path")
	flag.Parse()

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
