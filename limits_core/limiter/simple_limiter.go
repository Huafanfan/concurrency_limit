package limiter

import (
	"context"
	"git.garena.com/yifan.zhangyf/concurrency_limit/limits_core"
	"github.com/marusama/semaphore/v2"
)

type SimpleLimiter struct {
	aLimiter  *AbstractLimiter
	semaphore semaphore.Semaphore
}

func NewSimpleLimiter(limit limits_core.Limit) *SimpleLimiter {
	s := &SimpleLimiter{}
	s.aLimiter = NewAbstractLimiter(limit)
	s.semaphore = semaphore.New(s.aLimiter.GetLimit())
	return s
}

func (sLimiter *SimpleLimiter) Acquire(ctx context.Context) limits_core.Listener {
	if !sLimiter.semaphore.TryAcquire(1) {
		return sLimiter.aLimiter.CreateRejectedListener()
	}
	listener := sLimiter.aLimiter.CreateListener()
	return NewSimpleListener(listener, sLimiter)
}

func (sLimiter *SimpleLimiter) OnNewLimit(newLimit int) {
	sLimiter.aLimiter.OnNewLimit(newLimit)
	sLimiter.semaphore.SetLimit(newLimit)
}

type SimpleListener struct {
	sLimiter *SimpleLimiter
	delegate limits_core.Listener
}

func NewSimpleListener(delegate limits_core.Listener, sl *SimpleLimiter) limits_core.Listener {
	a := &SimpleListener{
		delegate: delegate,
		sLimiter: sl,
	}
	return a
}

func (sListener *SimpleListener) OnSuccess(endTime int64) {
	sListener.delegate.OnSuccess(endTime)
	sListener.sLimiter.semaphore.Release(1)
	sListener.UpdateLimit()
}

func (sListener *SimpleListener) OnIgnore(endTime int64) {
	sListener.delegate.OnIgnore(endTime)
	sListener.sLimiter.semaphore.Release(1)
	sListener.UpdateLimit()
}

func (sListener *SimpleListener) OnDropped(endTime int64) {
	sListener.delegate.OnDropped(endTime)
	sListener.sLimiter.semaphore.Release(1)
	sListener.UpdateLimit()
}

func (sListener *SimpleListener) UpdateLimit() {
	if sListener.sLimiter.aLimiter.GetLimit() != sListener.sLimiter.aLimiter.limitAlgorithm.GetLimit() {
		sListener.sLimiter.OnNewLimit(sListener.sLimiter.aLimiter.limitAlgorithm.GetLimit())
	}
}
