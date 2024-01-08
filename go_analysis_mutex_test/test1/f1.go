package test1

import (
	"fmt"
	"sync"
)

var m1 sync.Mutex // a,b,c
var a int
var b = map[int]int{}
var c string

func F1() {
	m1.Lock()
	defer m1.Unlock()
	a++
	b[1] = 1
	c = "111"
}

func F2() {
	func() {
		m1.Lock()
		defer m1.Unlock()
		_ = fmt.Sprintf("%v", a)

		go func() {
			c = "222"
		}()

	}()
	b = map[int]int{}
}
