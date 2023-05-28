package test3

import (
	"fmt"
	"sync"
)

var (
	m1 sync.RWMutex // a
	a  []string
)

func F3() {
	m1.Lock()
	defer m1.Unlock()
	a = make([]string, 0)
	// fmt.Println(a)

	fmt.Println("aaa")
}

func F32() {
	m1.Lock()
	m1.Unlock()
	fmt.Println(a)
}

func F33() {
	m1.Lock()
	m1.Unlock()

	m1.Lock()
	fmt.Println(a)
	m1.Unlock()
}

func F34() {
	F32()
}
