package limit

import (
	"fmt"
	"log"
	"math"

	"github.com/Huafanfan/concurrency_limit/limits_core"
	"github.com/Huafanfan/concurrency_limit/limits_core/limit/measurement"
)

const (
	DefaultGradientSmoothing = 0.2
	DefaultLongWindow        = 600
	DefaultTolerance         = 1.5
)

var (
	DefaultQueueSize = func(f float64) float64 { return math.Log10(f) }
)

type Gradient2Limit struct {
	*AbstractLimit
	EstimatedLimit float64
	LastRtt        int64
	LongRtt        measurement.Measurement
	MaxLimit       int
	MinLimit       int
	QueueSize      func(float64) float64
	Smoothing      float64
	Tolerance      float64
}

func NewGradient2Limit(initialLimit int, maxConcurrency int, minLimit int, queueSize func(float64) float64, smoothing float64, rttTolerance float64, longWindow int) limits_core.Limit {
	g := &Gradient2Limit{
		EstimatedLimit: float64(initialLimit),
		LastRtt:        0,
		LongRtt:        measurement.NewExpAvgMeasurement(longWindow, 10),
		MaxLimit:       maxConcurrency,
		MinLimit:       minLimit,
		QueueSize:      queueSize,
		Smoothing:      smoothing,
		Tolerance:      rttTolerance,
	}

	g.AbstractLimit = NewAbstractLimit(initialLimit)

	g.AbstractLimit._Update = g._Update
	return g
}

func (g *Gradient2Limit) _Update(startTime int64, rtt int64, inflight int, didDrop bool) int {
	queueSize := g.QueueSize(g.EstimatedLimit)
	g.LastRtt = rtt
	shortRtt := float64(rtt)
	longRtt := g.LongRtt.Add(float64(rtt))

	// If the long RTT is substantially larger than the short RTT then reduce the long RTT measurement.
	// This can happen when latency returns to normal after a prolonged prior of excessive load.  Reducing the
	// long RTT without waiting for the exponential smoothing helps bring the system back to steady state.
	if longRtt/shortRtt > 2 {
		g.LongRtt.Update(func(v float64) float64 {
			return v * 0.95
		})
	}

	// Don't grow the limit if we are app limited
	if float64(inflight) < g.EstimatedLimit/2 {
		return int(g.EstimatedLimit)
	}

	// Rtt could be higher than rtt_noload because of smoothing rtt noload updates
	// so set to 1.0 to indicate no queuing.  Otherwise calculate the slope and don't
	// allow it to be reduced by more than half to avoid aggressive load-shedding due to
	// outliers.
	gradient := math.Max(0.5, math.Min(1.0, g.Tolerance*longRtt/shortRtt))
	newLimit := g.EstimatedLimit*gradient + queueSize
	newLimit = g.EstimatedLimit*(1-g.Smoothing) + newLimit*g.Smoothing
	newLimit = math.Max(float64(g.MinLimit), math.Min(float64(g.MaxLimit), newLimit))

	if g.EstimatedLimit != newLimit {
		log.Printf(fmt.Sprintf("New limit=%v shortRtt=%v ms longRtt=%v ms queueSize=%v gradient=%v",
			newLimit,
			g.GetLastRtt()/1000000.0,
			g.GetRttNoLoad()/1000000,
			queueSize,
			gradient))
	}

	g.EstimatedLimit = newLimit

	return int(g.EstimatedLimit)
}

func (g *Gradient2Limit) GetLastRtt() int64 {
	return g.LastRtt
}

func (g *Gradient2Limit) GetRttNoLoad() int64 {
	return int64(g.LongRtt.Get())
}
