package sticklinegraph

import (
	"github.com/stephenlyu/tview/model"
	"github.com/stephenlyu/tds/util"
	"github.com/stephenlyu/tview/transform"
	"github.com/therecipe/qt/widgets"
	"github.com/therecipe/qt/gui"
	"github.com/stephenlyu/tview/graphs"
	"math"
	"github.com/stephenlyu/goformula/function"
	"github.com/stephenlyu/goformula/formulalibrary/base/formula"
)

type StickLineGraph struct {
	Model model.Model
	DrawAction formula.StickLine
	Scene *widgets.QGraphicsScene
	xTransformer transform.ScaleTransformer

	color *gui.QColor
	startIndex, endIndex int
	Lines map[int]*widgets.QGraphicsPathItem
}

func NewStickLineGraph(model model.Model, DrawAction formula.StickLine, color *gui.QColor, scene *widgets.QGraphicsScene, xTransformer transform.ScaleTransformer) *StickLineGraph {
	util.Assert(model != nil, "model != nil")

	this := &StickLineGraph{
		Model: model,
		DrawAction: DrawAction,
		Scene: scene,
		xTransformer: xTransformer,
		color: color,
		Lines: make(map[int]*widgets.QGraphicsPathItem),
	}
	this.init()
	return this
}

func (this *StickLineGraph) init() {
	this.Model.AddListener(this)
}
func (this *StickLineGraph) OnDataChanged() {
	for i, item := range this.Lines {
		if i >= this.Model.Count() {
			continue
		}
		this.updateStick(i, item)
	}
}

func (this *StickLineGraph) OnLastDataChanged() {
	if this.Model.Count() <= 0 {
		return
	}

	i := this.Model.Count() - 1
	item := this.ensureItem(i)
	this.updateStick(i, item)
}

func (this *StickLineGraph) GetValueRange(startIndex int, endIndex int) (float64, float64) {
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

func (this *StickLineGraph) ensureItem(i int) *widgets.QGraphicsPathItem {
	item, ok := this.Lines[i]
	if !ok {
		item = widgets.NewQGraphicsPathItem(nil)
		this.Scene.AddItem(item)
		this.Lines[i] = item
	}
	return item
}

func (this *StickLineGraph) updateStick(i int, item *widgets.QGraphicsPathItem) bool {
	x := (this.xTransformer.To(float64(i)) + this.xTransformer.To(float64(i + 1))) / 2
	value1 := this.Model.Transform(this.DrawAction.GetPrice1(i))
	value2 := this.Model.Transform(this.DrawAction.GetPrice2(i))

	if function.IsNaN(value1) || function.IsNaN(value2) {
		return false
	}

	path := gui.NewQPainterPath()

	path.MoveTo2(x, value1)
	path.LineTo2(x, value2)

	pen := gui.NewQPen3(this.color)
	item.SetPen(pen)
	graphs.SetPenWidth(pen, this.xTransformer, this.DrawAction.GetLineThick())

	item.SetPath(path)
	return true
}

func (this *StickLineGraph) adjustIndices(startIndex int, endIndex int) (int, int) {
	if startIndex > 0 {
		startIndex--
	}
	if endIndex < this.Model.Count() {
		endIndex++
	}
	return startIndex, endIndex
}

// 更新当前显示的K线
func (this *StickLineGraph) Update(startIndex int, endIndex int) {
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

	// 隐藏不需要显示的K线
	for _, item := range this.Lines {
		item.Hide()
	}

	// 更新需要显示的K线
	for i := startIndex; i < endIndex; i++ {
		cond := this.DrawAction.GetCond(i)
		if cond == 0 {
			continue
		}

		item := this.ensureItem(i)
		if this.updateStick(i, item) {
			item.Show()
		} else {
			item.Hide()
		}
	}
}

// 清除所有的K线
func (this *StickLineGraph) Clear() {
	if this.DrawAction.IsNoDraw() {
		return
	}

	for _, item := range this.Lines {
		this.Scene.RemoveItem(item)
	}
	this.Lines = make(map[int]*widgets.QGraphicsPathItem)
}

func (this *StickLineGraph) ShowInfo(index int, display graphs.InfoDisplay) {
}
