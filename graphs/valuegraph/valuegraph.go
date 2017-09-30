package valuegraph

import (
	"github.com/therecipe/qt/widgets"
	"github.com/therecipe/qt/gui"
	"github.com/stephenlyu/tview/constants"
	"github.com/therecipe/qt/core"
)

type ValueGraph struct {
	scene *widgets.QGraphicsScene
	w, h float64
	pen *gui.QPen
	brush *gui.QBrush

	bgRect *widgets.QGraphicsRectItem
	textItem *widgets.QGraphicsTextItem
}

func NewValueGraph(scene *widgets.QGraphicsScene, w float64, h float64) *ValueGraph {
	return &ValueGraph{
		scene: scene,
		w: w,
		h: h,
		pen: gui.NewQPen3(constants.DECORATOR_TEXT_COLOR),
		brush: gui.NewQBrush3(constants.VALUE_GRAPH_BG_COLOR, core.Qt__SolidPattern),
	}
}

func (this *ValueGraph) MeasureText(value string) *core.QRectF {
	fm := gui.NewQFontMetricsF(this.scene.Font())
	return fm.BoundingRect(value)
}

func (this *ValueGraph) Update(x float64, y float64, value string) {
	this.Clear()

	this.bgRect = this.scene.AddRect2(x, y, this.w, this.h, this.pen, this.brush)

	ti := this.scene.AddText(value, this.scene.Font())
	ti.SetDefaultTextColor(constants.DECORATOR_TEXT_COLOR)
	ti.SetTransform(gui.QTransform_FromScale(1.0, -1.0), false)
	// TODO: DON'T KNOW WHY?
	ti.SetPos2(x + 1, y + this.h + 2)
	this.textItem = ti
}

func (this *ValueGraph) Clear() {
	if this.bgRect != nil {
		this.scene.RemoveItem(this.bgRect)
		this.bgRect = nil
	}

	if this.textItem != nil {
		this.scene.RemoveItem(this.textItem)
		this.textItem = nil
	}
}
