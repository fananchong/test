package test5

import (
	"sync"
)

var (
	m1 sync.RWMutex // a
	a  []string
)

func F51() {
	m1.Lock()
	defer m1.Unlock()
	go func(aa []string) {}(a)

}

func F52() {
	m1.Lock()
	defer m1.Unlock()
	go func(aa []string) {

	}(a)
}

func F53() {
	m1.Lock()
	defer m1.Unlock()
	go func() {
		// a = append(a, "1")
		// a = make([]string, 0)

		func(aa []string) {}(a)
		// func(aa []string) {

		// }(a)

		// fmt.Println(a)
	}()
}
