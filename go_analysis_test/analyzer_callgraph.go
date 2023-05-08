package main

import (
	"fmt"
	"go/ast"
	"go/types"
	"math"
	"sort"
	"strings"

	"golang.org/x/tools/go/analysis"
)

type callGraphNode struct {
	parent   map[string]*callGraphNode
	children map[string]*callGraphNode
	name     string
}

func (c *callGraphNode) addChild(child *callGraphNode) {
	c.children[child.name] = child
	child.parent[c.name] = c
}

type callGraph struct {
	Nodes map[string]*callGraphNode
}

func newCallGraph() *callGraph {
	return &callGraph{
		Nodes: map[string]*callGraphNode{},
	}
}

func (cg *callGraph) addNode(pass *analysis.Pass, analyzer *CallGraphAnalyzer, file *ast.File, x ast.Node, nodeName string) {
	if _, ok := cg.Nodes[nodeName]; !ok {
		cg.Nodes[nodeName] = &callGraphNode{
			parent:   map[string]*callGraphNode{},
			children: map[string]*callGraphNode{},
			name:     nodeName,
		}
	}
	node := cg.Nodes[nodeName]
	if parent := cg.getParent(pass, analyzer, file, x); parent != nil {
		parent.addChild(node)
		if _, ok := cg.Nodes[parent.name]; !ok {
			cg.Nodes[parent.name] = parent
		}
	}
}

func (cg *callGraph) getParent(pass *analysis.Pass, analyzer *CallGraphAnalyzer, file *ast.File, node ast.Node) *callGraphNode {
	type funcInfo struct {
		Name     string
		PosBegin int
		PosEnd   int
	}
	var ninfo = &funcInfo{
		Name:     "",
		PosBegin: math.MinInt,
		PosEnd:   math.MaxInt,
	}
	pos := pass.Fset.Position(node.Pos())
	line := pos.Line

	for _, x := range funcInFile[pos.Filename] {
		posBegin := pass.Fset.Position(x.Body.Lbrace)
		posEnd := pass.Fset.Position(x.Body.Rbrace)
		if posBegin.Line < line && posEnd.Line > line && posBegin.Filename == pos.Filename && posBegin.Line > ninfo.PosBegin && posEnd.Line < ninfo.PosEnd {
			if x.Recv != nil {
				field := x.Recv.List[0]
				switch t := field.Type.(type) {
				case *ast.StarExpr:
					obj2 := pass.TypesInfo.ObjectOf(t.X.(*ast.Ident))
					ninfo.Name = getFuncName2(pass, analyzer, obj2, x.Name.Name)
				default:
					panic("[getParent] unknow type")
				}
			} else {
				ninfo.Name = getFuncName1(pass, analyzer, x.Name, pass.TypesInfo.ObjectOf(x.Name))
			}
			ninfo.PosBegin = posBegin.Line
			ninfo.PosEnd = posEnd.Line
		}
	}

	for _, x := range anonymousInFile[pos.Filename] {
		posBegin := pass.Fset.Position(x.Body.Lbrace)
		posEnd := pass.Fset.Position(x.Body.Rbrace)
		if posBegin.Line < line && posEnd.Line > line && posBegin.Filename == pos.Filename && posBegin.Line > ninfo.PosBegin && posEnd.Line < ninfo.PosEnd {
			ninfo.Name = ajustAnonymousName(posBegin, analyzer.goModuleName)
			ninfo.PosBegin = posBegin.Line
			ninfo.PosEnd = posEnd.Line
		}
	}

	if ninfo.Name != "" {
		if parent, ok := cg.Nodes[ninfo.Name]; ok {
			return parent
		}
		return &callGraphNode{
			parent:   map[string]*callGraphNode{},
			children: map[string]*callGraphNode{},
			name:     ninfo.Name,
		}
	}
	return nil
}

func (cg *callGraph) print() {
	keys := make([]string, 0, len(cg.Nodes))
	for k := range cg.Nodes {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		node := cg.Nodes[k]
		if len(node.parent) == 0 {
			printAllPaths(node, "")
		}
	}
}

func printAllPaths(node *callGraphNode, path string) {
	if strings.Contains(path, node.name+" ") {
		path += node.name + " <- [LOOP]"
		fmt.Println(path)
	} else {
		path += node.name + " -> "
		if len(node.children) == 0 {
			fmt.Println(path[:len(path)-4])
		} else {
			keys := make([]string, 0, len(node.children))
			for k := range node.children {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				child := node.children[k]
				printAllPaths(child, path)
			}
		}
	}
}

func getFuncName1(pass *analysis.Pass, analyzer *CallGraphAnalyzer, ident *ast.Ident, obj types.Object) string {
	var s string
	if _, ok := obj.(*types.Builtin); !ok {
		s = func2pkg[ident] + "." + obj.Name()
	} else {
		s = "[builtin]." + obj.Name()
	}
	return s
}

func getFuncName2(pass *analysis.Pass, analyzer *CallGraphAnalyzer, obj types.Object, name string) string {
	if obj.Type().String() != "invalid type" {
		s := obj.Type().String() + ":" + name
		if s[0] == '*' {
			s = s[1:]
		}
		if strings.HasPrefix(s, analyzer.goModuleName) {
			return s[len(analyzer.goModuleName)+1:]
		} else {
			return s
		}
	} else {
		s := ajustPkgName(obj.String(), analyzer.goModuleName) + "." + name
		return s
	}
}

func run(pass *analysis.Pass, analyzer *CallGraphAnalyzer) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.CallExpr:
				handleFuncNode(pass, analyzer, file, n, x, nil)
			}
			return true
		})
	}
	return nil, nil
}

func handleFuncNode(pass *analysis.Pass, analyzer *CallGraphAnalyzer, file *ast.File, rawNode ast.Node, n interface{}, sel *ast.Ident) {
	switch x := n.(type) {
	case *ast.CallExpr:
		handleFuncNode(pass, analyzer, file, rawNode, x.Fun, nil)
	case *ast.Ident:
		if obj := pass.TypesInfo.ObjectOf(x); obj != nil {
			if sel == nil {
				objname := getFuncName1(pass, analyzer, x, obj)
				analyzer.cg.addNode(pass, analyzer, file, rawNode, objname)
			} else {
				obj2 := pass.TypesInfo.ObjectOf(sel)
				objname := getFuncName2(pass, analyzer, obj, obj2.Name())
				analyzer.cg.addNode(pass, analyzer, file, rawNode, objname)
			}
		}
	case *ast.SelectorExpr:
		sels := getAllSel(x)
		if len(sels) == 1 {
			handleFuncNode(pass, analyzer, file, rawNode, x.X, x.Sel)
		} else {
			handleFuncNode(pass, analyzer, file, rawNode, sels[1], sels[0])
		}
	case *ast.FuncLit:
		pos := pass.Fset.Position(x.Pos())
		objname := ajustAnonymousName(pos, analyzer.goModuleName)
		analyzer.cg.addNode(pass, analyzer, file, rawNode, objname)
	}
}

func getAllSel(x *ast.SelectorExpr) []*ast.Ident {
	if v, ok := x.X.(*ast.SelectorExpr); ok {
		var s []*ast.Ident
		s = append([]*ast.Ident{x.Sel}, getAllSel(v)...)
		return s
	} else {
		return []*ast.Ident{x.Sel}
	}
}

type CallGraphAnalyzer struct {
	*analysis.Analyzer
	cg           *callGraph
	goModuleName string
}

func NewCallGraphAnalyzer(goModuleName string) *CallGraphAnalyzer {
	analyzer := &CallGraphAnalyzer{}
	analyzer.cg = newCallGraph()
	analyzer.Analyzer = &analysis.Analyzer{
		Name: "callgraph",
		Doc:  "prints the call graph",
		Run:  func(p *analysis.Pass) (interface{}, error) { return run(p, analyzer) },
	}
	analyzer.goModuleName = goModuleName
	return analyzer
}

func (analyzer *CallGraphAnalyzer) Print() {
	analyzer.cg.print()
}
