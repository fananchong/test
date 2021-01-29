package common

import (
	"container/heap"
	"context"
	"time"
)

// TimerJob TimerJob
type TimerJob interface {
	GetInterval() time.Duration
	Check() bool
	Fire()
	Do() error
	GetChan() <-chan struct{}
	DebugInfo() interface{} // error 时，用来打印日志
}

// TimerJobDispather 定时任务分发器
type TimerJobDispather struct {
	ctx   context.Context
	queue PriorityQueue
	c     chan TimerJob
}

// NewTimerJobDispather 构造函数
func NewTimerJobDispather(ctx context.Context) *TimerJobDispather {
	return &TimerJobDispather{
		ctx: ctx,
		c:   make(chan TimerJob),
	}
}

// Add 增加定时任务
func (dispather *TimerJobDispather) Add(timerJob TimerJob) {
	dispather.c <- timerJob
}

// Run Run
func (dispather *TimerJobDispather) Run() {
	ticker := time.NewTicker(25 * time.Millisecond)
	for {
		select {
		case <-dispather.ctx.Done():
			return
		case timerJob := <-dispather.c:
			if len(dispather.queue) == 0 {
				dispather.queue = append(dispather.queue, &PriorityItem{
					Value:    timerJob,
					Priority: time.Now().Add(timerJob.GetInterval()).UnixNano(),
				})
				heap.Init(&dispather.queue)
			} else {
				heap.Push(&dispather.queue, &PriorityItem{
					Value:    timerJob,
					Priority: time.Now().Add(timerJob.GetInterval()).UnixNano(),
				})
			}
		case now := <-ticker.C:
			for dispather.queue.Len() > 0 {
				v := dispather.queue[0]
				if v.Priority <= now.UnixNano() {
					timerJob := heap.Pop(&dispather.queue).(*PriorityItem)
					timerJob.Value.Fire() // XXXX 阻塞问题
					timerJob.Priority = timerJob.Priority + int64(timerJob.Value.GetInterval())
					heap.Push(&dispather.queue, timerJob)
				} else {
					break
				}
			}
		}
	}
}
