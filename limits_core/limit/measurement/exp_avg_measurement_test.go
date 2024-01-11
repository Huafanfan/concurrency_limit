package measurement

import "testing"

func TestWarmup(t *testing.T) {
	avg := NewExpAvgMeasurement(100, 10)
	expected := []float64{10.0, 10.5, 11, 11.5, 12, 12.5, 13, 13.5, 14, 14.5}
	for i := 0; i < 10; i++ {
		avg.Add(float64(i + 10))
		if expected[i]-avg.Get() > 0.01 {
			t.Fatal()
		}
	}
	avg.Add(100)
	if 16.2-avg.Get() > 0.01 {
		t.Fatal()
	}
}
