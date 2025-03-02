package limit

import (
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	"github.com/Huafanfan/concurrency_limit/limits_core"
)

var (
	DefaultTimeout      = 5 * time.Second.Nanoseconds()
	DefaultBackoffRatio = 0.9
)

type AIMDLimit struct {
	*AbstractLimit

	BackoffRatio float64
	Timeout      int64
	MaxLimit     int
	MinLimit     int

	mu sync.RWMutex
}

func NewAIMDLimit(initialLimit int, backoffRatio float64, timeout int64, maxLimit int, minLimit int) limits_core.Limit {
	a := &AIMDLimit{
		BackoffRatio: backoffRatio,
		Timeout:      timeout,
		MaxLimit:     maxLimit,
		MinLimit:     minLimit,
		mu:           sync.RWMutex{},
	}
	a.AbstractLimit = NewAbstractLimit(initialLimit)
	a.AbstractLimit._Update = a._Update
	return a
}

func (a *AIMDLimit) _Update(startTime int64, rtt int64, inflight int, didDrop bool) int {
	currentLimit := a.GetLimit()

	if didDrop || rtt > a.Timeout {
		currentLimit = int(float64(currentLimit) * a.BackoffRatio)
	} else if inflight*2 >= currentLimit {
		currentLimit = currentLimit + 1
	}

	log.Printf(fmt.Sprintf("New limit=%v", currentLimit))

	return int(math.Min(float64(a.MaxLimit), math.Max(float64(a.MinLimit), float64(currentLimit))))
}
