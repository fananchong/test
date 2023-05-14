package x1

import (
	"go_analysis_test_example/core/dir3"
)

var f = func() {
	dir3.MyVar2.FFF()
}

func F1() {
	dir3.MyVar1.FFF()

	f()

	func() {
		dir3.MyVar1.FFF()
	}()
}
