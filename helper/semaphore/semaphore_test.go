package semaphore

import (
	"testing"
	"time"
)

func TestSemaphore(t *testing.T) {
	s := New(2)
	timer := time.AfterFunc(time.Second, func() {
		panic("deadlock")
	})
	defer timer.Stop()

	s.Acquire()
	if !s.TryAcquire() {
		t.Fatalf("should acquire")
	}
	if s.TryAcquire() {
		t.Fatalf("should not acquire")
	}
	s.Release()
	s.Release()

	// This release should panic
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("should panic")
		}
	}()
	s.Release()
}
