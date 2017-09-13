package transform


type LogicTransformer struct {
	Scale float64
}

func NewLogicTransformer(scale float64) *LogicTransformer {
	return &LogicTransformer{Scale: scale}
}

func (this *LogicTransformer) To(value float64) float64 {
	if this.Scale == 0 {
		panic("Scale is 0")
	}
	return value / this.Scale
}

func (this *LogicTransformer) From(value float64) float64 {
	return value * this.Scale
}

func (this *LogicTransformer) SetScale(scale float64) {
	this.Scale = scale
}
