package model

import (
	"github.com/stephenlyu/tview/constants"
)

type KLineModel struct {
	BaseModel

	data *Data
}

func NewKLineModel(data *Data) *KLineModel {
	return &KLineModel{data: data}
}

func (this *KLineModel) Count() int {
	return this.data.Count()
}

func (this *KLineModel) GetRaw(index int) []float64 {
	if index < 0 || index >= this.data.Count() {
		panic("bad model index")
	}

	r := this.data.Get(index)

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

	return []float64{open, close, high, low}
}

func (this *KLineModel) Get(index int) []float64 {
	values := this.GetRaw(index)

	if this.scaleTransformer != nil {
		values[0] = this.scaleTransformer.To(values[0])
		values[1] = this.scaleTransformer.To(values[1])
		values[2] = this.scaleTransformer.To(values[2])
		values[3] = this.scaleTransformer.To(values[3])
	}

	return values
}

func (this *KLineModel) GetGraphTypes() []constants.GraphType {
	return []constants.GraphType{constants.GraphTypeKLine}
}
