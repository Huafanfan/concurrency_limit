package measurement

type Measurement interface {
	Add(sample float64) float64
	Get() float64
	Reset()
	Update(func(float64) float64)
}
