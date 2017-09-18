package graphview

import (
	"github.com/therecipe/qt/widgets"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/core"
	"github.com/stephenlyu/tview/views/mainwindow"
	"github.com/stephenlyu/tview/transform"
	"github.com/stephenlyu/tview/graphs/klinegraph"
	"github.com/stephenlyu/tview/graphs"
	"github.com/stephenlyu/tds/entity"
	"github.com/stephenlyu/tview/model"
	"github.com/stephenlyu/tview/constants"
	"github.com/stephenlyu/tview/graphs/formulagraph"
	"github.com/z-ray/log"
	"github.com/cznic/mathutil"
	"github.com/stephenlyu/tview/graphs/trackline"
	"github.com/stephenlyu/tds/util"
	"github.com/stephenlyu/tview/graphs/selectrect"
	"github.com/stephenlyu/tview/graphs/separatorgraph"
)

const (
	VISIBLE_KLINES_MIN = constants.VISIBLE_KLINES_MIN
	BEST_ITEM_WIDTH = constants.BEST_ITEM_WIDTH
)

const (
	H_MARGIN = constants.H_MARGIN
	V_MARGIN = constants.V_MARGIN
)

const KLINE_MODEL = "__kline__"

type Controller interface {
	SetVisibleRangeIndex(firstVisibleIndex int, lastVisibleIndex int)
	HandleKeyEvent(event *gui.QKeyEvent) bool

	TrackPoint(currentIndex int, x float64, y float64)
	CompleteTrackPoint()
}

//go:generate qtmoc
type GraphView struct {
	widgets.QGraphicsView
	MainWindow *mainwindow.MainWindow
	Decorator *GraphViewDecorator

	// Coordinate transformers

	ValueTransformer transform.Transformer				// 标准坐标或对数坐标
	XScaleTransformer transform.ScaleTransformer		// X轴缩放
	YScaleTransformer transform.ScaleTransformer		// Y轴缩放

	// Data

	Data *model.Data

	// Models

	Models map[string]model.Model
	SeparatorModel *model.SeparatorModel

	// Formula creators

	FormulaCreators map[string]model.FormulaCreator

	// Graphs

	Graphs map[string]graphs.Graph
	TrackLine *trackline.TrackLine
	SelectRect *selectrect.SelectRect
	SeparatorGraph *separatorgraph.SeparatorGraph

	// Controller

	Controller Controller

	// State Variables

	IsMainGraph bool									// 是否是主图
	IsLogCoordinate bool								// 是否对数坐标
	FirstVisibleIndex, LastVisibleIndex int
	ItemWidth float64									// 每个数据占用的屏幕宽度

	isTracking bool

	PressedPoint *core.QPointF

	yMax, yMin float64

	yValue float64										// 当前Y值
}

func CreateGraphView(isMain bool, decorator *GraphViewDecorator, parent widgets.QWidget_ITF) *GraphView {
	this := NewGraphView(parent)
	this.SetMouseTracking(true)
	this.Decorator = decorator
	this.IsMainGraph = isMain
	this.Models = make(map[string]model.Model)
	this.SeparatorModel = model.NewSeparatorModel()
	this.Graphs = make(map[string]graphs.Graph)
	this.FormulaCreators = make(map[string]model.FormulaCreator)
	this.init()

	this.ValueTransformer = transform.NewEQTransformer()
	this.XScaleTransformer = transform.NewLogicTransformer(1)
	this.YScaleTransformer = transform.NewLogicTransformer(1)

	this.TrackLine = trackline.NewTrackLine(this.Scene())
	this.SelectRect = selectrect.NewSelectRect(this.Scene())
	this.SeparatorGraph = separatorgraph.NewSeparatorGraph(this.Scene(), this.YScaleTransformer, this.SeparatorModel)

	// 设置YDecorator
	this.Decorator.YDecorator().SetModel(this.SeparatorModel)
	this.Decorator.YDecorator().SetTransformer(this.YScaleTransformer)

	return this
}

// Properties


func (this *GraphView) SetLogCoordinate(flag bool) {
	if flag == this.IsLogCoordinate {
		return
	}

	this.IsLogCoordinate = flag
	if this.IsLogCoordinate {
		this.ValueTransformer = transform.NewLogTransformer()
	} else {
		this.ValueTransformer = transform.NewEQTransformer()
	}

	this.UpdateUI()
}

// Data & Model routines

func (this *GraphView) reset() {
	// Clear graphs
	for _, graph := range this.Graphs {
		graph.Clear()
	}
	this.Graphs = make(map[string]graphs.Graph)
	this.TrackLine.Clear()

	// Clear data & models
	this.Data = nil
	this.Models = make(map[string]model.Model)

	// Reset State
	this.ItemWidth = 0
	this.FirstVisibleIndex = 0
	this.LastVisibleIndex = 0
	this.yMax = 0
	this.yMin = 0
}

func (this *GraphView) setModelTransformers(model model.Model) {
	model.SetValueTransformer(this.ValueTransformer)
	model.SetScaleTransformer(this.YScaleTransformer)
}

func (this *GraphView) SetData(data []entity.Record) {
	this.reset()
	this.Data = model.NewData(data)

	klineModel := model.NewKLineModel(this.Data)
	this.setModelTransformers(klineModel)
	this.Models[KLINE_MODEL] = klineModel

	if this.IsMainGraph {
		klineGraph := klinegraph.NewKLineGraph(klineModel, this.Scene(), this.XScaleTransformer)
		this.Graphs[KLINE_MODEL] = klineGraph
	}

	for name := range this.FormulaCreators {
		if name == KLINE_MODEL {
			continue
		}

		this.createFormulaGraph(name)
	}

	this.LastVisibleIndex = this.Data.Count() - 1
	this.ItemWidth = BEST_ITEM_WIDTH

	// Do layout
	this.Layout()
}

func (this *GraphView) RemoveFormula(name string) {
	delete(this.Models, name)
	delete(this.FormulaCreators, name)

	if graph, ok := this.Graphs[name]; ok {
		graph.Clear()
		delete(this.Graphs, name)
	}
}

func (this *GraphView) AddFormula(name string, args []float64) {
	if !model.GlobalLibrary.CanSupport(name) {
		log.Errorf("formula %s not supported", name)
		return
	}

	this.RemoveFormula(name)

	creatorFactory := model.GlobalLibrary.GetCreatorFactory(name)

	creator := creatorFactory.CreateFormulaCreator(args)

	this.FormulaCreators[name] = creator

	if this.Data != nil {
		this.createFormulaGraph(name)
	}
	this.Layout()
}

func (this *GraphView) createFormulaGraph(name string) {
	creator := this.FormulaCreators[name]

	var graphTypes []constants.GraphType
	switch name {
	case "MACD":
		graphTypes = []constants.GraphType{constants.GraphTypeLine, constants.GraphTypeLine, constants.GraphTypeStick}
	case "MA":
		graphTypes = []constants.GraphType{constants.GraphTypeLine, constants.GraphTypeLine, constants.GraphTypeLine, constants.GraphTypeLine}
	case "VOL":
		graphTypes = []constants.GraphType{constants.GraphTypeVolStick, constants.GraphTypeLine, constants.GraphTypeLine}
	}

	_, formula := creator.CreateFormula(this.Data)
	model := model.NewFormulaModel(formula, graphTypes)
	this.setModelTransformers(model)
	this.Models[name] = model
	this.Graphs[name] = formulagraph.NewFormulaGraph(model, this.Models[KLINE_MODEL], this.Scene(), this.XScaleTransformer)
}

// UI routines

func (this *GraphView) connectEvents() {
	this.ConnectKeyPressEvent(this.KeyPressEvent)
	this.ConnectResizeEvent(this.ResizeEvent)
	this.ConnectWheelEvent(this.WheelEvent)

	this.ConnectMousePressEvent(this.MousePressEvent)
	this.ConnectMouseReleaseEvent(this.MouseReleaseEvent)
	this.ConnectMouseMoveEvent(this.MouseMoveEvent)
	this.ConnectMouseDoubleClickEvent(this.MouseDoubleClickEvent)
}

func (this *GraphView) init() {
	this.SetRenderHint(gui.QPainter__TextAntialiasing, true)
	this.SetCacheMode(widgets.QGraphicsView__CacheBackground)
	this.SetHorizontalScrollBarPolicy(core.Qt__ScrollBarAlwaysOff)
	this.SetVerticalScrollBarPolicy(core.Qt__ScrollBarAlwaysOff)
	this.SetDragMode(widgets.QGraphicsView__NoDrag)
	this.SetAlignment(core.Qt__AlignLeft)
	this.Scale(1, -1)

	// 设置scene
	scene := widgets.NewQGraphicsScene(this)
	scene.SetBackgroundBrush(gui.NewQBrush4(core.Qt__black, core.Qt__SolidPattern))
	this.SetScene(scene)

	this.connectEvents()
}

func (this *GraphView) UpdateUI() {
	for _, graph := range this.Graphs {
		graph.Update(this.FirstVisibleIndex, this.LastVisibleIndex)
	}

	yMax := this.YScaleTransformer.To(this.yMax)
	yMin := this.YScaleTransformer.To(this.yMin)

	width := float64(this.Data.Count()) * this.ItemWidth
	usableWidth := this.Width() - 2 * H_MARGIN
	fullMode := width < float64(usableWidth)
	if fullMode {
		width = float64(usableWidth)
	}
	height := yMax - yMin
	this.Scene().SetSceneRect2(-H_MARGIN, yMin - V_MARGIN, width + 2 * H_MARGIN, height + 2 * V_MARGIN)

	yCenter := (yMax + yMin) / 2
	var xCenter float64
	if fullMode {
		xCenter = float64(this.Scene().Width() - 2 * H_MARGIN) / 2
	} else {
		xCenter = float64(this.LastVisibleIndex + 1) * this.ItemWidth + H_MARGIN - float64(this.Width()) / 2
	}

	this.CenterOn2(xCenter, yCenter)

	// 更新Separator lines，此时Scene的宽度已经确定

	this.SeparatorModel.Update(this.yMin, this.yMax, float64(this.Height() - 2 * V_MARGIN))
	this.SeparatorModel.NotifyDataChanged(this.yMin, this.yMax)

	this.Scene().Update(this.Scene().SceneRect())
}

func (this *GraphView) Layout() {
	if this.Data == nil || len(this.Graphs) == 0 {
		return
	}

	if this.ItemWidth <= 0 {
		this.ItemWidth = BEST_ITEM_WIDTH
	}

	// 计算ItemWidth
	width := float64(this.Width()) - 2 * H_MARGIN
	n := int(width / this.ItemWidth)
	if n < VISIBLE_KLINES_MIN {
		n = VISIBLE_KLINES_MIN
	}

	// 计算FirstVisibleIndex
	this.FirstVisibleIndex = this.LastVisibleIndex - n + 1
	if this.FirstVisibleIndex < 0 {
		this.FirstVisibleIndex = 0
	}

	// 设置X transformer Scale
	this.XScaleTransformer.SetScale(1 / this.ItemWidth)

	// 计算Y值范围

	if len(this.Graphs) > 0 {
		names := make([]string, len(this.Graphs))
		i := 0
		for name := range this.Graphs {
			names[i] = name
			i++
		}

		min, max := this.Graphs[names[0]].GetValueRange(this.FirstVisibleIndex, this.LastVisibleIndex + 1)

		for i := 1; i < len(names); i++ {
			min1, max1 := this.Graphs[names[i]].GetValueRange(this.FirstVisibleIndex, this.LastVisibleIndex + 1)
			if min1 < min {
				min = min1
			}
			if max1 > max {
				max = max1
			}
		}

		diff := max - min
		if diff <= 0 {
			diff = 1
		}
		height := float64(this.Height()) - 2 * V_MARGIN

		this.YScaleTransformer.SetScale(diff / height)

		this.yMax = max
		this.yMin = min
	}

	this.UpdateUI()
}

func (this *GraphView) updateYDecorator() {
	if this.Decorator != nil {
		this.Decorator.YDecorator().ShowValue(this.yValue)
	}
}

// StackedWidget method
func (this *GraphView) SetMainWindow(window *mainwindow.MainWindow) {
	this.MainWindow = window
}

// Event Handlers
func (this *GraphView) ResizeEvent(event *gui.QResizeEvent) {
	this.Layout()
}

func (this *GraphView) KeyPressEvent(event *gui.QKeyEvent) {
	if this.Controller != nil {
		if this.Controller.HandleKeyEvent(event) {
			return
		}
	}
}

func (this *GraphView) WheelEvent(event *gui.QWheelEvent) {
}

func (this *GraphView) MouseDoubleClickEvent(event *gui.QMouseEvent) {
	if event.Button() != core.Qt__LeftButton {
		return
	}

	if this.isTracking {
		if this.Controller != nil {
			this.Controller.CompleteTrackPoint()
		}
	} else {
		if this.Controller != nil {
			ptScene := this.MapToScene(event.Pos())
			currentIndex := int(this.XScaleTransformer.From(ptScene.X()))
			if currentIndex >= this.Data.Count() {
				currentIndex = this.Data.Count() - 1
			}

			if currentIndex < 0 {
				currentIndex = 0
			}
			this.Controller.TrackPoint(currentIndex, float64(event.GlobalX()), float64(event.GlobalY()))
		}
	}
}

func (this *GraphView) MouseMoveEvent(event *gui.QMouseEvent) {
	ptScene := this.MapToScene(event.Pos())
	this.yValue = this.YScaleTransformer.From(ptScene.Y())
	this.updateYDecorator()

	if this.isTracking {
		if this.Controller != nil {
			ptScene := this.MapToScene(event.Pos())
			currentIndex := int(this.XScaleTransformer.From(ptScene.X()))
			if currentIndex >= this.Data.Count() {
				currentIndex = this.Data.Count() - 1
			}

			if currentIndex < 0 {
				currentIndex = 0
			}
			this.Controller.TrackPoint(currentIndex, float64(event.GlobalX()), float64(event.GlobalY()))
		}
	} else if this.PressedPoint != nil {
		ptStart := this.PressedPoint
		ptScene := this.MapToScene(event.Pos())
		var x, y float64
		w := ptScene.X() - ptStart.X()
		if w < 0 {
			x = ptScene.X()
			w = -w
		} else {
			x = ptStart.X()
		}
		h := ptScene.Y() - ptStart.Y()
		if h < 0 {
			y = ptScene.Y()
			h = -h
		} else {
			y = ptStart.Y()
		}
		this.SelectRect.UpdateRect(x, y, w, h)
	}
}

func (this *GraphView) MousePressEvent(event *gui.QMouseEvent) {
	if event.Button() != core.Qt__LeftButton {
		return
	}
	this.PressedPoint = this.MapToScene(event.Pos())
}

func (this *GraphView) MouseReleaseEvent(event *gui.QMouseEvent) {
	this.PressedPoint = nil
	if this.SelectRect.GetRect() != nil {
		r := this.SelectRect.GetRect()
		this.SelectRect.Clear()

		startIndex := int(this.XScaleTransformer.From(r.X()))
		endIndex := int(this.XScaleTransformer.From(r.X() + r.Width()))

		if startIndex > 0 {
			startIndex--
		}

		if endIndex < this.Data.Count() - 1 {
			endIndex++
		}

		if this.Controller != nil {
			this.Controller.SetVisibleRangeIndex(startIndex, endIndex)
		}
	}
}

// Control routines

func (this *GraphView) GetItemWidth() float64 {
	if this.ItemWidth <= 0 {
		this.ItemWidth = BEST_ITEM_WIDTH
	}
	return this.ItemWidth
}

func (this *GraphView) GetUsableWidth() float64 {
	return float64(this.Width()) - 2 * H_MARGIN
}

// 获取主图收盘价的Global坐标
func (this *GraphView) GetItemXY(index int) (float64, float64) {
	util.Assert(this.IsMainGraph, "")
	util.Assert(index >= 0 && index < this.Data.Count(), "")

	r := this.Data.Get(index)

	close := r.GetClose()
	x := (this.XScaleTransformer.To(float64(index)) + this.XScaleTransformer.To(float64(index+1))) / 2
	y := this.YScaleTransformer.To(float64(close))

	pt := core.NewQPointF3(x, y)
	ptGlobal := this.MapToGlobal(this.MapFromScene(pt))

	return float64(ptGlobal.X()), float64(ptGlobal.Y())
}

func (this *GraphView) SetVisibleRange(lastVisibleIndex int, visibleCount int) {
	if this.Data == nil || this.Data.Count() == 0 {
		return
	}
	if lastVisibleIndex < 0 || lastVisibleIndex >= this.Data.Count() {
		return
	}

	if visibleCount <= 0 {
		return
	}

	firstVisibleIndex := int(mathutil.MaxInt32(0, int32(lastVisibleIndex - visibleCount + 1)))
	visibleCount = lastVisibleIndex - firstVisibleIndex + 1

	this.LastVisibleIndex = lastVisibleIndex
	usableWidth := float64(this.Width()) - 2 * H_MARGIN
	this.ItemWidth = usableWidth / float64(visibleCount)

	this.Layout()
}

func (this *GraphView) SetController(controller Controller) {
	this.Controller = controller
}

func (this *GraphView) TrackPoint(currentIndex int, x float64, y float64) {
	if currentIndex < this.FirstVisibleIndex {
		this.LastVisibleIndex -= (this.FirstVisibleIndex - currentIndex)
		this.Layout()
	} else if currentIndex > this.LastVisibleIndex {
		this.LastVisibleIndex = currentIndex
		this.Layout()
	}

	pt := this.MapFromGlobal(core.NewQPoint2(int(x), int(y)))
	ptScene := this.MapToScene(pt)

	this.TrackLine.UpdateTrackLine(ptScene.X(), ptScene.Y())
	this.isTracking = true

	this.yValue = this.YScaleTransformer.From(ptScene.Y())
	this.updateYDecorator()
}

func (this *GraphView) CompleteTrackPoint() {
	this.TrackLine.Clear()
	this.isTracking = false
}
