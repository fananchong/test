package main

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/ssa"
)

func (analyzer *StructFieldAnalyzer) runOne(prog *ssa.Program, pass *analysis.Pass) (interface{}, error) {
	// 获取含有 mutex 字段的结构体、以及结构体内 mutex 对应的字段 A
	analyzer.step1FindStructField(pass)
	return nil, nil
}

func (analyzer *StructFieldAnalyzer) step1FindStructField(pass *analysis.Pass) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch t := node.(type) {
			case *ast.TypeSpec:
				if structType, ok := t.Type.(*ast.StructType); ok {
					fields := structType.Fields.List
					for i, field := range fields {
						if isMutexType(field.Type) {
							m := getStructFieldByPos(analyzer.prog, pass.Fset.Position(field.Pos()))

							comment := ""
							if field.Comment != nil {
								comment = strings.ReplaceAll(field.Comment.Text(), " ", "")
								comment = strings.ReplaceAll(comment, "\n", "")
							}
							if comment == "" {
								pos := pass.Fset.Position(field.Pos())
								fmt.Printf("[mutex lint] %v:%v mutex 变量没有注释，指明它要保护的变量\n", pos.Filename, pos.Line)
								continue
							}
							varNames := strings.Split(comment, ",")
							if i+1+len(varNames) > len(fields) {
								pos := pass.Fset.Position(field.Pos())
								fmt.Printf("[mutex lint] %v:%v mutex 变量注释有误，它要保护的变量未声明\n", pos.Filename, pos.Line)
								continue
							}
							for j := 1; j <= len(varNames); j++ {
								if fields[i+j].Names[0].Name != varNames[j-1] {
									pos := pass.Fset.Position(fields[i+j].Pos())
									fmt.Printf("[mutex lint] %v:%v mutex 变量注释中的变量 %v ，声明不对\n", pos.Filename, pos.Line, varNames[j-1])
									break
								} else {
									v := getStructFieldByPos(analyzer.prog, pass.Fset.Position(fields[i+j].Pos()))
									analyzer.vars[v] = m
								}
							}
						}
					}
				}
			}
			return true
		})
	}
}

func (analyzer *StructFieldAnalyzer) step2FindCaller() {
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
		n := caller.Func.Name()
		if strings.HasSuffix(n, "f1") {
			fmt.Println(n)
		}
		for _, block := range caller.Func.Blocks {
			for _, instr := range block.Instrs {
				if fieldAddr, ok := instr.(*ssa.FieldAddr); ok && fieldAddr.X != nil {
					if pointerType, ok := fieldAddr.X.Type().Underlying().(*types.Pointer); ok {
						if structType, ok := pointerType.Elem().Underlying().(*types.Struct); ok {
							field := structType.Field(fieldAddr.Field)
							for k := range analyzer.vars {
								if k == field {
									if _, ok := analyzer.callers[k]; !ok {
										analyzer.callers[k] = make(map[*callgraph.Node]bool)
									}
									analyzer.callers[k][caller] = true
									break
								}
							}
						}
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

func (analyzer *StructFieldAnalyzer) step3CutCaller() {
	for v := range analyzer.callers {
		m := analyzer.vars[v]
		callers := analyzer.callers[v]

		for caller := range callers {
			if !checkStructFieldHaveMutex(analyzer.prog, caller, m, v) {
				if _, ok := analyzer.callers2[v]; !ok {
					analyzer.callers2[v] = make(map[*callgraph.Node]bool)
				}
				analyzer.callers2[v][caller] = true
			}
		}
	}
}

func (analyzer *StructFieldAnalyzer) step4CheckPath(myvar *types.Var, target *callgraph.Node, paths []*callgraph.Node, seen map[*callgraph.Node]bool, checkFail *string) {
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
	if checkStructFieldHaveMutex2(analyzer.prog, target, mymutex) {
		return
	}

	// 如果超出本包，则报错
	if target.Func.Pkg.Pkg != myvar.Pkg() {
		*checkFail = printPaht(newPaths, looped)
		return
	}

	// 如果已经是协程起点，则报错
	if isGoroutine(target.Func) {
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

type StructFieldAnalyzer struct {
	*analysis.Analyzer
	path     string
	cg       *callgraph.Graph
	prog     *ssa.Program
	vars     map[*types.Var]*types.Var // key : 变量； value mutex
	callers  map[*types.Var]map[*callgraph.Node]bool
	callers2 map[*types.Var]map[*callgraph.Node]bool
}

func NewStructFieldAnalyzer(path string, cg *callgraph.Graph, prog *ssa.Program) *StructFieldAnalyzer {
	analyzer := &StructFieldAnalyzer{
		path:     path,
		cg:       cg,
		prog:     prog,
		vars:     map[*types.Var]*types.Var{},
		callers:  map[*types.Var]map[*callgraph.Node]bool{},
		callers2: map[*types.Var]map[*callgraph.Node]bool{},
	}
	analyzer.Analyzer = &analysis.Analyzer{
		Name: "struct field",
		Doc:  "prints struct field",
		Run:  func(p *analysis.Pass) (interface{}, error) { return analyzer.runOne(prog, p) },
	}
	return analyzer
}

func (analyzer *StructFieldAnalyzer) Analysis() {
	err := Analysis(analyzer.path, analyzer.Analyzer)
	if err != nil {
		panic(err)
	}
	// 2. 获取哪些函数 B ，直接使用了相关字段
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

func (analyzer *StructFieldAnalyzer) Print() {
}
