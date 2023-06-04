package test7

import (
	"fmt"
	"sync"
)

var m1 sync.Mutex // a
var a int

func F8() {
	// a++
	fmt.Println(a)
}

func F82() {
	m1.Lock()
	F8()
	m1.Unlock()
	F8()
	F8()
}

func F83() {
	m1.Lock()
	defer m1.Unlock()
	for i := 0; i < a+100; i++ {
		F8()
	}
}
