package measurement

type MinimumMeasurement struct {
	Value float64
}

func (m *MinimumMeasurement) Add(sample float64) float64 {
	if m.Value == 0.0 || sample < m.Value {
		m.Value = sample
	}
	return m.Value
}

func (m *MinimumMeasurement) Get() float64 {
	return m.Value
}

func (m *MinimumMeasurement) Reset() {
	m.Value = 0.0
}

func (m *MinimumMeasurement) Update(f func(float64) float64) {
	m.Value = f(m.Value)
}
