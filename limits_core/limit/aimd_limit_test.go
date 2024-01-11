package limit

import (
	"testing"
	"time"
)

func TestDefault(t *testing.T) {
	limiter := NewAIMDLimit(10, DefaultBackoffRatio, DefaultTimeout, 200, 20)
	if limiter.GetLimit() != 10 {
		t.Fatal()
	}
}

func TestIncreaseOnSuccess(t *testing.T) {
	limiter := NewAIMDLimit(20, DefaultBackoffRatio, DefaultTimeout, 200, 20)
	limiter.OnSample(0, time.Millisecond.Nanoseconds(), 10, false)
	if limiter.GetLimit() != 21 {
		t.Fatal()
	}
}

func TestDecreaseOnDrops(t *testing.T) {
	limiter := NewAIMDLimit(30, DefaultBackoffRatio, DefaultTimeout, 200, 20)
	limiter.OnSample(0, 0, 0, true)
	if limiter.GetLimit() != 27 {
		t.Fatal()
	}
}

func TestSuccessOverflow(t *testing.T) {
	limiter := NewAIMDLimit(21, DefaultBackoffRatio, DefaultTimeout, 21, 0)
	limiter.OnSample(0, time.Millisecond.Nanoseconds(), 0, false)
	if limiter.GetLimit() != 21 {
		t.Fatal()
	}
}
