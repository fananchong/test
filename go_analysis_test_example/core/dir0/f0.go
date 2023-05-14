package dir0

import "go_analysis_test_example/core/dir3"

func F1() {
	dir3.MyVar1.FFF()
}

func F2() {
	dir3.MyVar2.FFF()
}

func F3() {
	dir3.MyVar1.FFF()
}
