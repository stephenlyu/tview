package model

import (
	"github.com/stephenlyu/goformula/stockfunc/formula"
	"github.com/stephenlyu/goformula/stockfunc/function"
	"github.com/stephenlyu/goformula/stockfunc"
)

var formulaFactory = stockfunc.NewFormulaFactory(true)

type FormulaCreator interface {
	CreateFormula(data *Data) (error, formula.Formula)
}

type FormulaCreatorFactory interface {
	CreateFormulaCreator(args []float64) FormulaCreator
}

type luaFormulaCreatorFactory struct {
	luaFile string
}

func NewLuaFormulaCreatorFactory(luaFile string) FormulaCreatorFactory {
	return &luaFormulaCreatorFactory{luaFile:luaFile}
}

func (this *luaFormulaCreatorFactory) CreateFormulaCreator(args []float64) FormulaCreator {
	return &luaFormulaCreator{
		factory: this,
		args: args,
	}
}

type luaFormulaCreator struct {
	factory *luaFormulaCreatorFactory
	args []float64
}

func (this *luaFormulaCreator) CreateFormula(data *Data) (error, formula.Formula) {
	return formulaFactory.NewLuaFormula(this.factory.luaFile, function.RecordVector(data.Records()), this.args)
}

type easyLangFormulaCreatorFactory struct {
	easyLangFile string
}

func NewEasyLangFormulaCreatorFactory(easyLangFile string) FormulaCreatorFactory {
	return &easyLangFormulaCreatorFactory{easyLangFile:easyLangFile}
}

func (this *easyLangFormulaCreatorFactory) CreateFormulaCreator(args []float64) FormulaCreator {
	return &easyLangFormulaCreator{
		factory: this,
		args: args,
	}
}

type easyLangFormulaCreator struct {
	factory *easyLangFormulaCreatorFactory
	args []float64
}

func (this *easyLangFormulaCreator) CreateFormula(data *Data) (error, formula.Formula) {
	return formulaFactory.NewEasyLangFormula(this.factory.easyLangFile, function.RecordVector(data.Records()), this.args)
}

type FormulaLibrary struct {
	formulas map[string]FormulaCreatorFactory
}

func newFormulaLibrary() *FormulaLibrary {
	return &FormulaLibrary {
		formulas: make(map[string]FormulaCreatorFactory),
	}
}

func (this *FormulaLibrary) Register(name string, creatorFactory FormulaCreatorFactory) {
	this.formulas[name] = creatorFactory
}

func (this *FormulaLibrary) Unregister(name string, creatorFactory FormulaCreatorFactory) {
	delete(this.formulas, name)
}

func (this *FormulaLibrary) GetCreatorFactory(name string) FormulaCreatorFactory {
	return this.formulas[name]
}

func (this *FormulaLibrary) CanSupport(name string) bool {
	_, ok := this.formulas[name]
	return ok
}

var GlobalLibrary = newFormulaLibrary()
