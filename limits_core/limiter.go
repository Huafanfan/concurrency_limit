package limits_core

type Limiter interface {
	Acquire() Listener
}

type Listener interface {
	OnSuccess(endTime int64)
	OnIgnore(endTime int64)
	OnDropped(endTime int64)
}
