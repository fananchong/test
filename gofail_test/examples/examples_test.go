package examples

import (
	"testing"

	gofail "github.com/etcd-io/gofail/runtime"
)

func TestWhatever(t *testing.T) {
	gofail.Enable("examples/ExampleString", `return("testtesttest")`)
	defer gofail.Disable("examples/ExampleString")
	if ExampleFunc() != "testtesttest" {
		t.Fatal("!!!")
	}
}
