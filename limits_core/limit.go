package limits_core

type Limit interface {
	GetLimit() int
	OnSample(startTime int64, rtt int64, inflight int, didDrop bool)
}
