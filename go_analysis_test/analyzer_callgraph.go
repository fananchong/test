package main

import (
	"fmt"
	"go/ast"
	"go/types"
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

func (cg *callGraph) addNode(pass *analysis.Pass, analyzer *CallGraphAnalyzer, file *ast.File, x ast.Node, nodeName string, obj types.Object) {
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
	var name string
	ast.Inspect(file, func(n ast.Node) bool {
		if n == nil {
			return false
		}
		switch x := n.(type) {
		case *ast.FuncDecl:
			if x.Pos() <= node.Pos() && x.End() >= node.End() {
				if x.Recv != nil {
					field := x.Recv.List[0]
					var ident *ast.Ident
					switch t := field.Type.(type) {
					case *ast.Ident:
						ident = t
						name = ident.Name + ":" + x.Name.Name
						panic("!!!!")
					case *ast.StarExpr:
						obj2 := pass.TypesInfo.ObjectOf(t.X.(*ast.Ident))
						name = getFuncName2(pass, analyzer, obj2, x.Name.Name)
					}
				} else {
					name = getFuncName1(pass, analyzer, x.Name, pass.TypesInfo.ObjectOf(x.Name))
				}

			}
			return false
		}
		return true
	})
	if name != "" {
		if parent, ok := cg.Nodes[name]; ok {
			return parent
		}
		return &callGraphNode{
			parent:   map[string]*callGraphNode{},
			children: map[string]*callGraphNode{},
			name:     name,
		}
	}
	return nil
}

func (cg *callGraph) print() {
	for _, node := range cg.Nodes {
		if len(node.parent) == 0 {
			printAllPaths(node, "")
		}
	}
}

func printAllPaths(node *callGraphNode, path string) {
	path += node.name + " -> "
	if len(node.children) == 0 {
		fmt.Println(path[:len(path)-4])
	} else {
		for _, child := range node.children {
			printAllPaths(child, path)
		}
	}
}

func getFuncName1(pass *analysis.Pass, analyzer *CallGraphAnalyzer, ident *ast.Ident, obj types.Object) string {
	s := func2pkg[ident] + "." + obj.Name()
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
		return obj.Name() + "." + name
	}
}

func run(pass *analysis.Pass, analyzer *CallGraphAnalyzer) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.CallExpr:
				if ident, ok := x.Fun.(*ast.Ident); ok {
					if obj := pass.TypesInfo.ObjectOf(ident); obj != nil {
						analyzer.cg.addNode(pass, analyzer, file, x, getFuncName1(pass, analyzer, ident, obj), obj)
					}
				} else if se, ok := x.Fun.(*ast.SelectorExpr); ok {
					se2 := getLastSelectorExpr(se)
					if obj := pass.TypesInfo.ObjectOf(se2.Sel); obj != nil {
						obj2 := pass.TypesInfo.ObjectOf(se2.X.(*ast.Ident))
						objname := getFuncName2(pass, analyzer, obj2, obj.Name())
						analyzer.cg.addNode(pass, analyzer, file, x, objname, obj)
					}
				}
			}
			return true
		})
	}
	return nil, nil
}

func getLastSelectorExpr(se *ast.SelectorExpr) *ast.SelectorExpr {
	if _, ok := se.X.(*ast.Ident); ok {
		return se
	} else if _, ok := se.X.(*ast.CallExpr); ok {
		panic("111")
	} else {
		return getLastSelectorExpr(se.X.(*ast.SelectorExpr))
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
