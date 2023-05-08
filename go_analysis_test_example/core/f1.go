package core

func Foo() {
}

func Bar() {
	C()
}

func C() {

	f := func() {
		Foo()
	}

	f()

	func() {
		Foo()
		func() {
			x := &X1{}
			x.Y.GGG()
			x.Y.Z.HHH()
			x.Y.Z.Hello()
			x.Y.Z.A1.Hello()
		}()
	}()

	x := &X1{}
	x.FFF()
	x.Y.GGG()
	x.Y.Z.HHH()
	x.Y.Z.Hello()
	x.Y.Z.A1.Hello()
}

type X1 struct {
	Y Y1
}

func (*X1) FFF() {
	Foo()
}

type Y1 struct {
	Z Z1
}

func (*Y1) GGG() {
	Foo()
}

type Z1 struct {
	*A1
}

func (*Z1) HHH() {
	Foo()
}

type A1 struct {
}

func (*A1) Hello() {

}
