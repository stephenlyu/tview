package drawtextgraph

import (
	"github.com/stephenlyu/tview/model"
	"github.com/stephenlyu/tds/util"
	"github.com/stephenlyu/tview/transform"
	"github.com/therecipe/qt/widgets"
	"github.com/therecipe/qt/gui"
	"github.com/stephenlyu/tview/graphs"
	"github.com/stephenlyu/goformula/stockfunc/formula"
	"github.com/stephenlyu/goformula/function"
)

type DrawTextGraph struct {
	Model                model.Model
	DrawAction           formula.DrawText
	Scene                *widgets.QGraphicsScene
	xTransformer         transform.ScaleTransformer

	Trans 				gui.QTransform_ITF
	color                *gui.QColor
	startIndex, endIndex int
	Texts                map[int]*widgets.QGraphicsTextItem
}

func NewDrawTextGraph(model model.Model, DrawAction formula.DrawText, color *gui.QColor, scene *widgets.QGraphicsScene, xTransformer transform.ScaleTransformer) *DrawTextGraph {
	util.Assert(model != nil, "model != nil")

	this := &DrawTextGraph{
		Model: model,
		DrawAction: DrawAction,
		Scene: scene,
		xTransformer: xTransformer,
		color: color,
		Texts: make(map[int]*widgets.QGraphicsTextItem),
		Trans: gui.QTransform_FromScale(1.0, -1.0),
	}
	this.init()
	return this
}

func (this *DrawTextGraph) init() {
	this.Model.AddListener(this)
}
func (this *DrawTextGraph) OnDataChanged() {
	for i, item := range this.Texts {
		if i >= this.Model.Count() {
			continue
		}
		this.updateText(i, item)
	}
}

func (this *DrawTextGraph) OnLastDataChanged() {
	if this.Model.Count() <= 0 {
		return
	}

	i := this.Model.Count() - 1
	item := this.ensureItem(i)
	this.updateText(i, item)
}

func (this *DrawTextGraph) GetValueRange(startIndex int, endIndex int) (float64, float64) {
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
		value = this.Model.TransformRaw(this.DrawAction.GetPrice(startIndex))

		if value > high {
			high = value
		}
		if value < low {
			low = value
		}
	}

	// 保留20个像素绘制文本
	low -= this.Model.TransformRawFrom(20)

	return low, high
}

func (this *DrawTextGraph) ensureItem(i int) *widgets.QGraphicsTextItem {
	item, ok := this.Texts[i]
	if !ok {
		item = this.Scene.AddText("", this.Scene.Font())
		item.SetDefaultTextColor(this.color)
		item.SetTransform(this.Trans, false)
		this.Texts[i] = item
	}
	return item
}

func (this *DrawTextGraph) updateText(i int, item *widgets.QGraphicsTextItem) bool {
	x := (this.xTransformer.To(float64(i)) + this.xTransformer.To(float64(i + 1))) / 2
	value := this.Model.Transform(this.DrawAction.GetPrice(i))

	if function.IsNaN(value) {
		return false
	}

	item.SetPlainText(this.DrawAction.GetText())
	item.SetPos2(x - item.BoundingRect().Width() / 2, value)

	return true
}

func (this *DrawTextGraph) adjustIndices(startIndex int, endIndex int) (int, int) {
	if startIndex > 0 {
		startIndex--
	}
	if endIndex < this.Model.Count() {
		endIndex++
	}
	return startIndex, endIndex
}

// 更新当前显示的K线
func (this *DrawTextGraph) Update(startIndex int, endIndex int) {
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
	for _, item := range this.Texts {
		item.Hide()
	}

	// 更新需要显示的K线
	for i := startIndex; i < endIndex; i++ {
		cond := this.DrawAction.GetCond(i)
		if cond == 0 {
			continue
		}

		item := this.ensureItem(i)
		if this.updateText(i, item) {
			item.Show()
		} else {
			item.Hide()
		}
	}
}

// 清除所有的K线
func (this *DrawTextGraph) Clear() {
	if this.DrawAction.IsNoDraw() {
		return
	}

	for _, item := range this.Texts {
		this.Scene.RemoveItem(item)
	}
	this.Texts = make(map[int]*widgets.QGraphicsTextItem)
}

func (this *DrawTextGraph) ShowInfo(index int, display graphs.InfoDisplay) {
}
