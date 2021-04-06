package main

import (
	"log"
	"runtime"
	"time"
)

// https://www.jianshu.com/p/43b0e12043b9

func f() {
	container := make([]int, 8)
	log.Println("> loop.")
	for i := 0; i < 32*1000*1000; i++ {
		container = append(container, i)
	}
	log.Println("< loop.")

}

func main() {
	log.Println("start.")
	f()

	log.Println("force gc.")
	runtime.GC()

	log.Println("done.")
	time.Sleep(1 * time.Hour)
}
