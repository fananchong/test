package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

func test_800b() {
	N := 1000000
	t1 := time.Now()
	wait := &sync.WaitGroup{}
	for i := 0; i < N/100; i++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			var a [100]int64
			for j := 0; j < 100; j++ {
				atomic.AddInt64(&a[len(a)-j-1], 1)
				runtime.Gosched()
			}
		}()
	}
	wait.Wait()
	t2 := time.Now()
	fmt.Printf("800b N=%d cost=%v avg=%v\n", N, t2.Sub(t1), t2.Sub(t1)/time.Duration(N))
}

func test_8k() {
	N := 1000000
	t1 := time.Now()
	wait := &sync.WaitGroup{}
	for i := 0; i < N/100; i++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			var a [1024]int64
			for j := 0; j < 100; j++ {
				atomic.AddInt64(&a[len(a)-j-1], 1)
				runtime.Gosched()
			}
		}()
	}
	wait.Wait()
	t2 := time.Now()
	fmt.Printf("8k N=%d cost=%v avg=%v\n", N, t2.Sub(t1), t2.Sub(t1)/time.Duration(N))
}

func test_80k() {
	N := 1000000
	t1 := time.Now()
	wait := &sync.WaitGroup{}
	for i := 0; i < N/100; i++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			var a [1024 * 10]int64
			for j := 0; j < 100; j++ {
				atomic.AddInt64(&a[len(a)-j-1], 1)
				runtime.Gosched()
			}
		}()
	}
	wait.Wait()
	t2 := time.Now()
	fmt.Printf("80k N=%d cost=%v avg=%v\n", N, t2.Sub(t1), t2.Sub(t1)/time.Duration(N))
}

func test_800k() {
	N := 1000000
	t1 := time.Now()
	wait := &sync.WaitGroup{}
	for i := 0; i < N/100; i++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			var a [1024 * 100]int64
			for j := 0; j < 100; j++ {
				atomic.AddInt64(&a[len(a)-j-1], 1)
				runtime.Gosched()
			}
		}()
	}
	wait.Wait()
	t2 := time.Now()
	fmt.Printf("800k N=%d cost=%v avg=%v\n", N, t2.Sub(t1), t2.Sub(t1)/time.Duration(N))
}

func test_4m() {
	N := 1000000
	t1 := time.Now()
	wait := &sync.WaitGroup{}
	for i := 0; i < N/100; i++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			var a [1024 * 500]int64
			for j := 0; j < 100; j++ {
				atomic.AddInt64(&a[len(a)-j-1], 1)
				runtime.Gosched()
			}
		}()
	}
	wait.Wait()
	t2 := time.Now()
	fmt.Printf("4m N=%d cost=%v avg=%v\n", N, t2.Sub(t1), t2.Sub(t1)/time.Duration(N))
}

func main() {
	test_800b()
	test_8k()
	test_80k()
	test_800k()
	test_4m()
}
