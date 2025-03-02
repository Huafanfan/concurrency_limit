package test

import (
	"context"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Huafanfan/concurrency_limit/limits_core/limit"
	"github.com/Huafanfan/concurrency_limit/limits_core/limit/window"
	"github.com/Huafanfan/concurrency_limit/limits_core/limiter"
	"github.com/marusama/semaphore/v2"
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

func TestSimpleLimiter(t *testing.T) {
	//gradient2Limit := limit.NewGradient2Limit(5000, 20000, 100, limit.DefaultQueueSize, limit.DefaultGradientSmoothing, limit.DefaultTolerance, limit.DefaultLongWindow)
	vegasLimit := limit.NewVegasLimit(5000, 20000, limit.DefaultAlphaFunc, limit.DefaultBetaFunc, limit.DefaultIncreaseFunc, limit.DefaultDecreaseFunc, limit.DefaultThresholdFunc, limit.DefaultVegasSmoothing, limit.DefaultProbeMultiplier)

	windowLimiter := limit.NewWindowedLimit(vegasLimit, 5*time.Millisecond.Nanoseconds(), 20*time.Millisecond.Nanoseconds(), 5, 10*time.Microsecond.Nanoseconds(), window.NewPercentileSampleWindowFactory(0.9, 5))
	simpleLimiter := limiter.NewSimpleLimiter(windowLimiter)

	var wg sync.WaitGroup
	for i := 0; i < 60; i++ {
		for j := 0; j < 10000; j++ {
			wg.Add(1)
			go func() {
				if listener := simpleLimiter.Acquire(context.Background()); listener != nil {
					// biz
					time.Sleep(1 * time.Millisecond)
					// success
					listener.OnSuccess(time.Now().UnixNano())
				}
				wg.Done()
			}()
		}
		time.Sleep(time.Second)
	}
	wg.Wait()

	println(windowLimiter.GetLimit())

	for i := 0; i < 60; i++ {
		for j := 0; j < 10000; j++ {
			wg.Add(1)
			go func() {
				if listener := simpleLimiter.Acquire(context.Background()); listener != nil {
					simpleListener := limiter.NewSimpleListener(listener, simpleLimiter)
					// biz
					time.Sleep(10 * time.Millisecond)
					// success
					simpleListener.OnDropped(time.Now().UnixNano())
				}
				wg.Done()
			}()
		}
		time.Sleep(time.Second)
	}
	wg.Wait()

	println(windowLimiter.GetLimit())

	for i := 0; i < 60; i++ {
		for j := 0; j < 10000; j++ {
			wg.Add(1)
			go func() {
				if listener := simpleLimiter.Acquire(context.Background()); listener != nil {
					simpleListener := limiter.NewSimpleListener(listener, simpleLimiter)
					// biz
					time.Sleep(1 * time.Millisecond)
					// success
					simpleListener.OnSuccess(time.Now().UnixNano())
				}
				wg.Done()
			}()
		}
		time.Sleep(time.Second)
	}
	wg.Wait()

	println(windowLimiter.GetLimit())
}
func TestPartitionedLimiter(t *testing.T) {
	//gradient2Limit := limit.NewGradient2Limit(5000, 20000, 100, limit.DefaultQueueSize, limit.DefaultGradientSmoothing, limit.DefaultTolerance, limit.DefaultLongWindow)
	vegasLimit := limit.NewVegasLimit(5000, 20000, limit.DefaultAlphaFunc, limit.DefaultBetaFunc, limit.DefaultIncreaseFunc, limit.DefaultDecreaseFunc, limit.DefaultThresholdFunc, limit.DefaultVegasSmoothing, limit.DefaultProbeMultiplier)

	windowLimiter := limit.NewWindowedLimit(vegasLimit, 5*time.Millisecond.Nanoseconds(), 20*time.Millisecond.Nanoseconds(), 5, 10*time.Microsecond.Nanoseconds(), window.NewPercentileSampleWindowFactory(0.9, 5))
	partitions := make([]*limiter.Partition, 3)
	partitions[0] = limiter.NewPartition("normal")
	partitions[0].Percent = 0.98
	partitions[1] = limiter.NewPartition("debug")
	partitions[1].Percent = 0.01
	partitions[2] = limiter.NewPartition("sample")
	partitions[2].Percent = 0.01

	partitionedLimiter := limiter.NewPartitionedLimiter(windowLimiter, partitions)

	var wg sync.WaitGroup
	for i := 0; i < 60; i++ {
		for j := 0; j < 10000; j++ {
			wg.Add(1)
			go func() {
				ctx := context.WithValue(context.Background(), "name", "normal")
				if listener := partitionedLimiter.Acquire(ctx); listener != nil {
					// biz
					time.Sleep(1 * time.Millisecond)
					// success
					listener.OnSuccess(time.Now().UnixNano())
				}
				wg.Done()
			}()
		}
		time.Sleep(time.Second)
	}
	wg.Wait()
}

func TestSemaphore(t *testing.T) {
	s := semaphore.New(10)
	if !s.TryAcquire(1) {
		t.Fatal()
	}
	if !s.TryAcquire(1) {
		t.Fatal()
	}
	if !s.TryAcquire(1) {
		t.Fatal()
	}
	if s.GetLimit() != 10 {
		t.Fatal()
	}
	if s.GetCount() != 3 {
		t.Fatal()
	}

	s.SetLimit(2)

	if s.GetLimit() != 2 {
		t.Fatal()
	}
	if s.GetCount() != 3 {
		t.Fatal()
	}
	if s.TryAcquire(1) {
		t.Fatal()
	}
	if s.TryAcquire(1) {
		t.Fatal()
	}

	s.Release(1)
	s.Release(1)
	s.Release(1)

	if s.GetLimit() != 2 {
		t.Fatal()
	}
	if s.GetCount() != 0 {
		t.Fatal()
	}
	if !s.TryAcquire(1) {
		t.Fatal()
	}
	if s.GetLimit() != 2 {
		t.Fatal()
	}
	if s.GetCount() != 1 {
		t.Fatal()
	}
}
