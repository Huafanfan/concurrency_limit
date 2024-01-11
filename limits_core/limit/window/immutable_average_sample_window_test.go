package window

import "testing"

var (
	bigRtt      = int64(5000)
	moderateRtt = int64(500)
	lowRtt      = int64(10)
)

func TestCalculateAverage(t *testing.T) {
	window := NewImmutableAverageSampleWindow()
	window = window.AddSample(bigRtt, 1, false)
	window = window.AddSample(moderateRtt, 1, false)
	window = window.AddSample(lowRtt, 1, false)
	if (bigRtt+moderateRtt+lowRtt)/3 != window.GetTrackedRttNanos() {
		t.Fatal("failed!", (bigRtt+moderateRtt+lowRtt)/3, "!=", window.GetTrackedRttNanos())
	}
}

func TestDroppedSampleShouldChangeTrackedAverage(t *testing.T) {
	window := NewImmutableAverageSampleWindow()
	window = window.AddSample(bigRtt, 1, false)
	window = window.AddSample(moderateRtt, 1, false)
	window = window.AddSample(lowRtt, 1, false)
	window = window.AddSample(bigRtt, 1, true)
	if (bigRtt+moderateRtt+lowRtt+bigRtt)/4 != window.GetTrackedRttNanos() {
		t.Fatal("failed!", (bigRtt+moderateRtt+lowRtt+bigRtt)/4, "!=", window.GetTrackedRttNanos())
	}
}
