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
	"os"
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

	// Find the node for the global variable "x".
	var x []*ssa.Global
	for _, pkg := range prog.AllPackages() {
		for _, member := range pkg.Members {
			if global, ok := member.(*ssa.Global); ok && strings.HasPrefix(global.Name(), "MyVar") {
				x = append(x, global)
			}
		}
	}

	type targetInfo struct {
		caller    *callgraph.Node
		globalvar *ssa.Global
		callee    *callgraph.Node
	}

	var targets []*targetInfo
	seen := make(map[*callgraph.Node]bool)
	if err := callgraph.GraphVisitEdges(cg, func(edge *callgraph.Edge) error {
		caller := edge.Caller
		if seen[caller] {
			return nil
		}
		seen[caller] = true
		for _, block := range caller.Func.Blocks {
			for _, instr := range block.Instrs {
				if call, ok := instr.(*ssa.Call); ok && len(call.Call.Args) > 0 {
					if call.Call.Signature().Recv() != nil {
						var embedded bool
						v := getGlobalValue(call.Call.Args[0], &embedded)
						for _, xv := range x {
							if xv == v {
								targets = append(targets, &targetInfo{caller, xv, edge.Callee})
								if embedded {
									embeddedCall[caller.ID] = append(embeddedCall[caller.ID], call)
								}
							}
						}
					}
				}
			}
		}
		return nil
	}); err != nil {
		return err
	}
	for _, target := range targets {
		printeAllPath(target.caller, []*callgraph.Node{target.callee})
	}
	return nil
}

var embeddedCall = map[int][]*ssa.Call{}

func getGlobalValue(arg ssa.Value, embedded *bool) *ssa.Global {
	if v, ok := arg.(*ssa.Global); ok {
		return v
	} else if v, ok := arg.(*ssa.UnOp); ok {
		return getGlobalValue(v.X, embedded)
	} else if v, ok := arg.(*ssa.FieldAddr); ok {
		*embedded = true
		return getGlobalValue(v.X, embedded)
	} else {
		return nil
	}
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

func printeAllPath(target *callgraph.Node, path []*callgraph.Node) {
	newPath := append([]*callgraph.Node{target}, path...)

	if len(newPath) > 2 {
		n0 := newPath[0]
		n1 := newPath[1]
		n2 := newPath[2]
		if calls, ok := embeddedCall[n0.ID]; ok {
			for _, call := range calls {
				if n1.Func == call.Call.Value {
					var find bool
					for i := 1; i < len(call.Call.Args); i++ {
						if n2.Func == call.Call.Args[i] {
							find = true
						}
					}
					if !find {
						return
					}
				}
			}
		}
	}

	if len(target.In) == 0 {
		if target.Func.Pkg.Pkg.Name() != "main" {
			return
		}
		if target.Func.Pkg.Pkg.Path() != "go_analysis_test_example/app1" {
			return
		}

		fmt.Printf(newPath[0].Func.String())
		for i := 1; i < len(newPath); i++ {
			fmt.Printf(" --> ")
			fmt.Printf(newPath[i].Func.String())
		}
		fmt.Printf("\n")
		return
	} else {
		for _, child := range target.In {
			printeAllPath(child.Caller, newPath)
		}
	}
}
