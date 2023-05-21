package main

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

// comment := ""
// if valueSpec.Comment != nil {
// 	comment = strings.Trim(valueSpec.Comment.Text(), " ")
// }
// fmt.Println(comment)

func runVarAnalyzer(pass *analysis.Pass, analyzer *VarAnalyzer) (interface{}, error) {
	for _, file := range pass.Files {
		mutexVars := map[string][]*ast.ValueSpec{}
		ast.Inspect(file, func(node ast.Node) bool {
			if genDecl, ok := node.(*ast.GenDecl); ok && genDecl.Tok == token.VAR {
				for _, spec := range genDecl.Specs {
					if valueSpec, ok := spec.(*ast.ValueSpec); ok {
						if valueSpec.Type == nil {
							continue
						}
						if isMutexType(valueSpec.Type) {
							mutexVars[valueSpec.Names[0].String()] = append(mutexVars[valueSpec.Names[0].String()], valueSpec)
						}
					}
				}
			}
			return true
		})
	}
	return nil, nil
}

type VarAnalyzer struct {
	*analysis.Analyzer
	vars map[*ast.ValueSpec]*ast.ValueSpec // key : 变量； value mutex
}

func NewVarAnalyzer() *VarAnalyzer {
	analyzer := &VarAnalyzer{
		vars: map[*ast.ValueSpec]*ast.ValueSpec{},
	}
	analyzer.Analyzer = &analysis.Analyzer{
		Name: "var",
		Doc:  "prints var",
		Run:  func(p *analysis.Pass) (interface{}, error) { return runVarAnalyzer(p, analyzer) },
	}
	return analyzer
}

func (analyzer *VarAnalyzer) Print() {
	// for k, v := range analyzer.vars {
	// 	fmt.Println(k, v)
	// }
}
