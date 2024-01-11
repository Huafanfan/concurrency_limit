package window

type PercentileSampleWindowFactory struct {
	Percentile float64
	WindowSize int
}

func NewPercentileSampleWindowFactory(percentile float64, windowSize int) SampleWindowFactory {
	return &PercentileSampleWindowFactory{
		Percentile: percentile,
		WindowSize: windowSize,
	}
}

func (f *PercentileSampleWindowFactory) NewInstance() SampleWindow {
	return NewImmutablePercentileSampleWindow(f.Percentile, f.WindowSize)
}
