package main

import (
	"flag"
	"fmt"
)

var path string

func main() {
	flag.StringVar(&path, "path", "", "package path")
	flag.Parse()
	analysis := NewVarAnalyzer()
	err := Analysis(path, analysis.Analyzer)
	if err != nil {
		fmt.Println(err)
	}
	analysis.Print()
}
