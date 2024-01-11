package measurement

type ExpAvgMeasurement struct {
	Value        float64
	Sum          float64
	Window       int
	WarmupWindow int
	Count        int
}

func NewExpAvgMeasurement(window int, warmupWindow int) Measurement {
	return &ExpAvgMeasurement{
		Sum:          0.0,
		Window:       window,
		WarmupWindow: warmupWindow,
	}
}

func (m *ExpAvgMeasurement) Add(sample float64) float64 {
	if m.Count < m.WarmupWindow {
		m.Count++
		m.Sum += sample
		m.Value = m.Sum / float64(m.Count)
	} else {
		f := factor(m.Window)
		m.Value = m.Value*(1-f) + sample*f
	}
	return m.Value
}

func (m *ExpAvgMeasurement) Get() float64 {
	return m.Value
}

func (m *ExpAvgMeasurement) Reset() {
	m.Value = 0.0
	m.Count = 0
	m.Sum = 0.0
}

func (m *ExpAvgMeasurement) Update(f func(float64) float64) {
	m.Value = f(m.Value)
}

func factor(n int) float64 {
	return 2.0 / float64(n+1)
}
