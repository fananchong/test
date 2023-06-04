package test11

import (
	"sync"
)

type xxx11 struct {
	A int
}
type A11 struct {
	sync.RWMutex // xxx11
	xxx11
}

var a = &A11{}

func F113() int {
	a.Lock()
	defer a.Unlock()
	return a.xxx11.A
}

// func F114() xxx11 {
// 	fmt.Println("a")
// 	return a.xxx11
// }

var m1 sync.Mutex //b
var b int

func F114() int {
	m1.Lock()
	defer m1.Unlock()
	return b
}

var m2 sync.Mutex // c
var c xxx11

func F115() int {
	m2.Lock()
	defer m2.Unlock()
	return c.A
}
