package main

func foo2() {
}

func bar2() {
	baz2()
}

func baz2() {

	f := func() {

	}

	f()

	func() {

		func() {

		}()
	}()
}

type xxxx struct {
}

func (*xxxx) fff() {

}
