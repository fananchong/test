package test6

import (
	"fmt"
	"sync"
)

type xxx struct {
}

type yyy struct {
}

type A1 struct {
	sync.RWMutex // xxx,yyy
	xxx
	*yyy
}

func (a *A1) g1() {
	a.Lock()
	defer a.Unlock()
	fmt.Println(a.xxx)
}

func (a *A1) g2() {
	fmt.Println(a.yyy)
}

func F4() {
	// m1.Lock()
	// defer m1.Unlock()
	// a = make([]int, 0)

	a1 := A1{}
	fmt.Println(a1.xxx)
	a1.g1()
	a1.g2()
}

func init() {
	F4()
}
