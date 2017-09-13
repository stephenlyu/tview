package model

import (
	"github.com/stephenlyu/tds/entity"
	"github.com/stephenlyu/tview/transform"
	"github.com/stephenlyu/tview/constants"
)

type KLineModel struct {
	data []entity.Record

	valueTransformer transform.Transformer
	scaleTransformer transform.ScaleTransformer

	listeners []ModelListener
}

func NewKLineModel(data []entity.Record) *KLineModel {
	return &KLineModel{data: data}
}

func (this *KLineModel) Count() int {
	return len(this.data)
}

func (this *KLineModel) Get(index int) []float64 {
	if index < 0 || index >= len(this.data) {
		panic("bad model index")
	}

	r := &this.data[index]

	open := float64(r.GetOpen())
	close := float64(r.GetClose())
	high := float64(r.GetHigh())
	low := float64(r.GetLow())

	if this.valueTransformer != nil {
		open = this.valueTransformer.To(open)
		close = this.valueTransformer.To(close)
		high = this.valueTransformer.To(high)
		low = this.valueTransformer.To(low)
	}

	if this.scaleTransformer != nil {
		open = this.scaleTransformer.To(open)
		close = this.scaleTransformer.To(close)
		high = this.scaleTransformer.To(high)
		low = this.scaleTransformer.To(low)
	}

	return []float64{open, close, high, low}
}

func (this *KLineModel) GetGraphTypes() []constants.GraphType {
	return []constants.GraphType{constants.GraphTypeKLine}
}

func (this *KLineModel) SetValueTransformer(transformer transform.Transformer) {
	this.valueTransformer = transformer
}

func (this *KLineModel) SetScaleTransformer(transformer transform.ScaleTransformer) {
	this.scaleTransformer = transformer
}

func (this *KLineModel) AddListener(listener ModelListener) {
	for _, l := range this.listeners {
		if l == listener {
			return
		}
	}
	this.listeners = append(this.listeners, listener)
}

func (this *KLineModel) RemoveListener(listener ModelListener) {
	for i, l := range this.listeners {
		if l == listener {
			this.listeners = append(this.listeners[:i], this.listeners[i+1:]...)
			return
		}
	}
}

func (this *KLineModel) NotifyDataChanged() {
	for _, listener := range this.listeners {
		listener.OnDataChanged()
	}
}

func (this *KLineModel) NotifyLastDataChanged() {
	for _, listener := range this.listeners {
		listener.OnLastDataChanged()
	}
}
