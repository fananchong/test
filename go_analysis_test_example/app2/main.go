package main

import (
	"go_analysis_test_example/core/dir3"
)

func main() {
	// dir0.F2()
	// x2.F2()
	// A2()

	// go func() {
	// 	dir0.F3()
	// }()
	dir3.MyVar3.Walk(F2)
}
func F2() {
	dir3.MyVar2.FFF()
}
