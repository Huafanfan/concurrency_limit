package limiter

import (
	"git.garena.com/yifan.zhangyf/concurrency_limit/limits_core"
	"github.com/marusama/semaphore/v2"
)

type SimpleLimiter struct {
	alimiter  *AbstractLimiter
	semaphore semaphore.Semaphore
}

func NewSimpleLimiter(limit limits_core.Limit) *SimpleLimiter {
	s := &SimpleLimiter{}
	s.alimiter = NewAbstractLimiter(limit)
	s.semaphore = semaphore.New(s.alimiter.GetLimit())
	return s
}

func (sLimiter *SimpleLimiter) Acquire() limits_core.Listener {
	if !sLimiter.semaphore.TryAcquire(1) {
		return sLimiter.alimiter.CreateRejectedListener()
	}
	return sLimiter.alimiter.CreateListener()
}

func (sLimiter *SimpleLimiter) OnNewLimit(newLimit int) {
	oldLimit := sLimiter.alimiter.GetLimit()
	sLimiter.alimiter.OnNewLimit(newLimit)

	if newLimit > oldLimit {
		sLimiter.semaphore.SetLimit(newLimit - oldLimit)
	} else {
		sLimiter.semaphore.SetLimit(oldLimit - newLimit)
	}
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
	if sListener.sLimiter.alimiter.GetLimit() != sListener.sLimiter.alimiter.limitAlgorithm.GetLimit() {
		sListener.sLimiter.OnNewLimit(sListener.sLimiter.alimiter.limitAlgorithm.GetLimit())
	}
}
