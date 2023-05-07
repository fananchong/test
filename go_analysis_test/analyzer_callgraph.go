package main

import (
	"fmt"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

type callGraphNode struct {
	parent   *callGraphNode
	children []*callGraphNode
	name     string
}

func (c *callGraphNode) addChild(child *callGraphNode) {
	c.children = append(c.children, child)
	child.parent = c
}

func (c *callGraphNode) print() {
	if c.parent == nil {
		fmt.Printf("%s", c.name)
	} else {
		fmt.Printf(" --> %s", c.name)
	}
	for _, child := range c.children {
		child.print()
	}
	if len(c.children) == 0 {
		fmt.Printf("\n")
	}
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
		cg.Nodes[nodeName] = &callGraphNode{name: nodeName}
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
						ident = t.X.(*ast.Ident)

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
		return &callGraphNode{name: name}
	}
	return nil
}

func (cg *callGraph) print() {
	for _, node := range cg.Nodes {
		if node.parent == nil {
			node.print()
		}
	}
}

func run(pass *analysis.Pass) (interface{}, error) {
	callGraph := newCallGraph()
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

	callGraph.print()
	return nil, nil
}

func GetCallGraphAnalyzer() *analysis.Analyzer {
	var analyzer = &analysis.Analyzer{
		Name: "callgraph",
		Doc:  "prints the call graph",
		Run:  run,
	}
	return analyzer
}
