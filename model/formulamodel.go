package model

import (
	"github.com/stephenlyu/goformula/stockfunc/formula"
	"github.com/stephenlyu/tview/constants"
)

type FormulaModel struct {
	BaseModel

	graphTypes []constants.GraphType
	Formula formula.Formula
}

func NewFormulaModel(f formula.Formula, graphTypes []constants.GraphType) *FormulaModel {
	return &FormulaModel{
		graphTypes: graphTypes,
		Formula: f,
	}
}

func (this *FormulaModel) Count() int {
	return this.Formula.Len()
}

func (this *FormulaModel) GetRaw(index int) []float64 {
	if index < 0 || index >= this.Count() {
		panic("bad model index")
	}

	values := this.Formula.Get(index)

	if this.valueTransformer != nil {
		for i := range values {
			values[i] = this.valueTransformer.To(values[i])
		}
	}

	return values
}

func (this *FormulaModel) Get(index int) []float64 {
	values := this.GetRaw(index)

	if this.scaleTransformer != nil {
		values[0] = this.scaleTransformer.To(values[0])
		values[1] = this.scaleTransformer.To(values[1])
		values[2] = this.scaleTransformer.To(values[2])
		values[3] = this.scaleTransformer.To(values[3])
	}

	return values
}

func (this *FormulaModel) GetGraphTypes() []constants.GraphType {
	return this.graphTypes
}

func (this *FormulaModel) NotifyDataChanged() {
	panic("!!!Not Supported!!!")
}

func (this *FormulaModel) NotifyLastDataChanged() {
	this.Formula.UpdateLastValue()

	for _, listener := range this.listeners {
		listener.OnLastDataChanged()
	}
}
