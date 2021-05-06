package main

import (
	"fmt"
	"time"
)

var counter = 0

func f1() {
	counter++
}

func f2() {
	counter++
}

func main() {
	fmt.Println("hello world")
	go f1()
	go f2()
	time.Sleep(1 * time.Second)
}
