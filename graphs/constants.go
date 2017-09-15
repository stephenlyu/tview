package graphs

import "github.com/therecipe/qt/gui"

var COLORS = []*gui.QColor{
	gui.NewQColor3(255, 255, 255, 255),
	gui.NewQColor3(255, 255,   0, 255),
	gui.NewQColor3(255,   0, 255, 255),
	gui.NewQColor3(  0, 255,   0, 255),
	gui.NewQColor3(  0,   0, 255, 255),
}

var (
	PositiveColor = gui.NewQColor3(0xD9, 0x11, 0x1B, 0xFF)
	NegativeColor = gui.NewQColor3(0x42, 0xFF, 0xFF, 0xFF)
)
