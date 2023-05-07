package core

func Foo() {
}

func Bar() {
	C()
}

func C() {

	f := func() {

	}

	f()

	func() {

		func() {

		}()
	}()

	x := &X1{}
	x.FFF()
}

type X1 struct {
}

func (*X1) FFF() {
	Foo()
}
