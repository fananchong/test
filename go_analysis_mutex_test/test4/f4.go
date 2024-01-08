package test4

import (
	"fmt"
	"sync"
)

type A1 struct {
	sync.RWMutex // B
	B            string
}

func (a *A1) g1() {
	a.Lock()
	defer a.Unlock()
	fmt.Println(a.B)
}

func (a *A1) g2() {
	fmt.Println(a.B)
}

func F4() {
	// m1.Lock()
	// defer m1.Unlock()
	// a = make([]int, 0)

	a1 := A1{}
	fmt.Println(a1.B)
	a1.g1()
	a1.g2()
}
