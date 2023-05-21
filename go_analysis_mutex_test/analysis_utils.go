package main

import "go/ast"

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
