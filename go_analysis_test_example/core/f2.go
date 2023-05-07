package core

func Foo2() {
}

func Bar2() {
	Baz2()
}

func Baz2() {

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
