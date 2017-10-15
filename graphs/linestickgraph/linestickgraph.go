package linestickgraph

import (
	"github.com/stephenlyu/tview/model"
	"github.com/stephenlyu/tds/util"
	"github.com/stephenlyu/tview/transform"
	"github.com/therecipe/qt/widgets"
	"github.com/stephenlyu/tview/constants"
	"github.com/therecipe/qt/gui"
	"github.com/stephenlyu/tview/graphs"
	"fmt"
	"github.com/stephenlyu/goformula/function"
	"github.com/therecipe/qt/core"
)


type LineStickGraph struct {
	Model model.Model
	ValueIndex int
	Color *gui.QColor
	Scene *widgets.QGraphicsScene
	xTransformer transform.ScaleTransformer

	startIndex, endIndex int
	XAxis *widgets.QGraphicsLineItem
	Lines map[int]*widgets.QGraphicsPathItem
	PathItem *widgets.QGraphicsPathItem
}

func NewLineStickGraph(model model.Model, valueIndex int, color *gui.QColor, scene *widgets.QGraphicsScene, xTransformer transform.ScaleTransformer) *LineStickGraph {
	util.Assert(model != nil, "model != nil")
	util.Assert(model.VarCount() > valueIndex, "len(model.GetGraphTypes()) > valueIndex")
	util.Assert(model.GraphType(valueIndex) == constants.GraphTypeLineStick, "model.GetGraphTypes()[valueIndex] == constants.GraphTypeLineStick")

	this := &LineStickGraph{
		Model: model,
		ValueIndex: valueIndex,
		Color: color,
		Scene: scene,
		xTransformer: xTransformer,
		Lines: make(map[int]*widgets.QGraphicsPathItem),
	}
	this.init()
	return this
}

func (this *LineStickGraph) init() {
	this.Model.AddListener(this)
}

func (this *LineStickGraph) OnDataChanged() {
	for i, item := range this.Lines {
		if i >= this.Model.Count() {
			continue
		}
		this.updateStick(i, item)
	}
}

func (this *LineStickGraph) OnLastDataChanged() {
	if this.Model.Count() <= 0 {
		return
	}

	i := this.Model.Count() - 1
	item := this.ensureItem(i)
	this.updateStick(i, item)
}

func (this *LineStickGraph) GetValueRange(startIndex int, endIndex int) (float64, float64) {
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

	// X轴
	if low > 0 {
		low = 0
	}

	return low, high
}

func (this *LineStickGraph) ensureItem(i int) *widgets.QGraphicsPathItem {
	item, ok := this.Lines[i]
	if !ok {
		item = widgets.NewQGraphicsPathItem(nil)
		this.Scene.AddItem(item)
		this.Lines[i] = item
	}
	return item
}

func (this *LineStickGraph) buildLine() {
	fmt.Println("buildLine")
	this.clearLine()

	if this.Model.Count() == 0 {
		return
	}

	path := gui.NewQPainterPath()

	needMove := true
	for i := this.startIndex; i < this.endIndex; i++ {
		x := (this.xTransformer.To(float64(i)) + this.xTransformer.To(float64(i + 1))) / 2
		values := this.Model.Get(i)
		v := values[this.ValueIndex]

		if function.IsNaN(v) {
			needMove = true
			continue
		}

		if needMove {
			path.MoveTo2(x, v)
			needMove = false
		} else {
			path.LineTo2(x, v)
		}
	}

	brush := gui.NewQBrush3(this.Color, core.Qt__NoBrush)
	pen := gui.NewQPen3(this.Color)

	this.PathItem = this.Scene.AddPath(path, pen, brush)
}

func (this *LineStickGraph) updateStick(i int, item *widgets.QGraphicsPathItem) {
	x := (this.xTransformer.To(float64(i)) + this.xTransformer.To(float64(i + 1))) / 2
	y := this.Model.Get(i)[this.ValueIndex]

	path := gui.NewQPainterPath()

	path.MoveTo2(x, 0)
	path.LineTo2(x, y)

	item.SetPen(gui.NewQPen3(this.Color))

	item.SetPath(path)
}

func (this *LineStickGraph) adjustIndices(startIndex int, endIndex int) (int, int) {
	if startIndex > 0 {
		startIndex--
	}
	if endIndex < this.Model.Count() {
		endIndex++
	}
	return startIndex, endIndex
}

// 更新当前显示的K线
func (this *LineStickGraph) Update(startIndex int, endIndex int) {
	if this.Model.NoDraw(this.ValueIndex) {
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

	if this.XAxis != nil {
		this.Scene.RemoveItem(this.XAxis)
		this.XAxis = nil
	}
	x1 := this.xTransformer.To(float64(startIndex))
	x2 := this.xTransformer.To(float64(endIndex))
	this.XAxis = this.Scene.AddLine2(x1, 0, x2, 0, gui.NewQPen3(gui.NewQColor3(255, 255, 255, 255)))
}

func (this *LineStickGraph) clearLine() {
	if this.PathItem != nil {
		this.Scene.RemoveItem(this.PathItem)
		this.PathItem = nil
	}
}

// 清除所有的K线
func (this *LineStickGraph) Clear() {
	if this.Model.NoDraw(this.ValueIndex) {
		return
	}
	for _, item := range this.Lines {
		this.Scene.RemoveItem(item)
	}
	if this.XAxis != nil {
		this.Scene.RemoveItem(this.XAxis)
		this.XAxis = nil
	}
	this.clearLine()

	this.Lines = make(map[int]*widgets.QGraphicsPathItem)
}

func (this *LineStickGraph) ShowInfo(index int, display graphs.InfoDisplay) {
	if this.Model.NoText(this.ValueIndex) {
		return
	}
	if index < 0 || index >= this.Model.Count() {
		return
	}

	name := this.Model.GetNames()[this.ValueIndex]
	v := this.Model.GetRaw(index)[this.ValueIndex]
	display.Add(fmt.Sprintf("%s: %.02f", name, v), this.Color)
}
