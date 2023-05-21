package test2

import "sync"

type A1 struct {
	M sync.RWMutex // A
	A int
}

func (a1 *A1) f1() {
	a1.M.RLock()
	defer a1.M.RUnlock()
	a1.A = 1
}

func f2() {
	var a1 A1
	a1.f1()
	a1.A = 2
}

func init() {
	f2()
	f3()
}

func f3() {
	var m2 sync.Mutex
	m2.Lock()
	defer m2.Unlock()
}
