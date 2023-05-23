package main

import (
	"flag"
	"fmt"
)

var path string

func main() {
	flag.StringVar(&path, "path", "", "package path")
	flag.Parse()

	cg, prog, err := doCallgraph("vta", false, []string{fmt.Sprintf("%s/...", path)})
	if err != nil {
		panic(err)
	}

	analyzer := NewVarAnalyzer(path, cg, prog)
	analyzer.Analysis()
	analyzer.Print()
}
