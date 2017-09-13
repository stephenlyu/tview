package graphview

import (
	"github.com/therecipe/qt/widgets"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/core"
	"github.com/stephenlyu/tview/views/mainwindow"
	"fmt"
	"github.com/stephenlyu/tview/transform"
	"github.com/stephenlyu/tview/graphs/klinegraph"
	"github.com/stephenlyu/tview/graphs"
	"github.com/stephenlyu/tds/entity"
	"github.com/stephenlyu/tview/model"
)

const (
	VISIBLE_KLINES_MIN = 16		// 最少可见K线数
	BEST_ITEM_WIDTH = 20		// Item最佳显示宽度
	MIN_ITEM_WIDTH = 3			// Item最小显示宽度, 1像素宽，1像素边距
)

const (
	H_MARGIN = 0
	V_MARGIN = 0
)

//go:generate qtmoc
type GraphView struct {
	widgets.QGraphicsView
	MainWindow *mainwindow.MainWindow

	// Coordinate transformers

	ValueTransformer transform.Transformer				// 标准坐标或对数坐标
	XScaleTransformer transform.ScaleTransformer		// X轴缩放
	YScaleTransformer transform.ScaleTransformer		// Y轴缩放

	// Data

	Data []entity.Record

	// Models

	Models []model.Model

	// Graphs

	Graphs []graphs.Graph

	// State Variables

	IsMainGraph bool									// 是否是主图
	IsLogCoordinate bool								// 是否对数坐标
	FirstVisibleIndex, LastVisibleIndex int
	ItemWidth float64									// 每个数据占用的屏幕宽度

	yMax, yMin float64
}

func CreateGraphView(isMain bool, parent widgets.QWidget_ITF) *GraphView {
	this := NewGraphView(parent)
	this.IsMainGraph = isMain
	this.init()

	this.ValueTransformer = transform.NewEQTransformer()
	this.XScaleTransformer = transform.NewLogicTransformer(1)
	this.YScaleTransformer = transform.NewLogicTransformer(1)

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
	this.Graphs = nil

	// Clear data & models
	this.Data = nil
	this.Models = nil

	// Reset State
	this.ItemWidth = 0
	this.FirstVisibleIndex = 0
	this.LastVisibleIndex = 0
	this.yMax = 0
	this.yMin = 0
}

func (this *GraphView) SetData(data []entity.Record) {
	this.reset()
	this.Data = data
	if this.IsMainGraph {
		klineModel := model.NewKLineModel(data)
		klineModel.SetValueTransformer(this.ValueTransformer)
		klineModel.SetScaleTransformer(this.YScaleTransformer)

		this.Models = append(this.Models, klineModel)

		klineGraph := klinegraph.NewKLineGraph(klineModel, this.Scene(), this.XScaleTransformer)
		this.Graphs = append(this.Graphs, klineGraph)
	}

	this.LastVisibleIndex = len(data) - 1
	this.ItemWidth = BEST_ITEM_WIDTH

	// Do layout
	this.layout()
}

// UI routines

func (this *GraphView) connectEvents() {
	this.ConnectKeyPressEvent(this.KeyPressEvent)
	this.ConnectShowEvent(this.ShowEvent)
	this.ConnectResizeEvent(this.ResizeEvent)
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

	width := float64(len(this.Data)) * this.ItemWidth
	height := yMax
	fmt.Println(this.Width(), this.Height(), width, this.SceneRect().X(), this.SceneRect().Y(), this.Scene().Width(), this.Scene().Height())
	this.Scene().SetSceneRect2(0, 0, width, height)

	yCenter := (yMax + yMin) / 2
	xCenter := float64(this.LastVisibleIndex + 1) * this.ItemWidth - float64(this.Width()) / 2
	fmt.Println(this.Width(), this.Height(), width, this.SceneRect().X(), this.SceneRect().Y(), this.Scene().Width(), this.Scene().Height(), yMax, yMin, yCenter, xCenter, float64(this.LastVisibleIndex + 1) * this.ItemWidth, float64(this.Width()) / 2)

	this.CenterOn2(xCenter, yCenter)

	this.Scene().Update(core.NewQRectF())
	fmt.Println(this.Width(), this.Height(), width, this.Scene().Width(), this.Scene().Height(), yMax, yMin, yCenter, xCenter, float64(this.LastVisibleIndex + 1) * this.ItemWidth, float64(this.Width()) / 2)
}

func (this *GraphView) layout() {
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
		min, max := this.Graphs[0].GetValueRange(this.FirstVisibleIndex, this.LastVisibleIndex + 1)

		for i := 1; i < len(this.Graphs); i++ {
			min1, max1 := this.Graphs[i].GetValueRange(this.FirstVisibleIndex, this.LastVisibleIndex + 1)
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

		fmt.Println("y scale:", diff / height)
		this.YScaleTransformer.SetScale(diff / height)

		this.yMax = max
		this.yMin = min
	}

	this.UpdateUI()
}

// StackedWidget method
func (this *GraphView) SetMainWindow(window *mainwindow.MainWindow) {
	this.MainWindow = window
}

// Event Handlers
func (this *GraphView) ResizeEvent(event *gui.QResizeEvent) {
	this.layout()
}

func (this *GraphView) ShowEvent(event *gui.QShowEvent) {
	fmt.Println("ShowEvent")
}

func (this *GraphView) KeyPressEvent(event *gui.QKeyEvent) {
	fmt.Println("KeyPressEvent")
}
