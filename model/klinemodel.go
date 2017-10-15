package model

import (
	"github.com/stephenlyu/tview/constants"
	"github.com/stephenlyu/goformula/stockfunc/formula"
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

func (this *KLineModel) GetNames() []string {
	return []string{"OPEN", "CLOSE", "HIGH", "LOW"}
}

func (this *KLineModel) VarCount() int {
	return 1
}

func (this *KLineModel) NoDraw(index int) bool {
	return false
}

func (this *KLineModel) NoText(index int) bool {
	return false
}

func (this *KLineModel) DrawAbove(index int) bool {
	return false
}

func (this *KLineModel) NoFrame(index int) bool {
	return false
}

func (this *KLineModel) Color(index int) *formula.Color {
	return nil
}

func (this *KLineModel) LineThick(index int) int {
	return 1
}

func (this *KLineModel) LineStyle(index int) int {
	return 0
}

func (this *KLineModel) GraphType(index int) int {
	return constants.GraphTypeKLine
}
