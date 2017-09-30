package klinegraph

import (
	"github.com/stephenlyu/tview/model"
	"github.com/stephenlyu/tds/util"
	"github.com/stephenlyu/tview/constants"
	"github.com/stephenlyu/tview/transform"
	"github.com/therecipe/qt/widgets"
	"github.com/stephenlyu/tview/graphs"
)

type KLineGraph struct {
	Model model.Model
	Scene *widgets.QGraphicsScene
	xTransformer transform.ScaleTransformer

	KLines map[int]*KLineItem
}

func NewKLineGraph(model model.Model, scene *widgets.QGraphicsScene, xTransformer transform.ScaleTransformer) *KLineGraph {
	util.Assert(model != nil, "model != nil")
	util.Assert(model.VarCount() == 1, "bad model graph types")
	util.Assert(model.GraphType(0) == constants.GraphTypeKLine, "bad model graph types")

	this := &KLineGraph{
		Model: model,
		Scene: scene,
		xTransformer: xTransformer,
		KLines: make(map[int]*KLineItem),
	}
	this.init()
	return this
}

func (this *KLineGraph) init() {
	this.Model.AddListener(this)
}

func (this *KLineGraph) OnDataChanged() {
	for i, item := range this.KLines {
		if i >= this.Model.Count() {
			continue
		}
		this.updateKLine(i, item)
	}
}

func (this *KLineGraph) OnLastDataChanged() {
	if this.Model.Count() <= 0 {
		return
	}

	i := this.Model.Count() - 1
	item := this.ensureItem(i)
	this.updateKLine(i, item)
}

func (this *KLineGraph) GetValueRange(startIndex int, endIndex int) (float64, float64) {
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

	values := this.Model.GetRaw(startIndex)
	util.Assert(len(values) == 4, "bad model values")

	high := values[2]
	low := values[3]

	for i := startIndex + 1; i < endIndex; i++ {
		values := this.Model.GetRaw(i)
		if values[2] > high {
			high = values[2]
		}
		if values[3] < low {
			low = values[3]
		}
	}

	return low, high
}

func (this *KLineGraph) ensureItem(i int) *KLineItem {
	item, ok := this.KLines[i]
	if !ok {
		item = NewKLineItem()
		this.Scene.AddItem(item.Item)
		this.KLines[i] = item
	}
	return item
}

func (this *KLineGraph) updateKLine(i int, item *KLineItem) {
	x := this.xTransformer.To(float64(i))
	w := this.xTransformer.To(1)

	values := this.Model.Get(i)
	util.Assert(len(values) == 4, "bad model values")
	item.Update(x, w, values[0], values[1], values[2], values[3])
}

func (this *KLineGraph) adjustIndices(startIndex int, endIndex int) (int, int) {
	if startIndex > 0 {
		startIndex--
	}
	if endIndex < this.Model.Count() {
		endIndex++
	}
	return startIndex, endIndex
}

// 更新当前显示的K线
func (this *KLineGraph) Update(startIndex int, endIndex int) {
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
func (this *KLineGraph) Clear() {
	for _, item := range this.KLines {
		this.Scene.RemoveItem(item.Item)
	}
	this.KLines = make(map[int]*KLineItem)
}

func (this *KLineGraph) ShowInfo(index int, display graphs.InfoDisplay) {
}