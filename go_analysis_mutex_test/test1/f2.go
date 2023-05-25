package test1

import "fmt"

func F2() {
	func() {
		m1.Lock()
		defer m1.Unlock()
		_ = fmt.Sprintf("%v", a)
	}()
	_ = fmt.Sprintf("%v", b)
	_ = fmt.Sprintf("%v", c)
}
