package limit

import (
	"fmt"
	"git.garena.com/yifan.zhangyf/concurrency_limit/limits_core"
	"math"
	"math/rand"
	"time"
)

const (
	DefaultVegasSmoothing  = 1.0
	DefaultProbeMultiplier = 30
)

var (
	DefaultAlphaFunc     = func(limit int) int { return int(3 * math.Log10(float64(limit))) }
	DefaultBetaFunc      = func(limit int) int { return int(6 * math.Log10(float64(limit))) }
	DefaultThresholdFunc = func(limit int) int { return int(math.Log10(float64(limit))) }
	DefaultIncreaseFunc  = func(limit float64) float64 { return limit + math.Log10(limit) }
	DefaultDecreaseFunc  = func(limit float64) float64 { return limit - math.Log10(limit) }
)

type VegasLimit struct {
	*AbstractLimit
	EstimatedLimit  float64
	RttNoLoad       int64
	MaxLimit        int
	Smoothing       float64
	AlphaFunc       func(int) int
	BetaFunc        func(int) int
	ThresholdFunc   func(int) int
	IncreaseFunc    func(float64) float64
	DecreaseFunc    func(float64) float64
	ProbeMultiplier int
	ProbeCount      int
	ProbeJitter     float64
}

func NewVegasLimit(initialLimit int, maxConcurrency int, alphaFunc func(int) int, betaFunc func(int) int, increaseFunc func(float64) float64, decreaseFunc func(float64) float64, thresholdFunc func(int) int, smoothing float64, probeMultiplier int) limits_core.Limit {
	v := &VegasLimit{
		EstimatedLimit:  float64(initialLimit),
		Smoothing:       smoothing,
		MaxLimit:        maxConcurrency,
		AlphaFunc:       alphaFunc,
		BetaFunc:        betaFunc,
		ThresholdFunc:   thresholdFunc,
		IncreaseFunc:    increaseFunc,
		DecreaseFunc:    decreaseFunc,
		ProbeMultiplier: probeMultiplier,
	}
	v.AbstractLimit = NewAbstractLimit(initialLimit)
	v.AbstractLimit._Update = v._Update
	v.ResetProbeJitter()

	return v
}

func (v *VegasLimit) ResetProbeJitter() {
	rand.Seed(time.Now().UnixNano())
	v.ProbeJitter = rand.Float64()*(1-0.5) + 0.5
}

func (v *VegasLimit) ShouldProbe() bool {
	return v.ProbeJitter*float64(v.ProbeMultiplier)*v.EstimatedLimit <= float64(v.ProbeCount)
}

func (v *VegasLimit) _Update(startTime int64, rtt int64, inflight int, didDrop bool) int {
	v.ProbeCount++
	if v.ShouldProbe() {
		GetUnifiedLogger().Debug(fmt.Sprintf("Probe MinRTT %v", rtt/1000000.0))
		v.ResetProbeJitter()
		v.ProbeCount = 0
		v.RttNoLoad = rtt
		return int(v.EstimatedLimit)
	}

	if v.RttNoLoad == 0 || rtt < v.RttNoLoad {
		GetUnifiedLogger().Debug(fmt.Sprintf("New MinRTT %v", rtt/1000000.0))
		v.RttNoLoad = rtt
		return int(v.EstimatedLimit)
	}
	return v.updateEstimatedLimit(rtt, inflight, didDrop)
}

func (v *VegasLimit) updateEstimatedLimit(rtt int64, inflight int, didDrop bool) int {
	queueSize := int(math.Ceil(v.EstimatedLimit * (1 - float64(v.RttNoLoad)/float64(rtt))))
	newLimit := 0.0
	if didDrop {
		newLimit = v.DecreaseFunc(v.EstimatedLimit)
	} else if inflight*2 < int(v.EstimatedLimit) {
		return int(v.EstimatedLimit)
	} else {
		alpha := v.AlphaFunc(int(v.EstimatedLimit))
		beta := v.BetaFunc(int(v.EstimatedLimit))
		threshold := v.ThresholdFunc(int(v.EstimatedLimit))

		if queueSize <= threshold {
			newLimit = v.EstimatedLimit + float64(beta)
		} else if queueSize < alpha {
			newLimit = v.IncreaseFunc(v.EstimatedLimit)
		} else if queueSize > beta {
			newLimit = v.DecreaseFunc(v.EstimatedLimit)
		} else {
			return int(v.EstimatedLimit)
		}
	}

	newLimit = math.Max(1, math.Min(float64(v.MaxLimit), newLimit))
	newLimit = (1-v.Smoothing)*v.EstimatedLimit + v.Smoothing*newLimit
	if int(newLimit) != int(v.EstimatedLimit) {
		GetUnifiedLogger().Debug(fmt.Sprintf("New limit=%v shortRtt=%v ms longRtt=%v ms queueSize=%v",
			newLimit,
			v.RttNoLoad/1000000.0,
			rtt/1000000,
			queueSize))
	}

	v.EstimatedLimit = newLimit
	return int(v.EstimatedLimit)
}
