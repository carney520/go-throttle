// Package throttle create a Go version throttle.
// When the passed function invoked repeatedly, will only actually call the
// original function at most once per every wait `time.Duration`.
// Useful for rate-limiting events that occur faster than you can keep up with.
package throttle

import (
	"sync"
	"time"
)

// Throttler define throttle interface
type Throttler interface {
	// Trigger trigger an invocation
	Trigger()
	// Stop cancel throttle
	Stop()
}

type throttled struct {
	// protect follow fields
	cond    *sync.Cond
	period  time.Duration
	stoped  bool
	waiting bool
	last    time.Time
}

func (t *throttled) Trigger() {
	t.cond.L.Lock()
	defer t.cond.L.Unlock()
	if !t.waiting && !t.stoped {
		delta := time.Now().Sub(t.last)
		if delta > t.period {
			t.waiting = true
			t.cond.Broadcast()
		} else {
			t.waiting = true
			// 等待剩余时间后唤醒
			time.AfterFunc(t.period-delta, t.cond.Broadcast)
		}
	}
}

func (t *throttled) Stop() {
	if t.stoped {
		return
	}
	t.cond.L.Lock()
	defer t.cond.L.Unlock()
	t.stoped = true
	// 通知已经关闭
	t.cond.Broadcast()
}

func (t *throttled) next() (goon bool) {
	t.cond.L.Lock()
	defer t.cond.L.Unlock()
	for !t.waiting && !t.stoped {
		t.cond.Wait()
	}

	if !t.stoped {
		t.waiting = false
		t.last = time.Now()
	}

	return !t.stoped
}

// New create a throttle
func New(wait time.Duration, callback func()) Throttler {
	th := &throttled{
		cond:   sync.NewCond(&sync.Mutex{}),
		period: wait,
	}

	go func() {
		for th.next() {
			callback()
		}
	}()

	return th
}
