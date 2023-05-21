package test1

import "sync"

var m sync.Mutex // a,b,c
var a int
var b = map[int]int{}
var c string

func f1() {
	m.Lock()
	defer m.Unlock()
	a++
	b[1] = 1
	c = "111"
}
