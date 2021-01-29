package goroutine_test

import (
	"context"
	"sync"
	"testing"
	"time"
)

func f1(ctx context.Context, wait *sync.WaitGroup) {
	wait.Done()
	a := 1
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			a++
		}
	}
}

func f2(ctx context.Context, wait *sync.WaitGroup) {
	a := 1
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			a++
			if a > 30 {
				wait.Done()
			}
		}
	}
}

var num = 70000

func Benchmark2(b *testing.B) {
	ctx, cancal := context.WithCancel(context.Background())
	wait := &sync.WaitGroup{}
	for i := 0; i < num; i++ {
		wait.Add(1)
		go f2(ctx, wait)
	}
	wait.Wait()
	cancal()
}

// Benchmark1 Benchmark1
func Benchmark1(b *testing.B) {
	ctx, cancal := context.WithCancel(context.Background())
	wait := &sync.WaitGroup{}
	for i := 0; i < num; i++ {
		wait.Add(1)
		go f1(ctx, wait)
	}
	wait.Wait()
	cancal()
}
