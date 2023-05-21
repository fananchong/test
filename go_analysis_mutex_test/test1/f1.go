package test1

import "sync"

var m1 sync.Mutex // a,b,c
var a int
var b = map[int]int{}
var c string

func f1() {
	m1.Lock()
	defer m1.Unlock()
	a++
	b[1] = 1
	c = "111"
}
