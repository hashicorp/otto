package localaddr

import (
	"container/heap"
	"net"
	"time"
)

// ipEntry is an entry in the priority queue of leased IP addresses.
type ipEntry struct {
	LeaseTime time.Time
	Value     net.IP

	// This is maintained by heap.Push and heap.Pop
	Index int
}

// ipQueue is an implementation of heap.Interface and holds ipEntrys.
type ipQueue []*ipEntry

func (q ipQueue) Len() int { return len(q) }

func (q ipQueue) Less(i, j int) bool {
	return q[i].LeaseTime.Before(q[j].LeaseTime)
}

func (q ipQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].Index = i
	q[j].Index = j
}

func (q *ipQueue) Push(x interface{}) {
	item := x.(*ipEntry)
	item.Index = len(*q)
	*q = append(*q, item)
}

func (q *ipQueue) Pop() interface{} {
	old := *q
	n := len(old)
	item := old[n-1]
	item.Index = -1
	*q = old[0 : n-1]
	return item
}

// Update updates the given entry. This entry must already be in the queue.
func (q *ipQueue) Update(item *ipEntry) {
	heap.Fix(q, item.Index)
}
