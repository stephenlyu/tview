package selectrect

import (
	"github.com/therecipe/qt/widgets"
	"github.com/therecipe/qt/gui"
	"github.com/stephenlyu/tview/constants"
	"github.com/therecipe/qt/core"
)

type SelectRect struct {
	Scene *widgets.QGraphicsScene
	pen *gui.QPen

	Item *widgets.QGraphicsRectItem
}

func NewSelectRect(scene *widgets.QGraphicsScene) *SelectRect {
	pen := gui.NewQPen3(constants.SELECT_RECT_COLOR)
	pen.SetWidth(1)
	pen.SetStyle(core.Qt__SolidLine)
	return &SelectRect{
		Scene: scene,
		pen: pen,
	}
}

func (this *SelectRect) UpdateRect(x float64, y float64, w float64, h float64) {
	this.Clear()

	if this.Item == nil {
		this.Item = this.Scene.AddRect2(x, y, w, h, this.pen, gui.NewQBrush2(core.Qt__NoBrush))
	} else {
		this.Item.Update2(x, y, w, h)
	}
}

func (this *SelectRect) GetRect() *core.QRectF {
	if this.Item == nil {
		return nil
	}
	return this.Item.Rect()
}

func (this *SelectRect) Clear() {
	if this.Item != nil {
		this.Scene.RemoveItem(this.Item)
		this.Item = nil
	}
}
