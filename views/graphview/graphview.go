package graphview

import (
	"github.com/therecipe/qt/widgets"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/core"
	"github.com/stephenlyu/tview/views/mainwindow"
	"fmt"
)

//go:generate qtmoc
type GraphView struct {
	widgets.QGraphicsView
	MainWindow *mainwindow.MainWindow
}

func CreateGraphView(parent widgets.QWidget_ITF) *GraphView {
	this := NewGraphView(parent)

	this.ConnectKeyPressEvent(this.KeyPressEvent)
	this.ConnectShowEvent(this.ShowEvent)
	this.ConnectResizeEvent(this.ResizeEvent)

	this.SetHorizontalScrollBarPolicy(core.Qt__ScrollBarAlwaysOff)
	this.SetVerticalScrollBarPolicy(core.Qt__ScrollBarAlwaysOff)

	this.Scale(1, -1)
	scene := widgets.NewQGraphicsScene(parent)
	scene.SetBackgroundBrush(gui.NewQBrush4(core.Qt__black, core.Qt__SolidPattern))
	this.SetScene(scene)

	r := this.Geometry()
	fmt.Println(r.X(), r.Y(), r.Width(), r.Height())
	scene.SetSceneRect2(float64(r.X()), float64(r.Y()), float64(r.X() + r.Width()), float64(r.Y() + r.Height()))

	pen := gui.NewQPen()
	pen.SetBrush(gui.NewQBrush4(core.Qt__red, core.Qt__SolidPattern))
	this.Scene().AddLine2(0, 0, 100, 100, pen)

	this.Scene().AddEllipse2(0, 0, 1, 1, pen, pen.Brush())

	return this
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