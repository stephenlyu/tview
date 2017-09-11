package views

import "github.com/therecipe/qt/widgets"

//go:generate qtmoc
type GraphView struct {
	widgets.QGraphicsView
}

