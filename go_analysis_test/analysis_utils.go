package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

func ajustPkgName(name string, goModuleName string) string {
	index1 := strings.Index(name, "\"")
	if index1 > 0 {
		index2 := strings.LastIndex(name, "\"")
		s := name[index1+1 : index2]
		if strings.HasPrefix(s, goModuleName) {
			return s[len(goModuleName)+1:]
		} else {
			return s
		}
	} else {
		v := strings.Split(name, " ")
		return v[1]
	}
}

func ajustAnonymousName(pos token.Position, goModuleName string) string {
	name := fmt.Sprintf("%v:%v", pos.Filename, pos.Line)
	token := goModuleName + "/"
	index := strings.Index(name, token)
	if index > 0 {
		s := fmt.Sprintf("[anonymous %v]", name[index+len(token):])
		return s
	} else {
		s := fmt.Sprintf("[anonymous %v]", name)
		return s
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
