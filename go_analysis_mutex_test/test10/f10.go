package test10

import (
	"fmt"
	"sync"
)

type xxx10 struct {
}
type A10 struct {
	sync.RWMutex // xxx10
	xxx10
}

func (a *A10) F10() {
	fmt.Println(a.xxx10)
}

var a = &A10{}

// func F102() {
// 	a.Lock()
// 	a.F10() // nolint: mutex_check
// 	a.Unlock()
// 	a.F10()
// }

func F103() {
	a.Lock()
	a.F10()
	a.Unlock()
	a.F10() // nolint: mutex_check
}

// func F104() {
// 	a.Lock()
// 	a.F10()
// 	a.Unlock()
// 	a.F10()
// }
