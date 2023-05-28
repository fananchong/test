package test4

import (
	"fmt"
	"sync"
)

type xxx struct {
}

type A1 struct {
	sync.RWMutex // xxx
	xxx
}

func (a *A1) g1() {
	a.Lock()
	defer a.Unlock()
	fmt.Println(a.xxx)
}

func (a *A1) g2() {
	fmt.Println(a.xxx)
}

func F5() {
	a1 := A1{}
	fmt.Println(a1.xxx)
	a1.g1()
	a1.g2()
}

func init() {
	F5()
}
