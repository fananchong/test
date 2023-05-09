package core

import "fmt"

func Foo2() {
	fmt.Println("")
}

func Bar2() {
	C2()
}

func C2() {

	f2 := func() {

	}

	f2()

	func() {

		func() {

		}()
	}()

	x := &X2{}
	x.FFF()
}

type X2 struct {
}

func (*X2) FFF() {
	Foo2()
}
