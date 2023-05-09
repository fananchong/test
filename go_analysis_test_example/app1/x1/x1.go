package x1

import "go_analysis_test_example/core"

var f = func() {
	core.MyVar1.FFF()
}

func F1() {
	core.MyVar1.FFF()

	f()

	func() {
		core.MyVar1.FFF()
	}()
}
