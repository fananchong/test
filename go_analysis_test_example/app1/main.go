package main

import (
	"go_analysis_test_example/core/dir0"
	"go_analysis_test_example/core/dir4"
)

func main() {

	dir4.SetHandler(2, func() {
		dir0.F1()
	})
	// dir4.GetHandler(1)()
}
