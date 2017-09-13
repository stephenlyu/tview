package transform


type Transformer interface {
	To(float64) float64
	From(float64) float64
}

type ScaleTransformer interface {
	Transformer
	SetScale(scale float64)
}
