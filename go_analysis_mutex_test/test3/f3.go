package test3

import (
	"fmt"
	"sync"
)

var (
	m1 sync.RWMutex // a
	a  map[int32]interface{}
)

var flag bool

var f = func(func()) {}

func F3() {

	if flag {

		f(func() {
			m1.RLock()
			defer m1.RUnlock()

			if !flag {
				for k, v := range a {
					fmt.Println(k, v)
				}
			}
		})

	}

	fmt.Println(a)

	if flag {

		m1.RLock()

		if !flag {
			for k, v := range a {
				fmt.Println(k, v)
			}
		}
		m1.RUnlock()
	}

	fmt.Println(a)
}

func F32() {
	m1.Lock()
	m1.Unlock()
	fmt.Println(a)
}

func F33() {
	m1.Lock()
	m1.Unlock()

	m1.Lock()
	fmt.Println(a)
	m1.Unlock()
}

func F34() {
	F32()
}
