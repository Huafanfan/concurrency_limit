package window

type SampleWindow interface {
	AddSample(rtt int64, inflight int, dropped bool) SampleWindow
	GetCandidateRttNanos() int64
	GetTrackedRttNanos() int64
	GetSampleCount() int
	GetMaxInFlight() int
	GetDidDrop() bool
}
