package main

import (
	"flag"
	"fmt"
)

var path string

func main() {
	flag.StringVar(&path, "path", "", "package path")
	flag.Parse()
	err := Analysis(path, GetCallGraphAnalyzer())
	if err != nil {
		fmt.Println(err)
	}
}
