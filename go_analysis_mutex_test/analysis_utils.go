package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/callgraph"
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

func checkVarHaveMutex(prog *ssa.Program, caller *callgraph.Node, m, v *types.Var) bool {
	findInstr := func(block *ssa.BasicBlock, v *types.Var) (instrs []ssa.Instruction) {
		for _, instr := range block.Instrs {
			for _, op := range instr.Operands(nil) {
				if varRef, ok := (*op).(*ssa.Global); ok && varRef.Object() == v {
					instrs = append(instrs, instr)
				}
			}
		}
		return
	}

	checkMutex := func(mInstr []ssa.Instruction) {
		if mInstr == nil {
			return
		}
		for _, instr := range mInstr {
			if c, ok := instr.(*ssa.Call); ok {
				fn := c.Common().StaticCallee()
				if fn.Name() == "Unlock" || fn.Name() == "RUnlock" {
					pos := prog.Fset.Position(instr.Pos())
					fmt.Printf("[mutex lint] %v:%v mutex 没有使用 defer 方式，调用 Unlock/RUnlock\n", pos.Filename, pos.Line)
					continue
				}
			}
		}
	}

	checkVar := func(mInstr, vInstr []ssa.Instruction) bool {
		if mInstr == nil && vInstr == nil {
			panic("不会走到这里，逻辑错误")
		}
		if mInstr != nil && vInstr != nil {
			mPos := prog.Fset.Position(mInstr[0].Pos())
			vPos := prog.Fset.Position(vInstr[0].Pos())
			if mPos.Line < vPos.Line {
				return true
			}
		}
		return false
	}

	var find bool
	for _, block := range caller.Func.Blocks {
		mInstr := findInstr(block, m)
		vInstr := findInstr(block, v)
		checkMutex(mInstr)
		if checkVar(mInstr, vInstr) {
			find = true
			break
		}
	}
	return find
}

func checkHaveMutex(prog *ssa.Program, caller *callgraph.Node, m *types.Var) bool {
	findInstr := func(block *ssa.BasicBlock, v *types.Var) (instrs []ssa.Instruction) {
		for _, instr := range block.Instrs {
			for _, op := range instr.Operands(nil) {
				if varRef, ok := (*op).(*ssa.Global); ok && varRef.Object() == v {
					instrs = append(instrs, instr)
				}
			}
		}
		return
	}

	checkMutex := func(mInstr []ssa.Instruction) {
		if mInstr == nil {
			return
		}
		for _, instr := range mInstr {
			if c, ok := instr.(*ssa.Call); ok {
				fn := c.Common().StaticCallee()
				if fn.Name() == "Unlock" || fn.Name() == "RUnlock" {
					pos := prog.Fset.Position(instr.Pos())
					fmt.Printf("[mutex lint] %v:%v mutex 没有使用 defer 方式，调用 Unlock/RUnlock\n", pos.Filename, pos.Line)
					continue
				}
			}
		}
	}

	var find bool
	for _, block := range caller.Func.Blocks {
		mInstr := findInstr(block, m)
		checkMutex(mInstr)
		if len(mInstr) > 0 {
			find = true
			break
		}
	}
	return find
}

func printPaht(newPath []*callgraph.Node, looped bool) string {
	s := newPath[0].Func.String()
	for i := 1; i < len(newPath); i++ {
		s += " --> " + newPath[i].Func.String()
	}
	if looped {
		s += " [LOOP]"
	}
	return s
}
