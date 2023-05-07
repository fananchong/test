package main

import (
	"go/ast"
	"os"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/packages"
)

func Analysis(path string, analyzer *analysis.Analyzer) error {
	err := os.Chdir(path)
	if err != nil {
		return err
	}
	packages, err := packages.Load(&packages.Config{
		Mode: packages.LoadAllSyntax,
	}, path+"/...")
	if err != nil {
		return err
	}
	pass := &analysis.Pass{
		Analyzer: analyzer,
		Files:    []*ast.File{},
		ResultOf: map[*analysis.Analyzer]interface{}{},
	}
	for _, pkg := range packages {
		if len(pkg.Errors) > 0 {
			return pkg.Errors[0]
		}
		pass.Fset = pkg.Fset
		pass.Files = pkg.Syntax
		pass.TypesInfo = pkg.TypesInfo
		_, err := analyzer.Run(pass)
		if err != nil {
			return err
		}
	}
	return nil
}
