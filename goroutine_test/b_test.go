package goroutine_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/fananchong/test/goroutine_test/common"
)

type myTimerJob struct {
	c  chan struct{}
	v  int
	no int
}

func newMyTimerJob(no int) *myTimerJob {
	return &myTimerJob{
		c:  make(chan struct{}, 1),
		no: no,
	}
}

func (job *myTimerJob) GetInterval() time.Duration {
	return 1 * time.Second
}

func (job *myTimerJob) Check() bool {
	return true
}
func (job *myTimerJob) Fire() {
	fmt.Printf("no=%d Fire\n", job.no)
	job.c <- struct{}{}
}
func (job *myTimerJob) Do() error {
	fmt.Printf("no=%d Do\n", job.no)
	job.v++
	return nil
}
func (job *myTimerJob) GetChan() <-chan struct{} {
	return job.c
}
func (job *myTimerJob) DebugInfo() interface{} {
	return &job.v
}

func f3(ctx context.Context, wait *sync.WaitGroup, timerJobObj common.TimerJob) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-timerJobObj.GetChan():
			timerJobObj.Do()
			if *timerJobObj.DebugInfo().(*int) > 2 {
				wait.Done()
				return
			}
		}
	}
}

var num3 = 2

func Benchmark3(b *testing.B) {
	ctx, cancal := context.WithCancel(context.Background())
	wait := &sync.WaitGroup{}
	dispather := common.NewTimerJobDispather(ctx)
	go dispather.Run()
	for i := 0; i < num3; i++ {
		wait.Add(1)
		job := newMyTimerJob(i)
		dispather.Add(job)
		go f3(ctx, wait, job)
	}
	wait.Wait()
	cancal()
}
