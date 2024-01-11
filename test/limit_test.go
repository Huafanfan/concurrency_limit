package test

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"git.garena.com/yifan.zhangyf/concurrency_limit/limits_core/limit"
	"git.garena.com/yifan.zhangyf/concurrency_limit/limits_core/limit/window"
	"git.garena.com/yifan.zhangyf/concurrency_limit/limits_core/limiter"
)

func TestLimit(t *testing.T) {
	//gradient2Limit := limit.NewGradient2Limit(1000, 10000, 100, limit.DefaultQueueSize, limit.DefaultSmoothing, limit.DefaultTolerance, limit.DefaultLongWindow)
	aimdLimit := limit.NewAIMDLimit(1000, limit.DefaultBackoffRatio, 10*time.Millisecond.Nanoseconds(), 10000, 100)
	windowLimiter := limit.NewWindowedLimit(aimdLimit, 100*time.Millisecond.Nanoseconds(), 200*time.Millisecond.Nanoseconds(), 5, 10*time.Millisecond.Nanoseconds(), window.NewPercentileSampleWindowFactory(0.9, 5))
	var inflight atomic.Int64
	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(index int) {
			climit := windowLimiter.GetLimit()
			if int(inflight.Load()) < climit {
				inflight.Add(1)
				now := time.Now()
				time.Sleep(30 * time.Millisecond)
				//time.Sleep(10 * time.Millisecond)
				//time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
				during := time.Since(now)
				inflight.Add(-1)
				windowLimiter.OnSample(now.UnixNano(), during.Nanoseconds(), int(inflight.Load()), false)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	println(windowLimiter.GetLimit())
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(index int) {
			climit := windowLimiter.GetLimit()
			if int(inflight.Load()) < climit {
				inflight.Add(1)
				now := time.Now()
				time.Sleep(15 * time.Millisecond)
				//time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
				during := time.Since(now)
				inflight.Add(-1)
				windowLimiter.OnSample(now.UnixNano(), during.Nanoseconds(), int(inflight.Load()), false)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	println(windowLimiter.GetLimit())

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(index int) {
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			climit := windowLimiter.GetLimit()
			if int(inflight.Load()) < climit {
				inflight.Add(1)
				now := time.Now()
				time.Sleep(100 * time.Millisecond)
				during := time.Since(now)
				inflight.Add(-1)
				windowLimiter.OnSample(now.UnixNano(), during.Nanoseconds(), int(inflight.Load()), true)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	println(windowLimiter.GetLimit())
}

func TestLimiter(t *testing.T) {
	gradient2Limit := limit.NewGradient2Limit(1000, 10000, 100, limit.DefaultQueueSize, limit.DefaultSmoothing, limit.DefaultTolerance, limit.DefaultLongWindow)
	windowLimiter := limit.NewWindowedLimit(gradient2Limit, 100*time.Millisecond.Nanoseconds(), 200*time.Millisecond.Nanoseconds(), 5, 10*time.Millisecond.Nanoseconds(), window.NewPercentileSampleWindowFactory(0.9, 5))
	simpleLimiter := limiter.NewSimpleLimiter(windowLimiter)
	if listener := simpleLimiter.Acquire(); listener != nil {
		simpleListener := limiter.NewSimpleListener(listener, simpleLimiter)
		// biz
		time.Sleep(time.Second)
		// success
		simpleListener.OnSuccess(time.Now().UnixNano())

		// ignore
		// simpleListener.OnIgnore(time.Now().UnixNano())

		// err
		// simpleListener.OnDropped(time.Now().UnixNano())
	} else {
		// concurrency limit
	}
}
