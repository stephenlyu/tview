package model

import (
	"github.com/stephenlyu/tview/transform"
	"github.com/stephenlyu/goformula/stockfunc/formula"
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

func (this *KLineModel) DrawActionCount() int {
	return 0
}

func (this *BaseModel) DrawAction(index int) formula.DrawAction {
	panic("Unimplemented")
	return nil
}

func (this *BaseModel) SetValueTransformer(transformer transform.Transformer) {
	this.valueTransformer = transformer
}

func (this *BaseModel) SetScaleTransformer(transformer transform.ScaleTransformer) {
	this.scaleTransformer = transformer
}

func (this *BaseModel) TransformRaw(v float64) float64 {
	if this.valueTransformer != nil {
		v = this.valueTransformer.To(v)
	}
	return v
}

func (this *BaseModel) Transform(v float64) float64 {
	v = this.TransformRaw(v)
	if this.scaleTransformer != nil {
		v = this.scaleTransformer.To(v)
	}
	return v
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
