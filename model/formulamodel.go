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
		for i := range values {
			values[i] = this.scaleTransformer.To(values[i])
		}
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
