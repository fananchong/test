package main

import (
	"go_workspace_test/core/dir1"

	"github.com/labstack/echo"
)

func main() {
	dir1.F1()
	_ = echo.New()
}
