package test1

import "fmt"

func f2() {
	{
		m.Lock()
		defer m.Unlock()
		fmt.Println(a)
	}
	b = make(map[int]int)
	c = "22"
}
