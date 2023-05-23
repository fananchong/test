package test1

import "fmt"

func f2() {
	func() {
		m1.Lock()
		defer m1.Unlock()
		fmt.Println(a)
	}()
	fmt.Println(b)
	fmt.Println(c)
}
