package limit

import (
	"sync"
)

type AbstractLimit struct {
	Limit   int
	_Update func(startTime int64, rtt int64, inflight int, didDrop bool) int
	mu      sync.RWMutex
}

func NewAbstractLimit(initialLimit int) *AbstractLimit {
	return &AbstractLimit{
		Limit: initialLimit,
	}
}

func (l *AbstractLimit) GetLimit() int {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.Limit
}

func (l *AbstractLimit) OnSample(startTime int64, rtt int64, inflight int, didDrop bool) {
	l.SetLimit(l._Update(startTime, rtt, inflight, didDrop))
}

func (l *AbstractLimit) SetLimit(newLimit int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if newLimit != l.Limit {
		l.Limit = newLimit
	}
}
