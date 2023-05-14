package x1

import (
	"go_analysis_test_example/core/dir3"
)

var F = func() {
	dir3.MyVar2.FFF()
}

func F1() {
	dir3.MyVar1.FFF()

	F()

	func() {
		dir3.MyVar1.FFF()
	}()
}
