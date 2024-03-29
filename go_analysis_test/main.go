package main

import (
	"flag"
	"fmt"
)

var path string
var goModuleName string

func main() {
	flag.StringVar(&path, "path", "", "package path")
	flag.StringVar(&goModuleName, "go_module", "", "go module name")
	flag.Parse()
	analysis := NewVarAnalyzer()
	err := Analysis(path, goModuleName, analysis.Analyzer)
	if err != nil {
		fmt.Println(err)
	}
	analysis.Print()
}
