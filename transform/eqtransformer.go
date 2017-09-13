package transform


type EQTransformer struct {
}

func NewEQTransformer() *EQTransformer {
	return &EQTransformer{}
}

func (this *EQTransformer) To(value float64) float64 {
	return value
}

func (this *EQTransformer) From(value float64) float64 {
	return value
}
