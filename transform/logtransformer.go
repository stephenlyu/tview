package transform

import "math"

type LogTransformer struct {
}

func NewLogTransformer() *LogTransformer {
	return &LogTransformer{}
}

func (this *LogTransformer) To(value float64) float64 {
	return math.Log(value)
}

func (this *LogTransformer) From(value float64) float64 {
	return math.Exp(value)
}
