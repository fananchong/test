package test11

import (
	"errors"
	"sync"
)

type xxx11 struct {
	sync.RWMutex // A
	A            map[int]int
}
type A11 struct {
	sync.RWMutex // xxx11
	xxx11
}

var a = &A11{}

func F113() map[int]int {
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

func F115() map[int]int {
	m2.Lock()
	defer m2.Unlock()
	return c.A // nolint: mutex_check
}

var m3 sync.Mutex // d
var d int

type E struct {
	k  uint64
	Ds []int
}

func F116() (int, *E, error) {
	m3.Lock()
	defer m3.Unlock()
	return d, &E{}, errors.New("text string")
}

var m4 sync.Mutex // e
var e func()

func F117() func() {
	m4.Lock()
	defer m4.Unlock()
	return e
}

var m5 sync.Mutex // f
var f error

func F118() error {
	m5.Lock()
	defer m5.Unlock()
	return f
}
