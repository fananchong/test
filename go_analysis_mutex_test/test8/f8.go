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
