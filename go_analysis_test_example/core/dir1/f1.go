package dir1

import "fmt"

type X1 struct {
}

func (*X1) FFF() {
	fmt.Println("")
}

type A1 struct {
}

func (*A1) FA() {
	fmt.Println("")
}

type B1 struct {
	*A1
}
