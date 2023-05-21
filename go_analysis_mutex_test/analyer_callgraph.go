package main

import (
	"fmt"
	"os/exec"
	"strings"

	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/callgraph/cha"
	"golang.org/x/tools/go/callgraph/rta"
	"golang.org/x/tools/go/callgraph/static"
	"golang.org/x/tools/go/callgraph/vta"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

func doCallgraph(algo string, tests bool, args []string) (*callgraph.Graph, *ssa.Program, error) {
	cfg := &packages.Config{
		Mode:  packages.LoadAllSyntax, //nolint:staticcheck
		Tests: tests,
	}

	initial, err := packages.Load(cfg, args...)
	if err != nil {
		return nil, nil, err
	}
	if packages.PrintErrors(initial) > 0 {
		return nil, nil, fmt.Errorf("packages contain errors")
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
			return nil, nil, err
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
		return nil, nil, fmt.Errorf("unknown algorithm: %s", algo)
	}

	cg.DeleteSyntheticNodes()

	return cg, prog, nil
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

func getGoPkg() map[string]struct{} {
	cmd := exec.Command("go", "list", "std")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("failed to execute command:", err)
		return nil
	}
	m := map[string]struct{}{}
	pkgList := strings.Split(string(output), "\n")
	for _, pkg := range pkgList {
		if pkg == "" {
			continue
		}
		m[pkg] = struct{}{}
		if strings.HasPrefix(pkg, "vendor/") {
			v := strings.TrimPrefix(pkg, "vendor/")
			m[v] = struct{}{}
		}
	}
	return m
}

var skipPkgs = getGoPkg()
var noSkipPkgs = make(map[string]struct{})

func isSkipByPkg(pkg string) bool {
	if _, ok := skipPkgs[pkg]; ok {
		return true
	}
	if _, ok := noSkipPkgs[pkg]; ok {
		return false
	}
	var skip bool
	for k := range skipPkgs {
		if strings.HasPrefix(pkg, k) {
			skipPkgs[pkg] = struct{}{}
			skip = true
		}
	}
	if !skip {
		noSkipPkgs[pkg] = struct{}{}
	}
	return skip
}

func isErrorPathByPkg(pkg string) bool {
	m := make(map[string]struct{})
	if _, ok := m[pkg]; ok {
		return true
	}
	return false
}
