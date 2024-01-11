package limit

import (
	"math"
	"sync"
	"sync/atomic"
	"time"

	"git.garena.com/yifan.zhangyf/concurrency_limit/limits_core"
	"git.garena.com/yifan.zhangyf/concurrency_limit/limits_core/limit/window"
)

var (
	DefaultMinWindowTime   = time.Second.Nanoseconds()
	DefaultMaxWindowTime   = time.Second.Nanoseconds()
	DefaultMinRttThershold = 100 * time.Millisecond.Nanoseconds()
	DefaultWindowSize      = 10
)

type WindowedLimit struct {
	Delegate            limits_core.Limit
	NextUpdateTime      int64
	MinWindowTime       int64
	MaxWindowTime       int64
	WindowSize          int
	MinRttThreshold     int64
	SampleWindowFactory window.SampleWindowFactory
	Sample              atomic.Value

	mu sync.Mutex
}

func NewWindowedLimit(delegate limits_core.Limit, MinWindowTime int64, MaxWindowTime int64, WindowSize int, MinRttThreshold int64, SampleWindowFactory window.SampleWindowFactory) limits_core.Limit {
	wl := &WindowedLimit{
		Delegate:            delegate,
		MinWindowTime:       MinWindowTime,
		MaxWindowTime:       MaxWindowTime,
		WindowSize:          WindowSize,
		MinRttThreshold:     MinRttThreshold,
		SampleWindowFactory: SampleWindowFactory,
	}
	wl.Sample.Store(SampleWindowFactory.NewInstance())
	return wl
}

func (w *WindowedLimit) OnSample(startTime int64, rtt int64, inflight int, didDrop bool) {
	endTime := startTime + rtt
	if rtt < w.MinRttThreshold {
		return
	}
	newSample := w.Sample.Load().(window.SampleWindow).AddSample(rtt, inflight, didDrop)
	w.Sample.Store(newSample)

	if endTime > w.NextUpdateTime {
		current := w.Sample.Load().(window.SampleWindow)
		w.Sample.Store(w.SampleWindowFactory.NewInstance())
		w.mu.Lock()
		w.NextUpdateTime = endTime + int64(math.Min(math.Max(float64(current.GetCandidateRttNanos()*2), float64(w.MinWindowTime)), float64(w.MaxWindowTime)))
		w.mu.Unlock()
		if w.isWindowReady(current) {
			w.Delegate.OnSample(startTime, current.GetTrackedRttNanos(), current.GetMaxInFlight(), current.GetDidDrop())
		}
	}
	GetUnifiedLogger().Flush()
}

func (w *WindowedLimit) isWindowReady(sample window.SampleWindow) bool {
	return sample.GetCandidateRttNanos() < math.MaxInt64 && sample.GetSampleCount() >= w.WindowSize
}

func (w *WindowedLimit) GetLimit() int {
	return w.Delegate.GetLimit()
}
