package drawlinegraph

import (
	"github.com/stephenlyu/tview/model"
	"github.com/stephenlyu/tds/util"
	"github.com/stephenlyu/tview/transform"
	"github.com/therecipe/qt/widgets"
	"github.com/therecipe/qt/gui"
	"github.com/stephenlyu/goformula/function"
	"github.com/therecipe/qt/core"
	"github.com/stephenlyu/tview/graphs"
	"math"
	"github.com/stephenlyu/goformula/formulalibrary/base/formula"
)

type DrawLineGraph struct {
	Model model.Model
	DrawAction formula.DrawLine
	Scene *widgets.QGraphicsScene
	xTransformer transform.ScaleTransformer

	color *gui.QColor
	startIndex, endIndex int
	PathItem *widgets.QGraphicsPathItem
}

func NewDrawLineGraph(model model.Model, DrawAction formula.DrawLine, color *gui.QColor, scene *widgets.QGraphicsScene, xTransformer transform.ScaleTransformer) *DrawLineGraph {
	util.Assert(model != nil, "model != nil")

	this := &DrawLineGraph{
		Model: model,
		DrawAction: DrawAction,
		Scene: scene,
		xTransformer: xTransformer,
		color: color,
	}
	this.init()
	return this
}

func (this *DrawLineGraph) init() {
	this.Model.AddListener(this)
}

func (this *DrawLineGraph) OnDataChanged() {
	this.buildLine()
}

func (this *DrawLineGraph) OnLastDataChanged() {
	if this.Model.Count() <= 0 {
		return
	}

	this.buildLine()
}

func (this *DrawLineGraph) GetValueRange(startIndex int, endIndex int) (float64, float64) {
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

	value1 := this.Model.TransformRaw(this.DrawAction.GetPrice1(startIndex))
	value2 := this.Model.TransformRaw(this.DrawAction.GetPrice2(startIndex))

	high := math.Max(value1, value2)
	low := math.Min(value1, value2)

	for i := startIndex + 1; i < endIndex; i++ {
		value1 = this.Model.TransformRaw(this.DrawAction.GetPrice1(i))
		value2 = this.Model.TransformRaw(this.DrawAction.GetPrice2(i))

		highV := math.Max(value1, value2)
		lowV := math.Min(value1, value2)

		if highV > high {
			high = highV
		}
		if lowV < low {
			low = lowV
		}
	}

	return low, high
}

func (this *DrawLineGraph) buildLine() {
	this.Clear()

	if this.Model.Count() == 0 {
		return
	}

	path := gui.NewQPainterPath()

	var prevCond1X, prevCond1Y float64
	var prevCond2X, prevCond2Y float64

	prevCond1X = -1
	prevCond2X = -1

	for i := this.startIndex; i < this.endIndex; i++ {
		x := (this.xTransformer.To(float64(i)) + this.xTransformer.To(float64(i + 1))) / 2
		value1 := this.Model.Transform(this.DrawAction.GetPrice1(i))

		cond1 := this.DrawAction.GetCond1(i)
		if function.IsTrue(cond1) {
			if prevCond2X > prevCond1X && !function.IsNaN(prevCond2Y) {
				path.MoveTo2(prevCond1X, prevCond1Y)
				path.LineTo2(prevCond2X, prevCond2Y)
				prevCond2X = -1
			}

			prevCond1X = x
			prevCond1Y = value1
			continue
		}

		if prevCond1X < 0 || function.IsNaN(prevCond1Y) {
			continue
		}

		cond2 := this.DrawAction.GetCond2(i)
		value2 := this.Model.Transform(this.DrawAction.GetPrice2(i))
		if function.IsTrue(cond2) {
			prevCond2X = x
			prevCond2Y = value2
		}
	}
	if prevCond2X > prevCond1X && !function.IsNaN(prevCond2Y) {
		path.MoveTo2(prevCond1X, prevCond1Y)
		path.LineTo2(prevCond2X, prevCond2Y)
		prevCond2X = -1
	}

	brush := gui.NewQBrush3(this.color, core.Qt__NoBrush)
	pen := gui.NewQPen3(this.color)
	graphs.SetPenWidth(pen, this.xTransformer, this.DrawAction.GetLineThick())

	this.PathItem = this.Scene.AddPath(path, pen, brush)
}

func (this *DrawLineGraph) adjustIndices(startIndex int, endIndex int) (int, int) {
	for startIndex > 0 {
		startIndex--
		cond1 := this.DrawAction.GetCond1(startIndex)
		if cond1 != 0 {
			break
		}
	}

	for endIndex < this.Model.Count() {
		cond2 := this.DrawAction.GetCond2(endIndex)
		if cond2 != 0 {
			break
		}
		endIndex++
	}

	return startIndex, endIndex
}

// 更新当前显示的K线
func (this *DrawLineGraph) Update(startIndex int, endIndex int) {
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
func (this *DrawLineGraph) Clear() {
	if this.DrawAction.IsNoDraw() {
		return
	}
	if this.PathItem != nil {
		this.Scene.RemoveItem(this.PathItem)
		this.PathItem = nil
	}
}

func (this *DrawLineGraph) ShowInfo(index int, display graphs.InfoDisplay) {
}
