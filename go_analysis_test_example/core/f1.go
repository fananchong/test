package core

func Foo() {
}

func Bar() {
	Baz()
}

func Baz() {

	f := func() {

	}

	f()

	func() {

		func() {

		}()
	}()

	x := &X2{}
	x.FFF()
}
