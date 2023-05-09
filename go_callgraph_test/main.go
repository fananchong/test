// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// callgraph: a tool for reporting the call graph of a Go program.
// See Usage for details, or run with -help.
package main // import "golang.org/x/tools/cmd/callgraph"

// TODO(adonovan):
//
// Features:
// - restrict graph to a single package
// - output
//   - functions reachable from root (use digraph tool?)
//   - unreachable functions (use digraph tool?)
//   - dynamic (runtime) types
//   - indexed output (numbered nodes)
//   - JSON output
//   - additional template fields:
//     callee file/line/col

import (
	"flag"
	"fmt"
	"go/build"
	"go/token"
	"io"
	"os"
	"os/exec"
	"sort"
	"strings"

	"golang.org/x/tools/go/buildutil"
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/callgraph/cha"
	"golang.org/x/tools/go/callgraph/rta"
	"golang.org/x/tools/go/callgraph/static"
	"golang.org/x/tools/go/callgraph/vta"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

// flags
var (
	algoFlag = flag.String("algo", "vta",
		`Call graph construction algorithm (static, cha, rta, vta)`)

	testFlag = flag.Bool("test", false,
		"Loads test code (*_test.go) for imported packages")
	excludePkgs = flag.String("exclude_pkgs", "",
		`Exclude some libraries specifically`)
)

func init() {
	flag.Var((*buildutil.TagsFlag)(&build.Default.BuildTags), "tags", buildutil.TagsFlagDoc)
}

const Usage = `callgraph: display the call graph of a Go program.

Usage:

  callgraph [-algo=static|cha|rta|vta|pta] [-test] [-format=...] package...

Flags:

-algo      Specifies the call-graph construction algorithm, one of:

            static      static calls only (unsound)
            cha         Class Hierarchy Analysis
            rta         Rapid Type Analysis
            vta         Variable Type Analysis
            pta         inclusion-based Points-To Analysis

           The algorithms are ordered by increasing precision in their
           treatment of dynamic calls (and thus also computational cost).
           RTA and PTA require a whole program (main or test), and
           include only functions reachable from main.

-test      Include the package's tests in the analysis.
`

func main() {
	flag.Parse()
	if err := doCallgraph(*algoFlag, *testFlag, flag.Args()); err != nil {
		fmt.Fprintf(os.Stderr, "callgraph: %s\n", err)
		os.Exit(1)
	}
}

var stdout io.Writer = os.Stdout

func doCallgraph(algo string, tests bool, args []string) error {
	cfg := &packages.Config{
		Mode:  packages.LoadAllSyntax,
		Tests: tests,
	}

	initial, err := packages.Load(cfg, args...)
	if err != nil {
		return err
	}
	if packages.PrintErrors(initial) > 0 {
		return fmt.Errorf("packages contain errors")
	}

	// Create and build SSA-form program representation.
	mode := ssa.InstantiateGenerics // instantiate generics by default for soundness
	prog, pkgs := ssautil.AllPackages(initial, mode)
	prog.Build()

	// -- call graph construction ------------------------------------------

	var cg *callgraph.Graph

	switch algo {
	case "static":
		cg = static.CallGraph(prog)

	case "cha":
		cg = cha.CallGraph(prog)

	case "rta":
		mains, err := mainPackages(pkgs)
		if err != nil {
			return err
		}
		var roots []*ssa.Function
		for _, main := range mains {
			roots = append(roots, main.Func("init"), main.Func("main"))
		}
		rtares := rta.Analyze(roots, true)
		cg = rtares.CallGraph

		// NB: RTA gives us Reachable and RuntimeTypes too.

	case "vta":
		cg = vta.CallGraph(ssautil.AllFunctions(prog), cha.CallGraph(prog))

	default:
		return fmt.Errorf("unknown algorithm: %s", algo)
	}

	cg.DeleteSyntheticNodes()

	// -- output------------------------------------------------------------

	excludePkgs := getexcludePkgs(*excludePkgs)
	for _, node := range cg.Nodes {
		if len(node.In) == 0 {
			printAllPaths(node, "", excludePkgs)
		}
	}

	// tmpl, err := template.New("-format").Parse("{{.Caller}}\t--{{.Dynamic}}-{{.Line}}:{{.Column}}-->\t{{.Callee}}")
	// if err != nil {
	// 	return fmt.Errorf("invalid -format template: %v", err)
	// }
	// var buf bytes.Buffer
	// data := Edge{fset: prog.Fset}

	// if err := callgraph.GraphVisitEdges(cg, func(edge *callgraph.Edge) error {
	// 	data.position.Offset = -1
	// 	data.edge = edge
	// 	data.Caller = edge.Caller.Func
	// 	data.Callee = edge.Callee.Func

	// 	buf.Reset()
	// 	if err := tmpl.Execute(&buf, &data); err != nil {
	// 		return err
	// 	}
	// 	stdout.Write(buf.Bytes())
	// 	if len := buf.Len(); len == 0 || buf.Bytes()[len-1] != '\n' {
	// 		fmt.Fprintln(stdout)
	// 	}
	// 	return nil
	// }); err != nil {
	// 	return err
	// }

	return nil
}

// mainPackages returns the main packages to analyze.
// Each resulting package is named "main" and has a main function.
func mainPackages(pkgs []*ssa.Package) ([]*ssa.Package, error) {
	var mains []*ssa.Package
	for _, p := range pkgs {
		if p != nil && p.Pkg.Name() == "main" && p.Func("main") != nil {
			mains = append(mains, p)
		}
	}
	if len(mains) == 0 {
		return nil, fmt.Errorf("no main packages")
	}
	return mains, nil
}

type Edge struct {
	Caller *ssa.Function
	Callee *ssa.Function

	edge     *callgraph.Edge
	fset     *token.FileSet
	position token.Position // initialized lazily
}

func (e *Edge) pos() *token.Position {
	if e.position.Offset == -1 {
		e.position = e.fset.Position(e.edge.Pos()) // called lazily
	}
	return &e.position
}

func (e *Edge) Filename() string { return e.pos().Filename }
func (e *Edge) Column() int      { return e.pos().Column }
func (e *Edge) Line() int        { return e.pos().Line }
func (e *Edge) Offset() int      { return e.pos().Offset }

func (e *Edge) Dynamic() string {
	if e.edge.Site != nil && e.edge.Site.Common().StaticCallee() == nil {
		return "dynamic"
	}
	return "static"
}

func (e *Edge) Description() string { return e.edge.Description() }

func printAllPaths(node *callgraph.Node, path string, excludePkgs map[string]struct{}) {
	if node.Func == nil {
		if len(path)-4 > 0 {
			fmt.Println(path[:len(path)-4])
		}
		return
	}
	fname := node.Func.String()
	if strings.Contains(path, fname+" ") {
		path += fname + " <- [LOOP]"
		fmt.Println(path)
	} else {
		pkgname := node.Func.Pkg.Pkg.Path()
		if _, ok := excludePkgs[pkgname]; ok {
			if len(path)-4 > 0 {
				fmt.Println(path[:len(path)-4])
			}
		} else {
			path += fname + " -> "
			if len(node.Out) == 0 {
				fmt.Println(path[:len(path)-4])
			} else {
				for _, k := range node.Out {
					child := k.Callee
					printAllPaths(child, path, excludePkgs)
				}
			}
		}
	}
}

func getexcludePkgs(excludePkgs string) map[string]struct{} {
	cmd := exec.Command("go", "list", "std")
	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	pkgs := make(map[string]struct{})
	for _, pkg := range strings.Split(string(output), "\n") {
		pkgs[pkg] = struct{}{}
		v := strings.ReplaceAll(pkg, "'", "")
		if strings.HasPrefix(pkg, "vendor/") {
			v = strings.TrimPrefix(pkg, "vendor/")
		}
		pkgs[v] = struct{}{}
	}
	cmd = exec.Command("go", "list", "-m", "-f", "'{{.Path}}'", "all")
	output, err = cmd.Output()
	if err != nil {
		panic(err)
	}
	for _, pkg := range strings.Split(string(output), "\n") {
		pkgs[pkg] = struct{}{}
		v := strings.ReplaceAll(pkg, "'", "")
		if strings.HasPrefix(pkg, "vendor/") {
			v = strings.TrimPrefix(pkg, "vendor/")
		}
		pkgs[v] = struct{}{}
	}
	for _, v := range strings.Split(excludePkgs, ",") {
		if v != "" {
			pkgs[v] = struct{}{}
		}
	}

	keys := make([]string, 0, len(pkgs))
	for k := range pkgs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if k != "" {
			fmt.Printf("[Exclude] package %v\n", k)
		}
	}
	return pkgs
}
