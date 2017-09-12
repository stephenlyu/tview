package mainwindow

import (
	"github.com/therecipe/qt/widgets"
	"github.com/stephenlyu/tview/uigen"
	"github.com/therecipe/qt/core"
)

type StackedWidget interface {
	widgets.QWidget_ITF
	SetMainWindow(window *MainWindow)
}

type MainWindow struct {
	uigen.UIMainwindowMainWindow
	Widget *widgets.QMainWindow

	StackWidget *widgets.QStackedWidget
}

var instance *MainWindow

func GetMainWindow(parent widgets.QWidget_ITF) *MainWindow {
	if instance != nil {
		return instance
	}

	window := &MainWindow{
		Widget: widgets.NewQMainWindow(parent, core.Qt__Window),
	}

	window.SetupUI(window.Widget)

	stackWidget := widgets.NewQStackedWidget(window.Widget)
	window.VerticalLayout.AddWidget(stackWidget, 0, 0)
	window.StackWidget = stackWidget

	instance = window
	return window
}

func (this *MainWindow) Push(widget StackedWidget) {
	if widget == nil {
		return
	}

	this.StackWidget.AddWidget(widget)
	this.StackWidget.SetCurrentWidget(widget)
	widget.SetMainWindow(this)
}

func (this *MainWindow) Pop() {
	if this.StackWidget.Count() == 0 {
		return
	}
	lastView := this.StackWidget.CurrentWidget()
	this.StackWidget.RemoveWidget(lastView)

	var i interface{} = lastView
	i.(StackedWidget).SetMainWindow(nil)
}

func (this *MainWindow) StackSize() int {
	return this.StackWidget.Count()
}
