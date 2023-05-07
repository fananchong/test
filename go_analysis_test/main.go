package main

import (
	"go/ast"
	"log"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/packages"
)

func main() {
	foo()
	bar()

	packages, err := packages.Load(&packages.Config{
		Mode: packages.LoadAllSyntax,
	}, "./...")
	if err != nil {
		log.Fatal(err)
	}
	pass := &analysis.Pass{
		Analyzer: analyzerCallgraph,
		Files:    []*ast.File{},
		ResultOf: map[*analysis.Analyzer]interface{}{},
	}
	for _, pkg := range packages {
		pass.Fset = pkg.Fset
		pass.Files = pkg.Syntax
		_, err := analyzerCallgraph.Run(pass)
		if err != nil {
			log.Fatal(err)
		}
	}
}
