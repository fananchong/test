package main

import (
	"fmt"
	"go/ast"
	"go/types"

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

func (cg *callGraph) addNode(pass *analysis.Pass, file *ast.File, x ast.Node, nodeName string, obj types.Object) {
	if _, ok := cg.Nodes[nodeName]; !ok {
		cg.Nodes[nodeName] = &callGraphNode{
			parent:   map[string]*callGraphNode{},
			children: map[string]*callGraphNode{},
			name:     nodeName,
		}
	}
	node := cg.Nodes[nodeName]
	if parent := cg.getParent(pass, file, x); parent != nil {
		parent.addChild(node)
		if _, ok := cg.Nodes[parent.name]; !ok {
			cg.Nodes[parent.name] = parent
		}
	}
}

func (cg *callGraph) getParent(pass *analysis.Pass, file *ast.File, node ast.Node) *callGraphNode {
	var name string
	ast.Inspect(file, func(n ast.Node) bool {
		if n == nil {
			return false
		}
		switch x := n.(type) {
		case *ast.FuncDecl:
			if x.Pos() <= node.Pos() && x.End() >= node.End() {
				name = x.Name.Name
				if x.Recv != nil {
					field := x.Recv.List[0]
					var ident *ast.Ident
					switch t := field.Type.(type) {
					case *ast.Ident:
						ident = t
						name = ident.Name + ":" + name
					case *ast.StarExpr:
						obj2 := pass.TypesInfo.ObjectOf(t.X.(*ast.Ident))
						if obj2.Type().String() != "invalid type" {
							name = obj2.Type().String() + ":" + name
						}

					}
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

func run(pass *analysis.Pass, callGraph *callGraph) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.CallExpr:
				if ident, ok := x.Fun.(*ast.Ident); ok {
					if obj := pass.TypesInfo.ObjectOf(ident); obj != nil {
						callGraph.addNode(pass, file, x, obj.Name(), obj)
					}
				} else if se, ok := x.Fun.(*ast.SelectorExpr); ok {
					if obj := pass.TypesInfo.ObjectOf(se.Sel); obj != nil {
						objname := obj.Name()
						obj2 := pass.TypesInfo.ObjectOf(se.X.(*ast.Ident))
						if obj2.Type().String() != "invalid type" {
							objname = obj2.Type().String() + ":" + obj.Name()
							if objname[0] == '*' {
								objname = objname[1:]
							}
						}
						callGraph.addNode(pass, file, x, objname, obj)
					}
				}
			}
			return true
		})
	}
	return nil, nil
}

type CallGraphAnalyzer struct {
	*analysis.Analyzer
	cg *callGraph
}

func NewCallGraphAnalyzer() *CallGraphAnalyzer {
	analyzer := &CallGraphAnalyzer{}
	analyzer.cg = newCallGraph()
	analyzer.Analyzer = &analysis.Analyzer{
		Name: "callgraph",
		Doc:  "prints the call graph",
		Run:  func(p *analysis.Pass) (interface{}, error) { return run(p, analyzer.cg) },
	}
	return analyzer
}

func (analyzer *CallGraphAnalyzer) Print() {
	analyzer.cg.print()
}
