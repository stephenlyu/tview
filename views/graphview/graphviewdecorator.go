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

	yDecorator *YDecorator
	topDecorator *TopDecorator
}

func CreateGraphViewDecorator(isMain bool, parent widgets.QWidget_ITF) *GraphViewDecorator {
	ret := NewGraphViewDecorator(parent, core.Qt__Widget)
	ret.SetupUI(&ret.QWidget)
	ret.init(isMain)
	return ret
}

func (this *GraphViewDecorator) init(isMain bool) {
	this.topDecorator = NewTopDecorator(this.InfoLayout)
	this.yDecorator = CreateYDecorator(this)
	this.YAxisLayout.AddWidget(this.yDecorator, 0, 0)
	this.yDecorator.SetStyleSheet("background-color: black;")

	this.graphView = CreateGraphView(isMain, this, this)
	this.ContentLayout.AddWidget(this.graphView, 0, 0)
}

func (this *GraphViewDecorator) GraphView() *GraphView {
	return this.graphView
}

func (this *GraphViewDecorator) YDecorator() *YDecorator {
	return this.yDecorator
}

func (this *GraphViewDecorator) TopDecorator() *TopDecorator {
	return this.topDecorator
}
