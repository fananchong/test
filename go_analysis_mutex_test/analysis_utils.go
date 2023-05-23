package main

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/ssa"
)

func isSyncMutexType(expr ast.Expr) bool {
	ident, ok := expr.(*ast.SelectorExpr)
	if !ok || ident.X == nil || ident.Sel == nil {
		return false
	}
	x, ok := ident.X.(*ast.Ident)
	sel := ident.Sel
	if !ok {
		return false
	}
	return sel.Name == "Mutex" && x.Name == "sync"
}

func isSyncRWMutexType(expr ast.Expr) bool {
	ident, ok := expr.(*ast.SelectorExpr)
	if !ok || ident.X == nil || ident.Sel == nil {
		return false
	}
	x, ok := ident.X.(*ast.Ident)
	sel := ident.Sel
	if !ok {
		return false
	}
	return sel.Name == "RWMutex" && x.Name == "sync"
}

func isMutexType(expr ast.Expr) bool {
	return isSyncMutexType(expr) || isSyncRWMutexType(expr)
}

func getGlobalVarByPos(prog *ssa.Program, pos token.Position) *types.Var {
	for _, pkg := range prog.AllPackages() {
		for _, member := range pkg.Members {
			if global, ok := member.(*ssa.Global); ok {
				p := prog.Fset.Position(global.Pos())
				if p == pos {
					return global.Object().(*types.Var)
				}
			}
		}
	}
	return nil
}
