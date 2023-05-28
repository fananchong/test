package main

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
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

func getStructFieldByPos(prog *ssa.Program, pos token.Position) *types.Var {
	for _, pkg := range prog.AllPackages() {
		for _, member := range pkg.Members {
			if obj, ok := member.(*ssa.Type); ok {
				if s, ok := obj.Type().Underlying().(*types.Struct); ok {
					for i := 0; i < s.NumFields(); i++ {
						field := s.Field(i)
						p := prog.Fset.Position(field.Pos())
						if p == pos {
							return field
						}
					}
				}
			}
		}
	}
	return nil
}

func checkGlobalVarHaveMutex(prog *ssa.Program, caller *callgraph.Node, mymutex, myvar *types.Var) bool {
	var find bool
	for _, block := range caller.Func.Blocks {
		mInstr := findInstrByGlobalVar(block, mymutex)
		vInstr := findInstrByGlobalVar(block, myvar)
		checkMutex(prog, mInstr)
		if checkVar(prog, mInstr, vInstr) {
			find = true
			break
		}
	}
	return find
}

func checkGlobalVarHaveMutex2(prog *ssa.Program, caller *callgraph.Node, m *types.Var) bool {
	var find bool
	for _, block := range caller.Func.Blocks {
		mInstr := findInstrByGlobalVar(block, m)
		checkMutex(prog, mInstr)
		if len(mInstr) > 0 {
			find = true
			break
		}
	}
	return find
}

func checkStructFieldHaveMutex(prog *ssa.Program, caller *callgraph.Node, mymutex, myvar *types.Var) bool {
	var find bool
	for _, block := range caller.Func.Blocks {
		mInstr := findInstrByStructField(block, mymutex)
		vInstr := findInstrByStructField(block, myvar)
		checkMutex(prog, block.Instrs)
		if checkVar(prog, mInstr, vInstr) {
			find = true
			break
		}
	}
	return find
}

func checkStructFieldHaveMutex2(prog *ssa.Program, caller *callgraph.Node, m *types.Var) bool {
	var find bool
	for _, block := range caller.Func.Blocks {
		mInstr := findInstrByStructField(block, m)
		checkMutex(prog, block.Instrs)
		if len(mInstr) > 0 {
			find = true
			break
		}
	}
	return find
}

func checkVar(prog *ssa.Program, mInstr, vInstr []ssa.Instruction) bool {
	if mInstr == nil && vInstr == nil {
		return false
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

func checkMutex(prog *ssa.Program, mInstr []ssa.Instruction) {
	if mInstr == nil {
		return
	}
	// for _, instr := range mInstr {
	// 	if c, ok := instr.(*ssa.Call); ok {
	// 		if c.Call.Value.Name() == "Unlock" || c.Call.Value.Name() == "RUnlock" {
	// 			pos := prog.Fset.Position(instr.Pos())
	// 			fmt.Printf("[mutex lint] %v:%v mutex 没有使用 defer 方式，调用 Unlock/RUnlock\n", pos.Filename, pos.Line)
	// 			continue
	// 		}
	// 	}
	// }
}

func findInstrByGlobalVar(block *ssa.BasicBlock, v *types.Var) (instrs []ssa.Instruction) {
	for _, instr := range block.Instrs {
		for _, op := range instr.Operands(nil) {
			if varRef, ok := (*op).(*ssa.Global); ok && varRef.Object() == v {
				instrs = append(instrs, instr)
			}
		}
	}
	return
}

func findInstrByStructField(block *ssa.BasicBlock, v *types.Var) (instrs []ssa.Instruction) {
	for _, instr := range block.Instrs {
		if fieldAddr, ok := instr.(*ssa.FieldAddr); ok && fieldAddr.X != nil {
			if pointerType, ok := fieldAddr.X.Type().Underlying().(*types.Pointer); ok {
				if structType, ok := pointerType.Elem().Underlying().(*types.Struct); ok {
					field := structType.Field(fieldAddr.Field)
					if field == v {
						instrs = append(instrs, instr)
					}
				}
			}
		}
	}
	return
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

func isGoroutine(fn *ssa.Function) bool {
	if fn.Referrers() != nil {
		for _, r := range *fn.Referrers() {
			if _, ok := r.(*ssa.Go); ok {
				return true
			}
		}
	}
	return false
}

func getGlobalVarByName(pass *analysis.Pass, file *ast.File, name string) *ast.ValueSpec {
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.VAR {
			continue
		}
		for _, spec := range genDecl.Specs {
			if valueSpec, ok := spec.(*ast.ValueSpec); ok {
				ident := valueSpec.Names[0]
				obj := pass.TypesInfo.Defs[ident]
				if obj == nil {
					continue
				}
				v, _ := obj.(*types.Var)
				isGlobal := !v.IsField() && !v.Embedded() && v.Parent() == pass.Pkg.Scope() // 全局变量
				if isGlobal && ident.Name == name {
					return valueSpec
				}
			}
		}
	}
	return nil
}
