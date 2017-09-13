package klinegraph

import (
	"github.com/therecipe/qt/widgets"
	"github.com/therecipe/qt/gui"
	"math"
	"github.com/therecipe/qt/core"
)

type KLineItem struct {
	Item *widgets.QGraphicsPathItem
}

func NewKLineItem() *KLineItem {
	return &KLineItem{Item: widgets.NewQGraphicsPathItem(nil)}
}

func (this *KLineItem) Update(x float64, w float64, open float64, close float64, high float64, low float64) {
	this.Item.SetPos2(x + w / 2, (high + low) / 2)

	kWidth := w * 2.0 / 3
	min := math.Min(close, open)
	max := math.Max(close, open)

	maxY := (high - low) / 2
	topY := (max - low) - maxY
	bottomY := (min - low) - maxY

	path := gui.NewQPainterPath()
	if kWidth < 3 {
		path.MoveTo2(0, -maxY)
		path.LineTo2(0, maxY)
	} else {
		path.MoveTo2(0, -maxY)
		path.LineTo2(0, bottomY)
		path.AddRect(core.NewQRectF4(-kWidth/2, bottomY, kWidth, topY-bottomY))
		path.MoveTo2(0, topY)
		path.LineTo2(0, maxY)
	}

	var brush *gui.QBrush
	if close > open {
		brush = gui.NewQBrush3(gui.NewQColor3(0xD9, 0x11, 0x1B, 0xFF), core.Qt__SolidPattern)
	} else if close == open {
		brush = gui.NewQBrush4(core.Qt__white, core.Qt__SolidPattern)
	} else {
		brush = gui.NewQBrush3(gui.NewQColor3(0x42, 0xFF, 0xFF, 0xFF), core.Qt__SolidPattern)
	}

	this.Item.SetBrush(brush)
	this.Item.SetPen(gui.NewQPen3(brush.Color()))

	this.Item.SetPath(path)
}
