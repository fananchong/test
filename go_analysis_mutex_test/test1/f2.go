package test1

import "fmt"

func f2() {
	{
		m1.Lock()
		defer m1.Unlock()
		fmt.Println(a)
	}
	b = make(map[int]int)
	c = "22"
}
