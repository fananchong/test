package examples

import (
	"testing"

	gofail "github.com/etcd-io/gofail/runtime"
)

func TestExampleFunc(t *testing.T) {
	gofail.Enable("examples/ExampleString", `return("testtesttest")`)
	defer gofail.Disable("examples/ExampleString")
	if ExampleFunc() != "testtesttest" {
		t.Fatal("!!!")
	}
}

func TestExampleFunc2(t *testing.T) {
	gofail.Enable("examples/ExampleString", `sleep(1000)->100.0%panic`)
	defer gofail.Disable("examples/ExampleString")
	if ExampleFunc() != "example" {
		t.Fatal("!!!")
	}
}

func TestExampleOneLineFunc(t *testing.T) {
	gofail.Enable("examples/ExampleOneLine", `return`)
	defer gofail.Disable("examples/ExampleOneLine")
	if ExampleOneLineFunc() != "abc" {
		t.Fatal("!!!")
	}
}

func TestExampleOneLineFunc2(t *testing.T) {
	gofail.Enable("examples/ExampleOneLine", `100.0%panic`)
	defer gofail.Disable("examples/ExampleOneLine")
	if ExampleOneLineFunc() != "abc" {
		t.Fatal("!!!")
	}
}

func TestExampleLabelsFunc(t *testing.T) {
	gofail.Enable("examples/ExampleLabels", `return`)
	defer gofail.Disable("examples/ExampleLabels")
	if ExampleLabelsFunc() != "ijijijijij" {
		t.Fatal("!!!")
	}
}

func TestExampleLabelsGoFunc(t *testing.T) {
	gofail.Enable("examples/ExampleLabelsGo", `return`)
	defer gofail.Disable("examples/ExampleLabelsGo")
	if ExampleLabelsGoFunc() != "ijijijijij" {
		t.Fatal("!!!")
	}
}
