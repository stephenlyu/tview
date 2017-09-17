package constants

import "github.com/therecipe/qt/gui"

type GraphType int
const (
	GraphTypeKLine = iota		// Kline, consume 4 values
	GraphTypeLine				// Line graph, consume 1 value
	GraphTypeStick				// Line graph, consume 1 value
	GraphTypeVolStick                // Line graph, consume 1 value
)


const (
	VISIBLE_KLINES_MIN = 16		// 最少可见K线数
	BEST_ITEM_WIDTH = 10		// Item最佳显示宽度
)

const (
	H_MARGIN = 10
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
)

const (
	SEPARATOR_GAP_MIN = 30
	SEPARATOR_MAX = 6
)