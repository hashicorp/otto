package localaddr

import (
	"container/heap"
	"net"
	"testing"
	"time"
)

func TestIPQueue_impl(t *testing.T) {
	var _ heap.Interface = new(ipQueue)

	now := time.Now().UTC()
	q := ipQueue(make([]*ipEntry, 0))
	q = append(q, &ipEntry{LeaseTime: now, Value: net.ParseIP("1.2.3.4")})
	q = append(q, &ipEntry{LeaseTime: now.Add(-1 * time.Minute), Value: net.ParseIP("2.3.4.5")})
	q = append(q, &ipEntry{LeaseTime: now.Add(-3 * time.Minute), Value: net.ParseIP("1.2.3.5")})
	q = append(q, &ipEntry{LeaseTime: now.Add(-2 * time.Minute), Value: net.ParseIP("1.2.3.6")})

	heap.Init(&q)
	actual := heap.Pop(&q)
	if actual.(*ipEntry).Value.String() != "1.2.3.5" {
		t.Fatalf("bad: %s", actual)
	}
}
