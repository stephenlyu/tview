package views

import "github.com/therecipe/qt/widgets"

type StackedWidget interface {
	widgets.QWidget_ITF
	SetMainWindow(window TopWindow)
}

type TopWindow interface {
	Push(widget StackedWidget)
	Pop()
}
