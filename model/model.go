package model

import (
	"github.com/stephenlyu/tview/constants"
	"github.com/stephenlyu/tview/transform"
)

type ModelListener interface {
	OnDataChanged()
	OnLastDataChanged()
}

type Model interface {
	Count() int									// data count
	Get(index int) []float64					// Get record at index
	GetRaw(index int) []float64					// Get un-scaled record at index
	GetGraphTypes() []constants.GraphType		// Get graph types of this model

	SetValueTransformer(transformer transform.Transformer)
	SetScaleTransformer(transformer transform.ScaleTransformer)

	AddListener(listener ModelListener)
	RemoveListener(listener ModelListener)

	NotifyDataChanged()
	NotifyLastDataChanged()
}
