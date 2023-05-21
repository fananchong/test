package main

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

func runVarAnalyzer(pass *analysis.Pass, analyzer *VarAnalyzer) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.CallExpr:
				handleFuncNodeVarAnalyzer(pass, analyzer, file, n, x, nil)
			}
			return true
		})
	}
	return nil, nil
}

func handleFuncNodeVarAnalyzer(pass *analysis.Pass, analyzer *VarAnalyzer, file *ast.File, rawNode ast.Node, n interface{}, sel *ast.Ident) {
	switch x := n.(type) {
	case *ast.CallExpr:
		handleFuncNodeVarAnalyzer(pass, analyzer, file, rawNode, x.Fun, nil)
	case *ast.Ident:
		if obj := pass.TypesInfo.ObjectOf(x); obj != nil {
			if sel != nil {
				analyzer.vars[x.Name] = obj.Type().String()
			}
		}
	case *ast.SelectorExpr:
		sels := getAllSel(x)
		if len(sels) == 1 {
			handleFuncNodeVarAnalyzer(pass, analyzer, file, rawNode, x.X, x.Sel)
		} else {
			handleFuncNodeVarAnalyzer(pass, analyzer, file, rawNode, sels[1], sels[0])
		}
	}
}

type VarAnalyzer struct {
	*analysis.Analyzer
	vars map[string]string
}

func NewVarAnalyzer() *VarAnalyzer {
	analyzer := &VarAnalyzer{
		vars: map[string]string{},
	}
	analyzer.Analyzer = &analysis.Analyzer{
		Name: "var",
		Doc:  "prints var",
		Run:  func(p *analysis.Pass) (interface{}, error) { return runVarAnalyzer(p, analyzer) },
	}
	return analyzer
}

func (analyzer *VarAnalyzer) Print() {
	for k, v := range analyzer.vars {
		fmt.Println(k, v)
	}
}
