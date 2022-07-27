package main

import (
	"fmt"
)

type mysturct struct {
	i int
	j float32
	k string
	g []string
	m map[string]int
}

var gA = &mysturct{
	i: 1,
	j: 2.999,
	k: "kkk",
	g: []string{"a", "b", "c"},
	m: map[string]int{
		"mk1": 10,
		"mk2": 20,
		"mk3": 30,
	},
}

func main() {
	fmt.Println(gA)

	b := &mysturct{
		i: 90,
		j: 10.8,
		k: "op900",
		g: []string{"ad3", "b33", "cce"},
		m: map[string]int{
			"mk090": 100,
			"mk223": 270,
			"mk455": 309,
		},
	}
	fmt.Println(b)
	for i := 0; i < 500; i++ {
		fmt.Println("looping")
		if i == 120 {
			panic(fmt.Errorf("i:%v", i))
		}
	}
	fmt.Println("Done")
}
