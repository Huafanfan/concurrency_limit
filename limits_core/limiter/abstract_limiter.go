package limiter

import (
	"sync/atomic"
	"time"

	"github.com/Huafanfan/concurrency_limit/limits_core"
)

type AbstractLimiter struct {
	limits_core.Limiter

	inFlight       atomic.Int64
	limitAlgorithm limits_core.Limit
	limit          atomic.Int64
}

func NewAbstractLimiter(limit limits_core.Limit) *AbstractLimiter {
	a := &AbstractLimiter{
		limitAlgorithm: limit,
	}
	a.limit.Store(int64(a.limitAlgorithm.GetLimit()))
	return a
}

func (aLimiter *AbstractLimiter) CreateRejectedListener() limits_core.Listener {
	return nil
}

func (aLimiter *AbstractLimiter) CreateListener() limits_core.Listener {
	l := &AbstractListener{
		StartTime:       time.Now().UnixNano(),
		CurrentInflight: int(aLimiter.inFlight.Add(1)),
		abstractLimiter: aLimiter,
	}

	return l
}

func (aLimiter *AbstractLimiter) GetLimit() int {
	return int(aLimiter.limit.Load())
}

func (aLimiter *AbstractLimiter) GetInflight() int {
	return int(aLimiter.inFlight.Load())
}

func (aLimiter *AbstractLimiter) OnNewLimit(newLimit int) {
	aLimiter.limit.Store(int64(newLimit))
}

type AbstractListener struct {
	StartTime       int64
	CurrentInflight int
	abstractLimiter *AbstractLimiter
}

func (aListener *AbstractListener) OnSuccess(endTime int64) {
	aListener.abstractLimiter.inFlight.Add(-1)
	aListener.abstractLimiter.limitAlgorithm.OnSample(aListener.StartTime, endTime-aListener.StartTime, aListener.CurrentInflight, false)
}
func (aListener *AbstractListener) OnIgnore(endTime int64) {
	aListener.abstractLimiter.inFlight.Add(-1)
}

func (aListener *AbstractListener) OnDropped(endTime int64) {
	aListener.abstractLimiter.inFlight.Add(-1)
	aListener.abstractLimiter.limitAlgorithm.OnSample(aListener.StartTime, endTime-aListener.StartTime, aListener.CurrentInflight, true)
}
