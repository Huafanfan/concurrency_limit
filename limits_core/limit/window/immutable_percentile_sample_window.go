package window

import (
	"math"
	"sort"

	"github.com/Huafanfan/concurrency_limit/limits_core/utils"
)

type ImmutablePercentileSampleWindow struct {
	MinRtt       int64
	MaxInFlight  int
	SampleCount  int
	DidDrop      bool
	Percentile   float64
	ObservedRtts *utils.AtomicLongArray
}

func NewImmutablePercentileSampleWindow(Percentile float64, WindowSize int) SampleWindow {
	return &ImmutablePercentileSampleWindow{
		MinRtt:       math.MaxInt64,
		MaxInFlight:  0,
		SampleCount:  0,
		DidDrop:      false,
		Percentile:   Percentile,
		ObservedRtts: utils.NewAtomicLongArray(WindowSize),
	}
}

func (w *ImmutablePercentileSampleWindow) AddSample(rtt int64, inflight int, dropped bool) SampleWindow {
	if w.SampleCount >= w.ObservedRtts.Length() {
		return w
	}

	w.ObservedRtts.Set(w.SampleCount, rtt)

	return &ImmutablePercentileSampleWindow{
		MinRtt:       int64(math.Min(float64(rtt), float64(w.MinRtt))),
		MaxInFlight:  int(math.Max(float64(inflight), float64(w.MaxInFlight))),
		DidDrop:      w.DidDrop || dropped,
		ObservedRtts: w.ObservedRtts,
		SampleCount:  w.SampleCount + 1,
		Percentile:   w.Percentile,
	}
}

func (w *ImmutablePercentileSampleWindow) GetCandidateRttNanos() int64 {
	return w.MinRtt
}

func (w *ImmutablePercentileSampleWindow) GetSampleCount() int {
	return w.SampleCount
}

func (w *ImmutablePercentileSampleWindow) GetMaxInFlight() int {
	return w.MaxInFlight
}

func (w *ImmutablePercentileSampleWindow) GetDidDrop() bool {
	return w.DidDrop
}

func (w *ImmutablePercentileSampleWindow) GetTrackedRttNanos() int64 {
	if w.SampleCount == 0 {
		return 0
	}
	copyOfObservedRtts := make([]int64, w.SampleCount)
	for i := 0; i < w.SampleCount; i++ {
		copyOfObservedRtts[i] = w.ObservedRtts.Get(i)
	}
	sort.Slice(copyOfObservedRtts, func(i, j int) bool {
		return copyOfObservedRtts[i] < copyOfObservedRtts[j]
	})

	rttIndex := math.Round(float64(w.SampleCount) * w.Percentile)
	zeroBasedRttIndex := rttIndex - 1
	return copyOfObservedRtts[int(zeroBasedRttIndex)]
}
