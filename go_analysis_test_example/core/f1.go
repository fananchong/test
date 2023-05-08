package core

func Foo() {
}

func Fxx() {
	Foo()
}

func Fyy() {
	Foo()
}

func Fhh() {
	Foo()
}

func Bar() {
	C()
}

func C() {

	d1 := func() {
		Fhh()
	}

	d1()

	d2 := Fxx
	d2()

	d3, d4 := func() {
		Foo()
	}, Fyy
	d3()

	d4()

	x := &X1{}

	d5 := x.Y.Z.A1.Hello
	d5()

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
