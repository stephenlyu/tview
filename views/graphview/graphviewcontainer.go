package graphview

import (
	"github.com/therecipe/qt/widgets"
	"github.com/therecipe/qt/core"
	"github.com/stephenlyu/tds/entity"
	"github.com/stephenlyu/tview/views/mainwindow"
)

const MAX_SECONDARY_GRAPHS = 5
const DEFAULT_SECONDARY_GRAPHS = 2

//go:generate qtmoc
type GraphViewContainer struct {
	widgets.QSplitter
	MainWindow *mainwindow.MainWindow

	// Children
	graphViews []*GraphView
	visibleIndices map[int]bool

	// Data & Model

	data []entity.Record
}

// Life Cycle Routines

func CreateGraphViewContainer(parent widgets.QWidget_ITF) *GraphViewContainer {
	ret := NewGraphViewContainer(parent)
	ret.SetOrientation(core.Qt__Vertical)
	ret.SetStyleSheet("QGraphicsView { background-color: black; }")
	ret.SetOpaqueResize(false)
	ret.SetHandleWidth(0)
	ret.SetStyleSheet("QSplitter::handle { background-color: black }")

	ret.init()
	return ret
}

func (this *GraphViewContainer) init() {
	this.visibleIndices = make(map[int]bool)
	this.graphViews = make([]*GraphView, MAX_SECONDARY_GRAPHS + 1)

	graphView := CreateGraphView(true, this)
	this.AddWidget(graphView)
	this.SetStretchFactor(0, 3)

	this.graphViews[0] = graphView

	for i := 0; i < MAX_SECONDARY_GRAPHS; i++ {
		graphView := CreateGraphView(false, this)
		this.AddWidget(graphView)
		this.SetStretchFactor(i + 1, 1)
		this.graphViews[i + 1] = graphView
	}

	for i := DEFAULT_SECONDARY_GRAPHS; i < MAX_SECONDARY_GRAPHS; i++ {
		this.HideSecondaryGraph(i)
	}
}

// StackedWidget method
func (this *GraphViewContainer) SetMainWindow(window *mainwindow.MainWindow) {
	this.MainWindow = window
}

// UI Control routines

func (this *GraphViewContainer) ShowGraph(index int) {
	this.graphViews[index].Show()
	this.visibleIndices[index] = true
}

func (this *GraphViewContainer) HideGraph(index int) {
	this.graphViews[index].Hide()
	this.visibleIndices[index] = false
}

func (this *GraphViewContainer) ShowSecondaryGraph(index int) {
	this.ShowGraph(index + 1)
}

func (this *GraphViewContainer) HideSecondaryGraph(index int) {
	this.HideGraph(index + 1)
}

// Data & model routines

func (this *GraphViewContainer) SetData(data []entity.Record) {
	this.data = data

	for _, view := range this.graphViews {
		view.SetData(data)
	}
}

// Graph Routines

func (this *GraphViewContainer) AddGraphFormula(index int, name string, args []float64) {
	if index < 0 || index >= len(this.graphViews) {
		return
	}

	this.graphViews[index].AddFormula(name, args)
}

func (this *GraphViewContainer) RemoveGraphFormula(index int, name string, args []float64) {
	if index < 0 || index >= len(this.graphViews) {
		return
	}

	this.graphViews[index].RemoveFormula(name)
}
