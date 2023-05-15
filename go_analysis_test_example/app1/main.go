package main

import (
	"go_analysis_test_example/core/dir0"
	"go_analysis_test_example/core/dir3"
)

func F() {
	dir3.MyVar3.Walk(dir0.F1)
}

func main() {

	// dir4.SetHandler(2, func() {
	// 	dir0.F1()
	// })
	// // dir4.GetHandler(1)()

	// e := echo.New()
	// e.POST("xx", xx)

	F()
}

// func xx(c echo.Context) error {
// 	return nil
// }
