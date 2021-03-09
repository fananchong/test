package main

import (
	"time"
)

func main() {
	// for i := 0; i < 10000; i++ {
	// 	time.NewTicker(10 * time.Millisecond)
	// 	time.NewTicker(25 * time.Millisecond)
	// }
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
		}
	}
}
