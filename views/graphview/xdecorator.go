package graphview

import (
	"github.com/therecipe/qt/widgets"
	"github.com/stephenlyu/tview/transform"
	"github.com/stephenlyu/tview/model"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/core"
	"github.com/stephenlyu/tview/constants"
	"github.com/stephenlyu/tview/graphs/valuegraph"
	"github.com/cznic/mathutil"
)

//go:generate qtmoc
type XDecorator struct {
	widgets.QGraphicsView

	Pen *gui.QPen

	Transformer transform.ScaleTransformer
	Model *model.KLineModel

	Items []widgets.QGraphicsItem_ITF
	ValueGraph *valuegraph.ValueGraph

	FirstVisibleIndex, LastVisibleIndex int
	ItemWidth float64									// 每个数据占用的屏幕宽度
}

func CreateXDecorator(parent widgets.QWidget_ITF) *XDecorator {
	ret := NewYDecorator(parent)
	ret.Transformer = transform.NewLogicTransformer(1)
	ret.init()
	return ret
}

func (this *XDecorator) init() {
	this.SetRenderHint(gui.QPainter__TextAntialiasing, true)
	this.SetCacheMode(widgets.QGraphicsView__CacheBackground)
	this.SetHorizontalScrollBarPolicy(core.Qt__ScrollBarAlwaysOff)
	this.SetVerticalScrollBarPolicy(core.Qt__ScrollBarAlwaysOff)
	this.SetDragMode(widgets.QGraphicsView__NoDrag)
	this.SetAlignment(core.Qt__AlignLeft)
	this.Scale(1, -1)

	// 设置scene
	scene := widgets.NewQGraphicsScene(this)
	scene.SetBackgroundBrush(gui.NewQBrush4(core.Qt__black, core.Qt__SolidPattern))
	this.SetScene(scene)

	this.Pen = gui.NewQPen3(constants.DECORATOR_TEXT_COLOR)
	this.ValueGraph = valuegraph.NewValueGraph(this.Scene(), float64(this.Width()), VALUE_GRAPH_HEIGHT)

	this.ConnectWheelEvent(this.WheelEvent)
}

func (this *XDecorator) SetTransformer(transformer transform.ScaleTransformer) {
	this.Transformer = transformer
}

func (this *XDecorator) SetModel(model *model.KLineModel) {
	this.Model = model
}



func (this *XDecorator) drawUI() {
	yMin = this.Transformer.To(yMin)
	yMax = this.Transformer.To(yMax)
	this.SetSceneRect2(0, yMin - constants.V_MARGIN, float64(this.Width()), yMax - yMin + 2 * constants.V_MARGIN)
	this.CenterOn2(float64(this.Width()) / 2, (yMax + yMin) / 2)

	trans := gui.QTransform_FromScale(1.0, -1.0)
	for i := 0; i < this.Model.Count(); i++ {
		rv := this.Model.Get(i)
		v := this.Transformer.To(rv)
		tick := this.Scene().AddLine2(0, v, constants.Y_TICK_WIDTH, v, this.Pen)
		this.Items = append(this.Items, tick)

		text := this.formatValue(rv)
		ti := this.Scene().AddText(text, this.Font())
		ti.SetDefaultTextColor(constants.DECORATOR_TEXT_COLOR)
		ti.AdjustSize()
		ti.SetTransform(trans, false)
		r := ti.BoundingRect()
		x := float64(constants.Y_TICK_WIDTH + 2)
		y := v + r.Height() / 2
		ti.SetPos2(x, y)
		this.Items = append(this.Items, ti)
	}
}

func (this *XDecorator) UpdateUI() {
	this.Clear()

	width := float64(this.Model.Count()) * this.ItemWidth
	usableWidth := this.Width() - 2 * H_MARGIN
	fullMode := width < float64(usableWidth)
	if fullMode {
		width = float64(usableWidth)
	}
	this.Scene().SetSceneRect2(-H_MARGIN, 0, width + 2 * H_MARGIN, this.Height())

	yCenter := this.Height() / 2
	var xCenter float64
	if fullMode {
		xCenter = float64(this.Scene().Width() - 2 * H_MARGIN) / 2
	} else {
		xCenter = float64(this.LastVisibleIndex + 1) * this.ItemWidth + H_MARGIN - float64(this.Width()) / 2
	}

	this.CenterOn2(xCenter, yCenter)

	this.drawUI()

	this.Scene().Update(this.Scene().SceneRect())
}

func (this *XDecorator) Layout() {
	if this.Model.Count() == 0 {
		return
	}

	if this.ItemWidth <= 0 {
		this.ItemWidth = BEST_ITEM_WIDTH
	}

	// 计算ItemWidth
	width := float64(this.Width()) - 2 * H_MARGIN
	n := int(width / this.ItemWidth)
	if n < VISIBLE_KLINES_MIN {
		n = VISIBLE_KLINES_MIN
	}

	// 计算FirstVisibleIndex
	this.FirstVisibleIndex = this.LastVisibleIndex - n + 1
	if this.FirstVisibleIndex < 0 {
		this.FirstVisibleIndex = 0
	}

	// 设置X transformer Scale
	this.Transformer.SetScale(1 / this.ItemWidth)

	this.UpdateUI()
}

func (this *XDecorator) SetVisibleRange(lastVisibleIndex int, visibleCount int) {
	if this.Model.Count() == 0 {
		return
	}
	if lastVisibleIndex < 0 || lastVisibleIndex >= this.Model.Count() {
		return
	}

	if visibleCount <= 0 {
		return
	}

	firstVisibleIndex := int(mathutil.MaxInt32(0, int32(lastVisibleIndex - visibleCount + 1)))
	visibleCount = lastVisibleIndex - firstVisibleIndex + 1

	this.LastVisibleIndex = lastVisibleIndex
	usableWidth := float64(this.Width()) - 2 * H_MARGIN
	this.ItemWidth = usableWidth / float64(visibleCount)

	this.Layout()
}

func (this *XDecorator) Clear() {
	for _, item := range this.Items {
		this.Scene().RemoveItem(item)
	}
	this.Items = nil

	this.ValueGraph.Clear()
}

func (this *XDecorator) WheelEvent(event *gui.QWheelEvent) {
}

func (this *XDecorator) ShowValue(index int) {
	this.ValueGraph.Clear()
	x := this.Transformer.To(value)

	this.ValueGraph.Update(0, y - VALUE_GRAPH_HEIGHT / 2, s)
}

func (this *XDecorator) HideValue() {
	this.ValueGraph.Clear()
}
