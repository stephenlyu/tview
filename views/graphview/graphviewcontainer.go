package graphview

import (
	"github.com/therecipe/qt/widgets"
	"github.com/therecipe/qt/core"
	"github.com/stephenlyu/tds/entity"
	"github.com/stephenlyu/tview/views/mainwindow"
	"github.com/therecipe/qt/gui"
)

const MAX_SECONDARY_GRAPHS = 5
const DEFAULT_SECONDARY_GRAPHS = 2

const (
	MIN_ITEM_WIDTH = 1			// Item最小显示宽度, 1像素宽，1像素边距
	MAX_ITEM_WIDTH = 100

	ZOOM_OUT = 1.2
)

//go:generate qtmoc
type GraphViewContainer struct {
	widgets.QSplitter
	MainWindow *mainwindow.MainWindow

	// Children
	graphViews []*GraphView
	visibleIndices map[int]bool

	// UI control variables

	itemWidth float64
	lastVisibleIndex int
	visibleCount int

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

	// Init control variables

	this.itemWidth = this.graphViews[0].GetItemWidth()
	for _, view := range this.graphViews {
		view.SetController(this)
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

	this.SetViewVisibleRange()
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

// Controller routines

func (this *GraphViewContainer) SetViewVisibleRange() {
	usableMainWidth := this.graphViews[0].GetUsableWidth()
	lastVisibleIndex := len(this.data) - 1
	visibleCount := int(usableMainWidth / this.itemWidth)
	for _, view := range this.graphViews {
		view.SetVisibleRange(lastVisibleIndex, visibleCount)
	}
	this.lastVisibleIndex = lastVisibleIndex
}

func (this *GraphViewContainer) doZoom(scale float64) {
	if scale == 1 {
		return
	}

	var itemWidth float64
	if scale < 1 {
		if this.itemWidth <= MIN_ITEM_WIDTH {
			return
		}

		itemWidth = this.itemWidth * scale
		if itemWidth < MIN_ITEM_WIDTH {
			itemWidth = MIN_ITEM_WIDTH
		}
	} else {
		if this.itemWidth >= MAX_ITEM_WIDTH {
			return
		}

		itemWidth = this.itemWidth * scale
		if itemWidth > MAX_ITEM_WIDTH {
			itemWidth = MAX_ITEM_WIDTH
		}
	}

	this.itemWidth = itemWidth
	this.SetViewVisibleRange()
}

func (this *GraphViewContainer) HandleKeyEvent(event *gui.QKeyEvent) bool {
	var key = core.Qt__Key(event.Key())
	if key == core.Qt__Key_Up {
		this.doZoom(ZOOM_OUT)
	} else if key == core.Qt__Key_Down {
		this.doZoom(1 / ZOOM_OUT)
	}
	return false
}
