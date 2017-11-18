package drawicongraph

import (
	"github.com/stephenlyu/tview/model"
	"github.com/stephenlyu/tds/util"
	"github.com/stephenlyu/tview/transform"
	"github.com/therecipe/qt/widgets"
	"github.com/therecipe/qt/gui"
	"github.com/stephenlyu/tview/graphs"
	"github.com/stephenlyu/goformula/function"
	"fmt"
	"github.com/therecipe/qt/core"
	"github.com/stephenlyu/goformula/formulalibrary/base/formula"
)

type DrawIconGraph struct {
	Model                model.Model
	DrawAction           formula.DrawIcon
	Scene                *widgets.QGraphicsScene
	xTransformer         transform.ScaleTransformer

	Trans                gui.QTransform_ITF
	color                *gui.QColor
	startIndex, endIndex int
	Icons                map[int]*widgets.QGraphicsPixmapItem
}

func NewDrawIconGraph(model model.Model, DrawAction formula.DrawIcon, color *gui.QColor, scene *widgets.QGraphicsScene, xTransformer transform.ScaleTransformer) *DrawIconGraph {
	util.Assert(model != nil, "model != nil")

	this := &DrawIconGraph{
		Model: model,
		DrawAction: DrawAction,
		Scene: scene,
		xTransformer: xTransformer,
		color: color,
		Icons: make(map[int]*widgets.QGraphicsPixmapItem),
		Trans: gui.QTransform_FromScale(1.0, -1.0),
	}
	this.init()
	return this
}

func (this *DrawIconGraph) init() {
	this.Model.AddListener(this)
}
func (this *DrawIconGraph) OnDataChanged() {
	for i, item := range this.Icons {
		if i >= this.Model.Count() {
			continue
		}
		this.updateIcon(i, item)
	}
}

func (this *DrawIconGraph) OnLastDataChanged() {
	if this.Model.Count() <= 0 {
		return
	}

	i := this.Model.Count() - 1
	item := this.ensureItem(i)
	this.updateIcon(i, item)
}

func (this *DrawIconGraph) GetValueRange(startIndex int, endIndex int) (float64, float64) {
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

	// 保留20个像素绘制Icon
	low -= this.Model.TransformRawFrom(20)

	return low, high
}

func (this *DrawIconGraph) ensureItem(i int) *widgets.QGraphicsPixmapItem {
	item, ok := this.Icons[i]
	if !ok {
		item = widgets.NewQGraphicsPixmapItem(nil)
		this.Scene.AddItem(item)
		item.SetTransform(this.Trans, false)
		this.Icons[i] = item
	}
	return item
}

func (this *DrawIconGraph) updateIcon(i int, item *widgets.QGraphicsPixmapItem) bool {
	x := (this.xTransformer.To(float64(i)) + this.xTransformer.To(float64(i + 1))) / 2
	value := this.Model.Transform(this.DrawAction.GetPrice(i))

	if function.IsNaN(value) {
		return false
	}

	iconPath := fmt.Sprintf(":/%d.png", (this.DrawAction.GetType() - 1) % 41)

	pixmap := gui.NewQPixmap()
	if !pixmap.Load(iconPath, "png", core.Qt__AutoColor) {
		fmt.Println("load icon fail")
	}
	item.SetPixmap(pixmap)
	item.SetPos2(x - item.BoundingRect().Width() / 2, value)

	return true
}

func (this *DrawIconGraph) adjustIndices(startIndex int, endIndex int) (int, int) {
	if startIndex > 0 {
		startIndex--
	}
	if endIndex < this.Model.Count() {
		endIndex++
	}
	return startIndex, endIndex
}

// 更新当前显示的K线
func (this *DrawIconGraph) Update(startIndex int, endIndex int) {
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
	for _, item := range this.Icons {
		item.Hide()
	}

	// 更新需要显示的K线
	for i := startIndex; i < endIndex; i++ {
		cond := this.DrawAction.GetCond(i)
		if !function.IsTrue(cond) {
			continue
		}

		item := this.ensureItem(i)
		if this.updateIcon(i, item) {
			item.Show()
		} else {
			item.Hide()
		}
	}
}

// 清除所有的K线
func (this *DrawIconGraph) Clear() {
	if this.DrawAction.IsNoDraw() {
		return
	}

	for _, item := range this.Icons {
		this.Scene.RemoveItem(item)
	}
	this.Icons = make(map[int]*widgets.QGraphicsPixmapItem)
}

func (this *DrawIconGraph) ShowInfo(index int, display graphs.InfoDisplay) {
}
