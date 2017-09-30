package model

import (
	"github.com/stephenlyu/tview/transform"
)

type BaseModel struct {
	valueTransformer transform.Transformer
	scaleTransformer transform.ScaleTransformer

	listeners []ModelListener
}

func (this *BaseModel) Count() int {
	panic("Unimplemented")
	return 0
}

func (this *BaseModel) GetRaw(index int) []float64 {
	panic("Unimplemented")
	return nil
}

func (this *BaseModel) Get(index int) []float64 {
	panic("Unimplemented")
	return nil
}

func (this *BaseModel) GetGraphTypes() []int {
	panic("Unimplemented")
	return nil
}

func (this *BaseModel) SetValueTransformer(transformer transform.Transformer) {
	this.valueTransformer = transformer
}

func (this *BaseModel) SetScaleTransformer(transformer transform.ScaleTransformer) {
	this.scaleTransformer = transformer
}

func (this *BaseModel) AddListener(listener ModelListener) {
	for _, l := range this.listeners {
		if l == listener {
			return
		}
	}
	this.listeners = append(this.listeners, listener)
}

func (this *BaseModel) RemoveListener(listener ModelListener) {
	for i, l := range this.listeners {
		if l == listener {
			this.listeners = append(this.listeners[:i], this.listeners[i+1:]...)
			return
		}
	}
}

func (this *BaseModel) NotifyDataChanged() {
	for _, listener := range this.listeners {
		listener.OnDataChanged()
	}
}

func (this *BaseModel) NotifyLastDataChanged() {
	for _, listener := range this.listeners {
		listener.OnLastDataChanged()
	}
}
