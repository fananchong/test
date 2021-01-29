package common

// An PriorityItem is something we manage in a priority queue.
type PriorityItem struct {
	Value    TimerJob // The value of the item; arbitrary.
	Priority int64    // The priority of the item in the queue.
}

// An PriorityQueue is a min-heap of ints.
type PriorityQueue []*PriorityItem

func (h PriorityQueue) Len() int { return len(h) }
func (h PriorityQueue) Less(i, j int) bool {
	return h[i].Priority < h[j].Priority
}
func (h PriorityQueue) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

// Push Push
func (h *PriorityQueue) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(*PriorityItem))
}

// Pop Pop
func (h *PriorityQueue) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
