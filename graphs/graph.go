package graphs

import (
	"github.com/therecipe/qt/gui"
	"github.com/stephenlyu/tview/transform"
	"github.com/therecipe/qt/core"
	"math"
	"fmt"
	"github.com/stephenlyu/goformula/formulalibrary/base/formula"
)

type InfoDisplay interface {
	Add(text string, color *gui.QColor)
}

type Graph interface {
	GetValueRange(startIndex int, endIndex int) (float64, float64)
	Update(startIndex int, endIndex int)
	Clear()

	ShowInfo(index int, display InfoDisplay)
}

func getLineWidth(xTransformer transform.ScaleTransformer, lineThick int) int {
	maxWidth := xTransformer.To(1) * 2 / 3

	width := maxWidth * float64(lineThick) / 9
	if width < 1.0 {
		width = 1.0
	} else if width > float64(lineThick) {
		width = float64(lineThick)
	}
	return int(width)
}

func SetPenWidth(pen *gui.QPen, xTransformer transform.ScaleTransformer, lineThick int) {
	pen.SetWidth(getLineWidth(xTransformer, lineThick))
}

func SetPenStyle(pen *gui.QPen, style int) {
	switch style {
	case formula.FORMULA_LINE_STYLE_CIRCLE_DOT:
		pen.SetStyle(core.Qt__DashDotLine)
	case formula.FORMULA_LINE_STYLE_CROSS_DOT:
		pen.SetStyle(core.Qt__DashLine)
	case formula.FORMULA_LINE_STYLE_DOT:
		pen.SetStyle(core.Qt__DotLine)
	case formula.FORMULA_LINE_STYLE_POINT_DOT:
		pen.SetStyle(core.Qt__DashDotDotLine)
	case formula.FORMULA_LINE_STYLE_SOLID:
		pen.SetStyle(core.Qt__SolidLine)
	}
}

func FormatValue(v float64, bits int) string {
	if math.IsNaN(v) {
		return "-"
	}

	switch bits {
	case 1:
		return fmt.Sprintf("%.1f", v)
	case 2:
		return fmt.Sprintf("%.2f", v)
	case 3:
		return fmt.Sprintf("%.3f", v)
	default:
		return fmt.Sprintf("%.0f", v)
	}
}
