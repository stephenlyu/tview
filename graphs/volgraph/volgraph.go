package volgraph

import (
	"github.com/stephenlyu/tview/model"
	"github.com/stephenlyu/tds/util"
	"github.com/stephenlyu/tview/transform"
	"github.com/therecipe/qt/widgets"
	"github.com/stephenlyu/tview/constants"
	"github.com/therecipe/qt/gui"
	"github.com/stephenlyu/tview/graphs"
	"github.com/therecipe/qt/core"
)


type VolStickGraph struct {
	Model model.Model
	KLineModel model.Model
	ValueIndex int
	Scene *widgets.QGraphicsScene
	xTransformer transform.ScaleTransformer

	startIndex, endIndex int
	Lines map[int]*widgets.QGraphicsPathItem
}

func NewVolStickGraph(model model.Model, valueIndex int, klineModel model.Model, scene *widgets.QGraphicsScene, xTransformer transform.ScaleTransformer) *VolStickGraph {
	util.Assert(model != nil, "model != nil")
	util.Assert(len(model.GetGraphTypes()) > valueIndex, "len(model.GetGraphTypes()) > valueIndex")
	util.Assert(model.GetGraphTypes()[valueIndex] == constants.GraphTypeVolStick, "model.GetGraphTypes()[valueIndex] == constants.GraphTypeVolStick")

	this := &VolStickGraph{
		Model: model,
		KLineModel: klineModel,
		ValueIndex: valueIndex,
		Scene: scene,
		xTransformer: xTransformer,
		Lines: make(map[int]*widgets.QGraphicsPathItem),
	}
	this.init()
	return this
}

func (this *VolStickGraph) init() {
	this.Model.AddListener(this)
}

func (this *VolStickGraph) OnDataChanged() {
	for i, item := range this.Lines {
		if i >= this.Model.Count() {
			continue
		}
		this.updateStick(i, item)
	}
}

func (this *VolStickGraph) OnLastDataChanged() {
	if this.Model.Count() <= 0 {
		return
	}

	i := this.Model.Count() - 1
	item := this.ensureItem(i)
	this.updateStick(i, item)
}

func (this *VolStickGraph) GetValueRange(startIndex int, endIndex int) (float64, float64) {
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
	util.Assert(len(values) > this.ValueIndex, "len(values) > this.ValueIndex")

	high := values[this.ValueIndex]
	low := values[this.ValueIndex]

	for i := startIndex + 1; i < endIndex; i++ {
		values := this.Model.GetRaw(i)
		v := values[this.ValueIndex]
		if v > high {
			high = v
		}
		if v < low {
			low = v
		}
	}

	if low > 0 {
		low = 0
	}

	return low, high
}

func (this *VolStickGraph) ensureItem(i int) *widgets.QGraphicsPathItem {
	item, ok := this.Lines[i]
	if !ok {
		item = widgets.NewQGraphicsPathItem(nil)
		this.Scene.AddItem(item)
		this.Lines[i] = item
	}
	return item
}

func (this *VolStickGraph) updateStick(i int, item *widgets.QGraphicsPathItem) {
	x := (this.xTransformer.To(float64(i)) + this.xTransformer.To(float64(i + 2))) / 2
	w := this.xTransformer.To(1) / 3
	values := this.Model.Get(i)
	y := values[this.ValueIndex]

	path := gui.NewQPainterPath()

	path.AddRect2(x - w, 0, w * 2, y)

	klineValues := this.KLineModel.GetRaw(i)
	open := klineValues[0]
	close := klineValues[1]

	if close < open {
		item.SetBrush(gui.NewQBrush3(graphs.NegativeColor, core.Qt__SolidPattern))
		item.SetPen(gui.NewQPen3(graphs.NegativeColor))
	} else {
		item.SetBrush(gui.NewQBrush3(graphs.PositiveColor, core.Qt__SolidPattern))
		item.SetPen(gui.NewQPen3(graphs.PositiveColor))
	}

	item.SetPath(path)
}

func (this *VolStickGraph) adjustIndices(startIndex int, endIndex int) (int, int) {
	if startIndex > 0 {
		startIndex--
	}
	if endIndex < this.Model.Count() {
		endIndex++
	}
	return startIndex, endIndex
}

// 更新当前显示的K线
func (this *VolStickGraph) Update(startIndex int, endIndex int) {
	if startIndex < 0 {
		startIndex = 0
	}

	if endIndex > this.Model.Count() {
		endIndex = this.Model.Count()
	}

	startIndex, endIndex = this.adjustIndices(startIndex, endIndex)

	// 隐藏不需要显示的K线
	for i, item := range this.Lines {
		if i < startIndex || i >= endIndex {
			item.Hide()
		}
	}

	// 更新需要显示的K线
	for i := startIndex; i < endIndex; i++ {
		item := this.ensureItem(i)
		item.Show()
		this.updateStick(i, item)
	}
}

// 清除所有的K线
func (this *VolStickGraph) Clear() {
	for _, item := range this.Lines {
		this.Scene.RemoveItem(item)
	}
	this.Lines = make(map[int]*widgets.QGraphicsPathItem)
}
