package graphview

import (
	"github.com/therecipe/qt/widgets"
	"github.com/stephenlyu/tview/uigen"
	"github.com/therecipe/qt/core"
)

//go:generate qtmoc
type GraphViewDecorator struct {
	widgets.QWidget
	uigen.UIGraphViewDecorator

	graphView *GraphView
}

func CreateGraphViewDecorator(isMain bool, parent widgets.QWidget_ITF) *GraphViewDecorator {
	ret := NewGraphViewDecorator(parent, core.Qt__Widget)
	ret.SetupUI(&ret.QWidget)
	ret.init(isMain)
	return ret
}

func (this *GraphViewDecorator) init(isMain bool) {
	this.graphView = CreateGraphView(isMain, this, this)
	this.ContentLayout.AddWidget(this.graphView, 0, 0)
	this.SetStyleSheet("GraphViewDecorator{background-color: black;}")
}

func (this *GraphViewDecorator) GraphView() *GraphView {
	return this.graphView
}
