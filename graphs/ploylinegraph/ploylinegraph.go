package ploylinegraph

import (
	"github.com/stephenlyu/tview/model"
	"github.com/stephenlyu/tds/util"
	"github.com/stephenlyu/tview/transform"
	"github.com/therecipe/qt/widgets"
	"github.com/therecipe/qt/gui"
	"github.com/stephenlyu/goformula/function"
	"github.com/therecipe/qt/core"
	"github.com/stephenlyu/tview/graphs"
	"github.com/stephenlyu/goformula/stockfunc/formula"
)

type PloyLineGraph struct {
	Model model.Model
	DrawAction formula.PloyLine
	Scene *widgets.QGraphicsScene
	xTransformer transform.ScaleTransformer

	color *gui.QColor
	startIndex, endIndex int
	PathItem *widgets.QGraphicsPathItem
}

func NewPloyLineGraph(model model.Model, DrawAction formula.PloyLine, color *gui.QColor, scene *widgets.QGraphicsScene, xTransformer transform.ScaleTransformer) *PloyLineGraph {
	util.Assert(model != nil, "model != nil")

	this := &PloyLineGraph{
		Model: model,
		DrawAction: DrawAction,
		Scene: scene,
		xTransformer: xTransformer,
		color: color,
	}
	this.init()
	return this
}

func (this *PloyLineGraph) init() {
	this.Model.AddListener(this)
}

func (this *PloyLineGraph) OnDataChanged() {
	this.buildLine()
}

func (this *PloyLineGraph) OnLastDataChanged() {
	if this.Model.Count() <= 0 {
		return
	}

	this.buildLine()
}

func (this *PloyLineGraph) GetValueRange(startIndex int, endIndex int) (float64, float64) {
	if this.Model.Count() == 0 {
		return 0, 0
	}

	if startIndex < 0 {
		startIndex = 0
	}

	if endIndex > this.Model.Count() {
		endIndex = this.Model.Count()
	}

	startIndex, endIndex = this.adjustIndices(startIndex, endIndex)

	value := this.Model.TransformRaw(this.DrawAction.GetPrice(startIndex))

	high := value
	low := value

	for i := startIndex + 1; i < endIndex; i++ {
		v := this.Model.TransformRaw(this.DrawAction.GetPrice(i))
		if v > high {
			high = v
		}
		if v < low {
			low = v
		}
	}

	return low, high
}

func (this *PloyLineGraph) buildLine() {
	this.Clear()

	if this.Model.Count() == 0 {
		return
	}

	path := gui.NewQPainterPath()

	needMove := true
	for i := this.startIndex; i < this.endIndex; i++ {
		x := (this.xTransformer.To(float64(i)) + this.xTransformer.To(float64(i + 1))) / 2
		v := this.Model.Transform(this.DrawAction.GetPrice(i))

		cond := this.DrawAction.GetCond(i)

		if function.IsNaN(v) {
			needMove = true
			continue
		}

		if cond == 0 {
			continue
		}

		if needMove {
			path.MoveTo2(x, v)
			needMove = false
		} else {
			path.LineTo2(x, v)
		}
	}

	brush := gui.NewQBrush3(this.color, core.Qt__NoBrush)
	pen := gui.NewQPen3(this.color)
	graphs.SetPenWidth(pen, this.xTransformer, this.DrawAction.GetLineThick())

	this.PathItem = this.Scene.AddPath(path, pen, brush)
}

func (this *PloyLineGraph) adjustIndices(startIndex int, endIndex int) (int, int) {
	if startIndex > 0 {
		startIndex--
	}
	if endIndex < this.Model.Count() {
		endIndex++
	}
	return startIndex, endIndex
}

// 更新当前显示的K线
func (this *PloyLineGraph) Update(startIndex int, endIndex int) {
	if this.DrawAction.IsNoDraw() {
		return
	}
	if startIndex < 0 {
		startIndex = 0
	}

	if endIndex > this.Model.Count() {
		endIndex = this.Model.Count()
	}

	startIndex, endIndex = this.adjustIndices(startIndex, endIndex)
	this.startIndex = startIndex
	this.endIndex = endIndex

	// 更新需要显示的K线
	this.buildLine()
}

// 清除所有的K线
func (this *PloyLineGraph) Clear() {
	if this.DrawAction.IsNoDraw() {
		return
	}
	if this.PathItem != nil {
		this.Scene.RemoveItem(this.PathItem)
		this.PathItem = nil
	}
}

func (this *PloyLineGraph) ShowInfo(index int, display graphs.InfoDisplay) {
}
