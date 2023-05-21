package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/ssa"
)

// comment := ""
// if valueSpec.Comment != nil {
// 	comment = strings.Trim(valueSpec.Comment.Text(), " ")
// }
// fmt.Println(comment)

func (analyzer *VarAnalyzer) runOne(pass *analysis.Pass) (interface{}, error) {
	// 找出 mutex 全局变量
	for _, file := range pass.Files {
		for i, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.VAR {
				continue
			}
			var mutexValueSpec *ast.ValueSpec
			for _, spec := range genDecl.Specs {
				if valueSpec, ok := spec.(*ast.ValueSpec); ok {
					if !isMutexType(valueSpec.Type) {
						continue
					}
					ident := valueSpec.Names[0]
					obj := pass.TypesInfo.Defs[ident]
					if obj == nil {
						continue
					}
					v, _ := obj.(*types.Var)
					isGlobal := !v.IsField() && !v.Embedded() && v.Parent() == pass.Pkg.Scope() // 全局变量
					if isGlobal {
						mutexValueSpec = valueSpec
					}
				}
			}
			if mutexValueSpec != nil {
				pos := pass.Fset.Position(mutexValueSpec.Pos())
				comment := ""
				if mutexValueSpec.Comment != nil {
					comment = strings.ReplaceAll(mutexValueSpec.Comment.Text(), " ", "")
					comment = strings.ReplaceAll(comment, "\n", "")
				}
				if comment == "" {
					fmt.Printf("[mutex lint] %v:%v mutex 变量没有注释，指明它要保护的变量\n", pos.Filename, pos.Line)
					os.Exit(1)
				}
				varNames := strings.Split(comment, ",")
				if i+1+len(varNames) > len(file.Decls) {
					fmt.Printf("[mutex lint] %v:%v mutex 变量注释有误，它要保护的变量未声明\n", pos.Filename, pos.Line)
					os.Exit(1)
				}
				for j := 1; j <= len(varNames); j++ {
					genDecl, ok := file.Decls[i+j].(*ast.GenDecl)
					if !ok || genDecl.Tok != token.VAR {
						fmt.Printf("[mutex lint] %v:%v mutex 变量注释中的变量 %v ，声明不对\n", pos.Filename, pos.Line+j, varNames[j-1])
						os.Exit(1)
					}
					spec := genDecl.Specs[0]
					if valueSpec, ok := spec.(*ast.ValueSpec); !ok || valueSpec.Names[0].Name != varNames[j-1] {
						pos := pass.Fset.Position(spec.Pos())
						fmt.Printf("[mutex lint] %v:%v mutex 变量注释中的变量 %v ，声明不对\n", pos.Filename, pos.Line, varNames[j-1])
						os.Exit(1)
					} else {
						analyzer.vars[valueSpec] = mutexValueSpec
					}
				}
			}
		}
	}

	return nil, nil
}

type VarAnalyzer struct {
	*analysis.Analyzer
	cg   *callgraph.Graph
	prog *ssa.Program
	vars map[*ast.ValueSpec]*ast.ValueSpec // key : 变量； value mutex
}

func NewVarAnalyzer(cg *callgraph.Graph, prog *ssa.Program) *VarAnalyzer {
	analyzer := &VarAnalyzer{
		cg:   cg,
		prog: prog,
		vars: map[*ast.ValueSpec]*ast.ValueSpec{},
	}
	analyzer.Analyzer = &analysis.Analyzer{
		Name: "var",
		Doc:  "prints var",
		Run:  func(p *analysis.Pass) (interface{}, error) { return analyzer.runOne(p) },
	}
	return analyzer
}

func (analyzer *VarAnalyzer) Print() {
	// for k, v := range analyzer.vars {
	// 	fmt.Println(k, v)
	// }
}
