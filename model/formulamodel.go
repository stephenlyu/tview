package model

import (
	"github.com/stephenlyu/goformula/stockfunc/formula"
)

type FormulaModel struct {
	BaseModel

	varNames []string

	formula.Formula
}

func NewFormulaModel(f formula.Formula) *FormulaModel {
	varNames := make([]string, f.VarCount())
	for i := 0; i < f.VarCount(); i++ {
		varNames[i] = f.VarName(i)
	}

	return &FormulaModel{
		varNames: varNames,
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

func (this *FormulaModel) GetNames() []string {
	return this.varNames
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

func (this *FormulaModel) VarCount() int {
	return this.Formula.VarCount()
}

func (this *FormulaModel) DrawActionCount() int {
	return len(this.Formula.DrawActions())
}

func (this *FormulaModel) DrawAction(index int) formula.DrawAction {
	if index < 0 || index >= this.DrawActionCount() {
		panic("bad draw action index")
	}
	return this.Formula.DrawActions()[index]
}

func (this *FormulaModel) NoDraw(index int) bool {
	return this.Formula.NoDraw(index)
}

func (this *FormulaModel) NoText(index int) bool {
	return this.Formula.NoText(index)
}

func (this *FormulaModel) DrawAbove(index int) bool {
	return this.Formula.DrawAbove(index)
}

func (this *FormulaModel) NoFrame(index int) bool {
	return this.Formula.NoFrame(index)
}

func (this *FormulaModel) Color(index int) *formula.Color {
	return this.Formula.Color(index)
}

func (this *FormulaModel) LineThick(index int) int {
	return this.Formula.LineThick(index)
}

func (this *FormulaModel) LineStyle(index int) int {
	return this.Formula.LineStyle(index)
}

func (this *FormulaModel) GraphType(index int) int {
	return this.Formula.GraphType(index)
}
