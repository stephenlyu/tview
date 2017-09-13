package graphview

import (
	"github.com/therecipe/qt/widgets"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/core"
	"github.com/stephenlyu/tview/views/mainwindow"
	"fmt"
	"time"
	"math/rand"
	"github.com/stephenlyu/tview/transform"
)

//go:generate qtmoc
type GraphView struct {
	widgets.QGraphicsView
	MainWindow *mainwindow.MainWindow

	// Coordinate transformers

	ValueTransformer transform.Transformer				// 标准坐标或对数坐标
	XScaleTransformer transform.ScaleTransformer		// X轴缩放
	YScaleTransformer transform.ScaleTransformer		// Y轴缩放


}

func CreateGraphView(parent widgets.QWidget_ITF) *GraphView {
	this := NewGraphView(parent)

	this.ConnectKeyPressEvent(this.KeyPressEvent)
	this.ConnectShowEvent(this.ShowEvent)
	this.ConnectResizeEvent(this.ResizeEvent)

	this.SetRenderHint(gui.QPainter__TextAntialiasing, true)
	this.SetCacheMode(widgets.QGraphicsView__CacheBackground)
	this.SetHorizontalScrollBarPolicy(core.Qt__ScrollBarAlwaysOff)
	this.SetVerticalScrollBarPolicy(core.Qt__ScrollBarAlwaysOff)
	this.SetAlignment(core.Qt__AlignLeft)
	this.Scale(1, -1)

	scene := widgets.NewQGraphicsScene(parent)
	scene.SetBackgroundBrush(gui.NewQBrush4(core.Qt__black, core.Qt__SolidPattern))
	this.SetScene(scene)

	r := this.Geometry()
	fmt.Println(r.X(), r.Y(), r.Width(), r.Height())
	scene.SetSceneRect2(float64(r.X()), float64(r.Y()), float64(r.X() + r.Width()), float64(r.Y() + r.Height()))

	timer := core.NewQTimer(this)
	timer.ConnectTimeout(this.addLines)
	timer.SetSingleShot(true)
	timer.Start(100)

	return this
}

func (this *GraphView) addLines() {
	pen := gui.NewQPen()
	pen.SetBrush(gui.NewQBrush4(core.Qt__red, core.Qt__SolidPattern))
	start := time.Now().UnixNano()
	for i := 0; i < 10000; i++ {
		this.Scene().AddLine2(float64(rand.Intn(1000)), float64(rand.Intn(1000)), float64(rand.Intn(1000)), float64(rand.Intn(1000)), pen)
	}
	fmt.Println("time cost: ", (time.Now().UnixNano()  - start) / int64(time.Millisecond))
	this.Scene().Update(this.Scene().SceneRect())
}

func (this *GraphView) SetMainWindow(window *mainwindow.MainWindow) {
	this.MainWindow = window
}

func (this *GraphView) ResizeEvent(event *gui.QResizeEvent) {
	fmt.Println("ResizeEvent")
	r := this.Geometry()
	fmt.Println(r.X(), r.Y(), r.Width(), r.Height())
	this.Scene().SetSceneRect2(float64(r.X()), float64(r.Y()), float64(r.X() + r.Width()), float64(r.Y() + r.Height()))
}

func (this *GraphView) ShowEvent(event *gui.QShowEvent) {
	fmt.Println("ShowEvent")
}

func (this *GraphView) KeyPressEvent(event *gui.QKeyEvent) {
	fmt.Println("KeyPressEvent")
}