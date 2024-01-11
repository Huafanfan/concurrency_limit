package window

import (
	"math"
)

type ImmutableAverageSampleWindow struct {
	MinRtt      int64
	Sum         int64
	MaxInFlight int
	SampleCount int
	DidDrop     bool
}

func NewImmutableAverageSampleWindow() SampleWindow {
	return &ImmutableAverageSampleWindow{
		MinRtt:      math.MaxInt64,
		MaxInFlight: 0,
		SampleCount: 0,
		DidDrop:     false,
	}
}

func (w *ImmutableAverageSampleWindow) AddSample(rtt int64, inflight int, dropped bool) SampleWindow {
	return &ImmutableAverageSampleWindow{
		MinRtt:      int64(math.Min(float64(rtt), float64(w.MinRtt))),
		Sum:         w.Sum + rtt,
		MaxInFlight: int(math.Max(float64(inflight), float64(w.MaxInFlight))),
		SampleCount: w.SampleCount + 1,
		DidDrop:     w.DidDrop || dropped,
	}
}

func (w *ImmutableAverageSampleWindow) GetCandidateRttNanos() int64 {
	return w.MinRtt
}

func (w *ImmutableAverageSampleWindow) GetSampleCount() int {
	return w.SampleCount
}

func (w *ImmutableAverageSampleWindow) GetMaxInFlight() int {
	return w.MaxInFlight
}

func (w *ImmutableAverageSampleWindow) GetDidDrop() bool {
	return w.DidDrop
}

func (w *ImmutableAverageSampleWindow) GetTrackedRttNanos() int64 {
	if w.SampleCount == 0 {
		return 0
	} else {
		return w.Sum / int64(w.SampleCount)
	}
}
