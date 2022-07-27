package main

import (
	"fmt"
)

func main() {
	for i := 0; i < 500; i++ {
		fmt.Println("looping")
		if i == 120 {
			panic(fmt.Errorf("i:%v", i))
		}
	}
	fmt.Println("Done")
}
