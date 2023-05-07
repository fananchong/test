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
}

type X2 struct {
}

func (*X2) FFF() {
	Baz2()
}
