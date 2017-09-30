package graphview

import (
	"github.com/therecipe/qt/widgets"
	"github.com/therecipe/qt/core"
	"github.com/stephenlyu/tds/entity"
	"github.com/stephenlyu/tview/views/mainwindow"
	"github.com/therecipe/qt/gui"
	"github.com/stephenlyu/tview/constants"
	"fmt"
	"github.com/stephenlyu/tds/period"
)

const MAX_SECONDARY_GRAPHS = 5
const DEFAULT_SECONDARY_GRAPHS = 2


//go:generate qtmoc
type GraphViewContainer struct {
	widgets.QWidget
	MainWindow *mainwindow.MainWindow

	// Children
	splitter *widgets.QSplitter
	xBar *XBar
	graphViews []*GraphView
	visibleIndices map[int]bool

	// UI control variables

	itemWidth float64
	lastVisibleIndex int
	firstVisibleIndex int
	visibleCount int

	// Data & Model

	data []entity.Record
	currentIndex int
}

// Life Cycle Routines

func CreateGraphViewContainer(parent widgets.QWidget_ITF) *GraphViewContainer {
	ret := NewGraphViewContainer(parent, core.Qt__Widget)
	ret.currentIndex = -1
	ret.lastVisibleIndex = -1

	ret.init()
	return ret
}

func (this *GraphViewContainer) createGraphView(isMain bool) *GraphView {
	decorator := CreateGraphViewDecorator(isMain, this)
	graphView := decorator.GraphView()
	this.splitter.AddWidget(decorator)
	return graphView
}

func (this *GraphViewContainer) init() {
	// 创建Layout
	layout := widgets.NewQVBoxLayout()
	layout.SetContentsMargins(0, 0, 0, 0)
	layout.SetSpacing(0)
	this.SetLayout(layout)

	// Create QSplitter
	this.splitter = widgets.NewQSplitter(this)
	this.splitter.SetOrientation(core.Qt__Vertical)
	this.splitter.SetOpaqueResize(false)
	this.splitter.SetHandleWidth(1)
	this.splitter.SetStyleSheet("QSplitter::handle { background-color: gray; }")
	layout.AddWidget(this.splitter, 0, 0)

	// Create XBar
	this.xBar = CreateXBar(this)
	layout.AddWidget(this.xBar, 0, 0)

	// Create Graph views
	this.visibleIndices = make(map[int]bool)
	this.graphViews = make([]*GraphView, MAX_SECONDARY_GRAPHS + 1)

	// Create main graph view
	graphView := this.createGraphView(true)
	graphView.SetName("MainGraphView")
	this.splitter.SetStretchFactor(0, 3)

	this.graphViews[0] = graphView

	// Create secondary graph view
	for i := 0; i < MAX_SECONDARY_GRAPHS; i++ {
		graphView := this.createGraphView(false)
		graphView.SetName(fmt.Sprintf("SecondaryGraphView%d", i + 1))
		this.splitter.SetStretchFactor(i + 1, 1)
		this.graphViews[i + 1] = graphView
	}

	for i := 0; i < DEFAULT_SECONDARY_GRAPHS + 1; i++ {
		this.ShowGraph(i)
	}

	for i := DEFAULT_SECONDARY_GRAPHS; i < MAX_SECONDARY_GRAPHS; i++ {
		this.HideSecondaryGraph(i)
	}

	// Init control variables

	this.itemWidth = this.graphViews[0].GetItemWidth()
	for _, view := range this.graphViews {
		view.SetController(this)
	}

	this.ConnectResizeEvent(this.ResizeEvent)
}

// StackedWidget method
func (this *GraphViewContainer) SetMainWindow(window *mainwindow.MainWindow) {
	this.MainWindow = window
}

// UI Control routines

func (this *GraphViewContainer) ShowGraph(index int) {
	this.graphViews[index].Decorator.Show()
	this.visibleIndices[index] = true
}

func (this *GraphViewContainer) HideGraph(index int) {
	this.graphViews[index].Decorator.Hide()
	this.visibleIndices[index] = false
}

func (this *GraphViewContainer) ShowSecondaryGraph(index int) {
	this.ShowGraph(index + 1)
}

func (this *GraphViewContainer) HideSecondaryGraph(index int) {
	this.HideGraph(index + 1)
}

// Data & model routines

func (this *GraphViewContainer) SetData(data []entity.Record, p period.Period) {
	this.data = data

	for _, view := range this.graphViews {
		view.SetData(data, p)
	}
	this.xBar.SetData(data, p)

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
	lastVisibleIndex := this.lastVisibleIndex
	if lastVisibleIndex < 0 {
		lastVisibleIndex = len(this.data) - 1
	}
	visibleCount := int(usableMainWidth / this.itemWidth)
	if lastVisibleIndex < visibleCount {
		lastVisibleIndex = visibleCount
	}
	if lastVisibleIndex >= len(this.data) {
		lastVisibleIndex = len(this.data) - 1
	}

	// 确保能够容纳整数根K线
	this.itemWidth = usableMainWidth / float64(visibleCount)

	for _, view := range this.graphViews {
		view.SetVisibleRange(lastVisibleIndex, visibleCount)
	}
	this.xBar.SetVisibleRange(lastVisibleIndex, visibleCount)
	this.lastVisibleIndex = lastVisibleIndex
	this.firstVisibleIndex = lastVisibleIndex - visibleCount
	if this.firstVisibleIndex < 0 {
		this.firstVisibleIndex = 0
	}
}

func (this *GraphViewContainer) doZoom(scale float64) {
	if scale == 1 {
		return
	}

	var itemWidth float64
	if scale < 1 {
		if this.itemWidth <= constants.MIN_ITEM_WIDTH {
			return
		}

		itemWidth = this.itemWidth * scale
		if itemWidth < constants.MIN_ITEM_WIDTH {
			itemWidth = constants.MIN_ITEM_WIDTH
		}
	} else {
		if this.itemWidth >= constants.MAX_ITEM_WIDTH {
			return
		}

		itemWidth = this.itemWidth * scale
		if itemWidth > constants.MAX_ITEM_WIDTH {
			itemWidth = constants.MAX_ITEM_WIDTH
		}
	}

	this.itemWidth = itemWidth
	this.SetViewVisibleRange()
}

func (this *GraphViewContainer) ResizeEvent(event *gui.QResizeEvent) {
	this.SetViewVisibleRange()
}

func (this *GraphViewContainer) HandleKeyEvent(event *gui.QKeyEvent) bool {
	var key = core.Qt__Key(event.Key())
	if key == core.Qt__Key_Up {
		this.doZoom(constants.ZOOM_OUT)
	} else if key == core.Qt__Key_Down {
		this.doZoom(1 / constants.ZOOM_OUT)
	} else if key == core.Qt__Key_Left {
		if this.currentIndex > 0 {
			x, y := this.graphViews[0].GetItemXY(this.currentIndex - 1)
			this.TrackPoint(this.currentIndex - 1, x, y)
			if this.currentIndex - 1 < this.firstVisibleIndex {
				this.firstVisibleIndex = this.currentIndex - 1
			}
		}
	} else if key == core.Qt__Key_Right {
		if this.currentIndex < len(this.data) - 1 {
			x, y := this.graphViews[0].GetItemXY(this.currentIndex + 1)
			this.TrackPoint(this.currentIndex + 1, x, y)
			if this.currentIndex + 1 > this.lastVisibleIndex {
				this.lastVisibleIndex = this.currentIndex + 1
			}
		}
	}
	return false
}

// Track Point
// x: global x coordinate
// y: global y coordinate
func (this *GraphViewContainer) TrackPoint(currentIndex int, x float64, y float64) {
	for i, view := range this.graphViews {
		if !this.visibleIndices[i] {
			continue
		}
		this.currentIndex = currentIndex
		view.TrackPoint(currentIndex, x, y)
	}
	this.xBar.TrackPoint(currentIndex, x, y)
}

func (this *GraphViewContainer) TrackXY(globalX float64, globalY float64) {
	this.xBar.TrackXY(globalX, globalY)
}

func (this *GraphViewContainer) CompleteTrackPoint() {
	for _, view := range this.graphViews {
		view.CompleteTrackPoint()
	}
	this.xBar.CompleteTrackPoint()
	this.currentIndex = -1
}

func (this *GraphViewContainer) SetVisibleRangeIndex(firstVisibleIndex int, lastVisibleIndex int) {
	visibleCount := lastVisibleIndex - firstVisibleIndex + 1
	if visibleCount < constants.VISIBLE_KLINES_MIN {
		return
	}

	usableMainWidth := this.graphViews[0].GetUsableWidth()
	this.itemWidth = usableMainWidth / float64(visibleCount)

	for _, view := range this.graphViews {
		view.SetVisibleRange(lastVisibleIndex, visibleCount)
	}
	this.lastVisibleIndex = lastVisibleIndex
}
