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
	var name string
	pos := pass.Fset.Position(node.Pos())
	line := pos.Line
	for _, x := range funcInFile[pos.Filename] {
		posBegin := pass.Fset.Position(x.Body.Lbrace)
		posEnd := pass.Fset.Position(x.Body.Rbrace)
		if posBegin.Line <= line && posEnd.Line >= line && posBegin.Filename == pos.Filename {
			if x.Recv != nil {
				field := x.Recv.List[0]
				switch t := field.Type.(type) {
				case *ast.StarExpr:
					obj2 := pass.TypesInfo.ObjectOf(t.X.(*ast.Ident))
					name = getFuncName2(pass, analyzer, obj2, x.Name.Name)
				default:
					panic("[getParent] unknow type")
				}
			} else {
				name = getFuncName1(pass, analyzer, x.Name, pass.TypesInfo.ObjectOf(x.Name))
			}
			break
		}
	}
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
	if strings.Contains(path, node.name+" ") {
		path += node.name + "[LOOP]"
		fmt.Println(path)
	} else {
		path += node.name + " -> "
		if len(node.children) == 0 {
			fmt.Println(path[:len(path)-4])
		} else {
			for _, child := range node.children {
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
		return obj.Name() + "." + name
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
		handleFuncNode(pass, analyzer, file, rawNode, x.X, x.Sel)
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
