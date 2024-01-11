package window

type AverageSampleWindowFactory struct {
}

func NewAverageSampleWindowFactory() SampleWindowFactory {
	return &AverageSampleWindowFactory{}
}

func (f *AverageSampleWindowFactory) NewInstance() SampleWindow {
	return NewImmutableAverageSampleWindow()
}
