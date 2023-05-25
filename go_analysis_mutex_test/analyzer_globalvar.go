package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/ssa"
)

func (analyzer *VarAnalyzer) runOne(prog *ssa.Program, pass *analysis.Pass) (interface{}, error) {
	// 1. 获取需要加锁的全局变量 A
	analyzer.step1FindGlobalVar(pass)
	return nil, nil
}

func (analyzer *VarAnalyzer) step1FindGlobalVar(pass *analysis.Pass) {
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
			if mutexValueSpec == nil {
				continue
			}
			pos := pass.Fset.Position(mutexValueSpec.Pos())
			mutexVar := getGlobalVarByPos(analyzer.prog, pos)
			comment := ""
			if mutexValueSpec.Comment != nil {
				comment = strings.ReplaceAll(mutexValueSpec.Comment.Text(), " ", "")
				comment = strings.ReplaceAll(comment, "\n", "")
			}
			if comment == "" {
				fmt.Printf("[mutex lint] %v:%v mutex 变量没有注释，指明它要保护的变量\n", pos.Filename, pos.Line)
				continue
			}
			varNames := strings.Split(comment, ",")
			if i+1+len(varNames) > len(file.Decls) {
				fmt.Printf("[mutex lint] %v:%v mutex 变量注释有误，它要保护的变量未声明\n", pos.Filename, pos.Line)
				continue
			}
			for j := 1; j <= len(varNames); j++ {
				genDecl, ok := file.Decls[i+j].(*ast.GenDecl)
				if !ok || genDecl.Tok != token.VAR {
					fmt.Printf("[mutex lint] %v:%v mutex 变量注释中的变量 %v ，声明不对\n", pos.Filename, pos.Line+j, varNames[j-1])
					break
				}
				spec := genDecl.Specs[0]
				if valueSpec, ok := spec.(*ast.ValueSpec); !ok || valueSpec.Names[0].Name != varNames[j-1] {
					pos := pass.Fset.Position(spec.Pos())
					fmt.Printf("[mutex lint] %v:%v mutex 变量注释中的变量 %v ，声明不对\n", pos.Filename, pos.Line, varNames[j-1])
					break
				} else {
					pos := pass.Fset.Position(spec.Pos())
					v := getGlobalVarByPos(analyzer.prog, pos)
					analyzer.vars[v] = mutexVar
				}
			}
		}
	}
}

func (analyzer *VarAnalyzer) step2FindCaller() {
	seen := make(map[*callgraph.Node]bool)
	f := func(edge *callgraph.Edge) error {
		caller := edge.Caller
		if seen[caller] {
			return nil
		}
		if caller.Func == nil {
			return nil
		}
		seen[caller] = true
		if caller.Func.Name() == "init" {
			return nil
		}
		usesVar := func(instr ssa.Instruction, v *types.Var) bool {
			for _, op := range instr.Operands(nil) {
				if varRef, ok := (*op).(*ssa.Global); ok && varRef.Object() == v {
					return true
				}
			}
			return false
		}
		for _, block := range caller.Func.Blocks {
			for _, instr := range block.Instrs {
				for k := range analyzer.vars {
					if usesVar(instr, k) {
						if _, ok := analyzer.callers[k]; !ok {
							analyzer.callers[k] = make(map[*callgraph.Node]bool)
						}
						analyzer.callers[k][caller] = true
					}
				}
			}
		}
		return nil
	}
	if err := callgraph.GraphVisitEdges(analyzer.cg, f); err != nil {
		return
	}
	return
}

func (analyzer *VarAnalyzer) step3CutCaller() {
	for v := range analyzer.callers {
		m := analyzer.vars[v]
		callers := analyzer.callers[v]

		for caller := range callers {
			if !checkVarHaveMutex(analyzer.prog, caller, m, v) {
				if _, ok := analyzer.callers2[v]; !ok {
					analyzer.callers2[v] = make(map[*callgraph.Node]bool)
				}
				analyzer.callers2[v][caller] = true
			}
		}
	}
}

func (analyzer *VarAnalyzer) step4CheckPath(myvar *types.Var, target *callgraph.Node, paths []*callgraph.Node, seen map[*callgraph.Node]bool, checkFail *string) {
	if seen[target] {
		return
	}
	seen[target] = true

	if *checkFail != "" {
		return
	}

	newPaths := append([]*callgraph.Node{target}, paths...)
	var looped bool
	for _, v := range paths {
		if v.Func == target.Func {
			looped = true
			break
		}
	}

	// 检查是否有 mutex
	mymutex := analyzer.vars[myvar]
	if checkHaveMutex(analyzer.prog, target, mymutex) {
		return
	}

	// 如果超出本包，则报错
	if target.Func.Pkg.Pkg != myvar.Pkg() {
		*checkFail = printPaht(newPaths, looped)
		return
	}

	if len(target.In) == 0 || looped {
		*checkFail = printPaht(newPaths, looped)
		return
	} else {
		for _, in := range target.In {
			analyzer.step4CheckPath(myvar, in.Caller, newPaths, seen, checkFail)
		}
	}
}

type VarAnalyzer struct {
	*analysis.Analyzer
	path     string
	cg       *callgraph.Graph
	prog     *ssa.Program
	vars     map[*types.Var]*types.Var // key : 变量； value mutex
	callers  map[*types.Var]map[*callgraph.Node]bool
	callers2 map[*types.Var]map[*callgraph.Node]bool
}

func NewVarAnalyzer(path string, cg *callgraph.Graph, prog *ssa.Program) *VarAnalyzer {
	analyzer := &VarAnalyzer{
		path:     path,
		cg:       cg,
		prog:     prog,
		vars:     map[*types.Var]*types.Var{},
		callers:  map[*types.Var]map[*callgraph.Node]bool{},
		callers2: map[*types.Var]map[*callgraph.Node]bool{},
	}
	analyzer.Analyzer = &analysis.Analyzer{
		Name: "var",
		Doc:  "prints var",
		Run:  func(p *analysis.Pass) (interface{}, error) { return analyzer.runOne(prog, p) },
	}
	return analyzer
}

func (analyzer *VarAnalyzer) Analysis() {
	err := Analysis(analyzer.path, analyzer.Analyzer)
	if err != nil {
		panic(err)
	}
	// 2. 获取哪些函数 B ，直接使用了全局变量 A
	analyzer.step2FindCaller()
	// 3. 剔除 B 中有加锁的函数，得 C
	analyzer.step3CutCaller()
	// 4. 查看调用关系，逆向检查上级调用是否加锁
	for v, nodes := range analyzer.callers2 {
		for node := range nodes {
			var checkFail string
			analyzer.step4CheckPath(v, node, []*callgraph.Node{}, map[*callgraph.Node]bool{}, &checkFail)
			if checkFail != "" {
				pos := analyzer.prog.Fset.Position(v.Pos())
				fmt.Printf("[mutex lint] %v:%v 没有调用 mutex lock 。调用链：%v\n", pos.Filename, pos.Line, checkFail)
				break
			}
		}
	}
}

func (analyzer *VarAnalyzer) Print() {
}
