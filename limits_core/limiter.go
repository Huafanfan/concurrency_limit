package limits_core

import "context"

type Limiter interface {
	Acquire(ctx context.Context) Listener
}

type Listener interface {
	OnSuccess(endTime int64)
	OnIgnore(endTime int64)
	OnDropped(endTime int64)
}
