package main

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

func myanalysis(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			if call, ok := node.(*ast.CallExpr); ok {
				if ident, ok := call.Fun.(*ast.Ident); ok && ident.Obj != nil {
					pos := pass.Fset.Position(call.Pos())
					line := pos.Line
					var funcname string
					for _, decl := range file.Decls {
						funcDecl, ok := decl.(*ast.FuncDecl)
						if !ok || funcDecl.Body == nil {
							continue
						}

						posBegin := pass.Fset.Position(funcDecl.Body.Lbrace)
						posEnd := pass.Fset.Position(funcDecl.Body.Rbrace)
						if posBegin.Line > line || posEnd.Line < line {
							continue
						}
						if posBegin.Filename != pos.Filename {
							continue
						}
						funcname = funcDecl.Name.Name
						// fmt.Printf("Function %s calls fmt.Println at line %d\n", funcDecl.Name.Name, line)
						break
					}

					fmt.Printf("%s %s calls %s\n", pos, funcname, ident.Obj.Name)
				}
			}
			return true
		})
	}
	return nil, nil
}

var analyzerCallgraph = &analysis.Analyzer{
	Name: "callgraph",
	Doc:  "prints the call graph",
	Run:  myanalysis,
}
