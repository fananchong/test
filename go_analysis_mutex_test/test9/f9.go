package test7

import (
	"fmt"
	"sync"
)

type xxx9 struct {
}
type A9 struct {
	sync.RWMutex // xxx9
	xxx9
}

func (a *A9) F9() {
	fmt.Println(a.xxx9)
}

var a = &A9{}

// func F92() {
// 	a.Lock()
// 	a.F9()
// 	a.Unlock()
// 	a.F9()
// }

func F93() {
	a.Lock()
	defer a.Unlock()
	for i := 0; i < 100; i++ {
		a.F9()
	}
}
