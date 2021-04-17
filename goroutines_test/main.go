package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var counter = int64(0)

func main() {
	N := 6397677
	t1 := time.Now()
	wait := &sync.WaitGroup{}
	for i := 0; i < N; i++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			atomic.AddInt64(&counter, 1)
		}()
	}
	wait.Wait()
	t2 := time.Now()
	fmt.Printf("N=%d cost=%v avg=%v\n", N, t2.Sub(t1), t2.Sub(t1)/time.Duration(N))
}
