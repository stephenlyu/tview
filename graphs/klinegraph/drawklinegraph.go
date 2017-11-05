package klinegraph

import (
	"github.com/stephenlyu/tview/model"
	"github.com/stephenlyu/tds/util"
	"github.com/stephenlyu/tview/transform"
	"github.com/therecipe/qt/widgets"
	"github.com/stephenlyu/tview/graphs"
	"github.com/stephenlyu/goformula/stockfunc/formula"
)

type DrawKLineGraph struct {
	Model model.Model
	DrawAction formula.DrawKLine
	Scene *widgets.QGraphicsScene
	xTransformer transform.ScaleTransformer

	KLines map[int]*KLineItem
}

func NewDrawKLineGraph(model model.Model, DrawAction formula.DrawKLine, scene *widgets.QGraphicsScene, xTransformer transform.ScaleTransformer) *DrawKLineGraph {
	util.Assert(model != nil, "model != nil")

	this := &DrawKLineGraph{
		Model: model,
		DrawAction: DrawAction,
		Scene: scene,
		xTransformer: xTransformer,
		KLines: make(map[int]*KLineItem),
	}
	this.init()
	return this
}

func (this *DrawKLineGraph) init() {
	this.Model.AddListener(this)
}

func (this *DrawKLineGraph) OnDataChanged() {
	for i, item := range this.KLines {
		if i >= this.Model.Count() {
			continue
		}
		this.updateKLine(i, item)
	}
}

func (this *DrawKLineGraph) OnLastDataChanged() {
	if this.Model.Count() <= 0 {
		return
	}

	i := this.Model.Count() - 1
	item := this.ensureItem(i)
	this.updateKLine(i, item)
}

func (this *DrawKLineGraph) GetValueRange(startIndex int, endIndex int) (float64, float64) {
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

	high := this.Model.TransformRaw(this.DrawAction.GetHigh(startIndex))
	low := this.Model.TransformRaw(this.DrawAction.GetLow(startIndex))

	for i := startIndex + 1; i < endIndex; i++ {
		vHigh, vLow := this.Model.TransformRaw(this.DrawAction.GetHigh(i)),
			this.Model.TransformRaw(this.DrawAction.GetLow(i))

		if vHigh > high {
			high = vHigh
		}
		if vLow < low {
			low = vLow
		}
	}

	return low, high
}

func (this *DrawKLineGraph) ensureItem(i int) *KLineItem {
	item, ok := this.KLines[i]
	if !ok {
		item = NewKLineItem()
		this.Scene.AddItem(item.Item)
		this.KLines[i] = item
	}
	return item
}

func (this *DrawKLineGraph) updateKLine(i int, item *KLineItem) {
	x := this.xTransformer.To(float64(i))
	w := this.xTransformer.To(1)

	open, close, high, low := this.Model.Transform(this.DrawAction.GetOpen(i)),
		this.Model.Transform(this.DrawAction.GetClose(i)),
		this.Model.Transform(this.DrawAction.GetHigh(i)),
		this.Model.Transform(this.DrawAction.GetLow(i))
	item.Update1(x, w, open, close, high, low)
}

func (this *DrawKLineGraph) adjustIndices(startIndex int, endIndex int) (int, int) {
	if startIndex > 0 {
		startIndex--
	}
	if endIndex < this.Model.Count() {
		endIndex++
	}
	return startIndex, endIndex
}

// 更新当前显示的K线
func (this *DrawKLineGraph) Update(startIndex int, endIndex int) {
	if startIndex < 0 {
		startIndex = 0
	}

	if endIndex > this.Model.Count() {
		endIndex = this.Model.Count()
	}

	startIndex, endIndex = this.adjustIndices(startIndex, endIndex)

	// 隐藏不需要显示的K线
	for i, item := range this.KLines {
		if i < startIndex || i >= endIndex {
			item.Item.Hide()
		}
	}

	// 更新需要显示的K线
	for i := startIndex; i < endIndex; i++ {
		item := this.ensureItem(i)
		item.Item.Show()
		this.updateKLine(i, item)
	}
}

// 清除所有的K线
func (this *DrawKLineGraph) Clear() {
	for _, item := range this.KLines {
		this.Scene.RemoveItem(item.Item)
	}
	this.KLines = make(map[int]*KLineItem)
}

func (this *DrawKLineGraph) ShowInfo(index int, display graphs.InfoDisplay) {
}