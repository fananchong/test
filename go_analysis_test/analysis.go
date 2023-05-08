package main

import (
	"go/ast"
	"os"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/packages"
)

func Analysis(path, goModuleName string, analyzer *analysis.Analyzer) error {
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
	initCacheData(packages, goModuleName)
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
		pass.Pkg = pkg.Types
		_, err := analyzer.Run(pass)
		if err != nil {
			return err
		}
	}
	return nil
}

var func2pkg = make(map[*ast.Ident]string)
var funcInFile = make(map[string][]*ast.FuncDecl)
var anonymousInFile = make(map[string][]*ast.FuncLit)

func initCacheData(pkgs []*packages.Package, goModuleName string) {
	for _, pkgInfo := range pkgs {
		for _, f := range pkgInfo.Syntax {
			ast.Inspect(f, func(n ast.Node) bool {
				switch x := n.(type) {
				case *ast.FuncLit:
					pos := pkgInfo.Fset.Position(x.Body.Lbrace)
					anonymousInFile[pos.Filename] = append(anonymousInFile[pos.Filename], x)
				case *ast.CallExpr:
					if ident, ok := x.Fun.(*ast.Ident); ok {
						if obj := pkgInfo.TypesInfo.ObjectOf(ident); obj != nil {
							for _, obj2 := range pkgInfo.TypesInfo.Defs {
								if obj2 != nil && obj2.Name() == obj.Name() {
									func2pkg[ident] = pkgInfo.Name
								}
							}
						}
					}
				case *ast.FuncDecl:
					if x.Recv == nil {
						if obj := pkgInfo.TypesInfo.ObjectOf(x.Name); obj != nil {
							for _, obj2 := range pkgInfo.TypesInfo.Defs {
								if obj2 != nil && obj2.Name() == obj.Name() {
									func2pkg[x.Name] = pkgInfo.Name
								}
							}
						}
					}

					pos := pkgInfo.Fset.Position(x.Body.Lbrace)
					funcInFile[pos.Filename] = append(funcInFile[pos.Filename], x)
				}
				return true
			})
		}
	}
}
