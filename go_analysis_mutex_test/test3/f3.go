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
	// m1.Lock()
	// defer m1.Unlock()
	a = make([]string, 0)
	// fmt.Println(a)

	fmt.Println("aaa")
}

func init() {
	F3()
}
