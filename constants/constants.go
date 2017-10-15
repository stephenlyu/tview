package constants

import (
	"github.com/therecipe/qt/gui"
	"github.com/stephenlyu/goformula/stockfunc/formula"
)

const (
	GraphTypeLine = formula.FORMULA_GRAPH_LINE							// Line graph, consume 1 value
	GraphTypeColorStick = formula.FORMULA_GRAPH_COLOR_STICK             // Line graph, consume 1 value
	GraphTypeStick = formula.FORMULA_GRAPH_STICK             			// Line graph, consume 1 value
	GraphTypeLineStick = formula.FORMULA_GRAPH_LINE_STICK             	// Line graph, consume 1 value
	GraphTypeVolStick = formula.FORMULA_GRAPH_VOL_STICK           		// Line graph, consume 1 value

	GraphTypeKLine = 99
)


const (
	VISIBLE_KLINES_MIN = 16		// 最少可见K线数
	BEST_ITEM_WIDTH = 10		// Item最佳显示宽度
)

const (
	H_MARGIN = 0
	V_MARGIN = 10
)

const (
	MIN_ITEM_WIDTH = 1			// Item最小显示宽度, 1像素宽，1像素边距
	MAX_ITEM_WIDTH = 100

	ZOOM_OUT = 1.2
)

var (
	TRACK_LINE_COLOR = gui.NewQColor3(255, 255, 255, 0x9F)
	SELECT_RECT_COLOR = gui.NewQColor3(255, 255, 255, 0x9F)
	SEPARATOR_LINE_COLOR = gui.NewQColor3(255, 0, 0, 0x9F)
	DECORATOR_TEXT_COLOR = gui.NewQColor3(255, 0, 0, 255)
	VALUE_GRAPH_BG_COLOR = gui.NewQColor3(0, 0, 111, 255)
)

const (
	SEPARATOR_GAP_MIN = 30
	SEPARATOR_MAX = 6
)

const (
	Y_TICK_WIDTH = 4
	X_TICK_HEIGHT = 4
)