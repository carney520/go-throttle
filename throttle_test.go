package throttle

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func times(n int, dur time.Duration, fun func()) {
	for i := 0; i < n; i++ {
		fun()
		if dur > 0 {
			time.Sleep(dur)
		}
	}
}

func ExampleNew() {
	tr := New(100*time.Millisecond, func() {
		fmt.Println("some event")
	})
	defer tr.Stop()
	times(6, 40*time.Millisecond, tr.Trigger)
	time.Sleep(100 * time.Microsecond)
	// Output: some event
	// some event
	// some event
}

func TestThrottleEqualTime(t *testing.T) {
	var count int32
	tr := New(100*time.Millisecond, func() {
		atomic.AddInt32(&count, 1)
	})
	defer tr.Stop()
	times(5, 150*time.Millisecond, tr.Trigger)
	time.Sleep(100 * time.Microsecond)

	if c := atomic.LoadInt32(&count); c != 5 {
		t.Errorf("throttle invoked error: expect 5 but got %d\n", c)
	}
}

func TestThrottleUnequalTime(t *testing.T) {
	var count int32
	tr := New(100*time.Millisecond, func() {
		atomic.AddInt32(&count, 1)
	})
	defer tr.Stop()
	times(3, 150*time.Millisecond, func() {
		times(3, 0, tr.Trigger)
	})
	time.Sleep(100 * time.Microsecond)

	if c := atomic.LoadInt32(&count); c != 3 {
		t.Errorf("throttle invoked error: expect 3 but got %d\n", c)
	}
}

func TestThrottleStop(t *testing.T) {
	var count int32
	tr := New(100*time.Millisecond, func() {
		atomic.AddInt32(&count, 1)
	})
	tr.Trigger()
	time.Sleep(120 * time.Millisecond)
	tr.Stop()
	times(3, 100*time.Millisecond, tr.Trigger)
	if c := atomic.LoadInt32(&count); c != 1 {
		t.Errorf("throttle invoked error: expect 1 but got %d\n", c)
	}
}
