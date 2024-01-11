package window

import "testing"

func TestCalculateP50(t *testing.T) {
	window := NewImmutablePercentileSampleWindow(0.5, 10)
	window = window.AddSample(bigRtt, 0, false)
	window = window.AddSample(moderateRtt, 0, false)
	window = window.AddSample(lowRtt, 0, false)
	if moderateRtt != window.GetTrackedRttNanos() {
		t.Fatal()
	}
}

func TestDroppedSampleShouldChangeTrackedRtt(t *testing.T) {
	window := NewImmutablePercentileSampleWindow(0.5, 10)
	window = window.AddSample(lowRtt, 1, false)
	window = window.AddSample(bigRtt, 1, true)
	window = window.AddSample(bigRtt, 1, true)
	if bigRtt != window.GetTrackedRttNanos() {
		t.Fatal()
	}
}

func TestP999ReturnsSlowestObservedRtt(t *testing.T) {
	window := NewImmutablePercentileSampleWindow(0.999, 10)
	window = window.AddSample(bigRtt, 1, false)
	window = window.AddSample(moderateRtt, 1, false)
	window = window.AddSample(lowRtt, 1, false)
	if bigRtt != window.GetTrackedRttNanos() {
		t.Fatal()
	}
}

func TestRttObservationOrderDoesntAffectResultValue(t *testing.T) {
	window := NewImmutablePercentileSampleWindow(0.999, 10)
	window = window.AddSample(moderateRtt, 1, false)
	window = window.AddSample(lowRtt, 1, false)
	window = window.AddSample(bigRtt, 1, false)
	window = window.AddSample(lowRtt, 1, false)
	if bigRtt != window.GetTrackedRttNanos() {
		t.Fatal()
	}
}
