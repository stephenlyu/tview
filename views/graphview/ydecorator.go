package graphview

import (
	"github.com/therecipe/qt/widgets"
	"github.com/stephenlyu/tview/transform"
	"github.com/stephenlyu/tview/model"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/core"
	"math"
	"fmt"
	"github.com/stephenlyu/tview/constants"
	"github.com/stephenlyu/tview/graphs/valuegraph"
)

const VALUE_GRAPH_HEIGHT = 20

//go:generate qtmoc
type YDecorator struct {
	widgets.QGraphicsView

	Pen *gui.QPen

	Transformer transform.ScaleTransformer
	Model *model.SeparatorModel

	Items []widgets.QGraphicsItem_ITF
	ValueGraph *valuegraph.ValueGraph
}

func CreateYDecorator(parent widgets.QWidget_ITF) *YDecorator {
	ret := NewYDecorator(parent)
	ret.init()
	return ret
}

func (this *YDecorator) init() {
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

func (this *YDecorator) SetTransformer(transformer transform.ScaleTransformer) {
	this.Transformer = transformer
}

func (this *YDecorator) SetModel(model *model.SeparatorModel) {
	if this.Model != nil {
		this.Model.RemoveListener(this)
	}

	this.Model = model
	if model != nil {
		this.Model.AddListener(this)
	}
}

func (this *YDecorator) OnModelChanged(yMin float64, yMax float64) {
	this.UpdateUI(yMin, yMax)
}

func (this *YDecorator) UpdateUI(yMin float64, yMax float64) {
	this.Clear()

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

func (this *YDecorator) Clear() {
	for _, item := range this.Items {
		this.Scene().RemoveItem(item)
	}
	this.Items = nil

	this.ValueGraph.Clear()
}

func (this *YDecorator) WheelEvent(event *gui.QWheelEvent) {
}

func (this *YDecorator) ShowValue(value float64) {
	this.ValueGraph.Clear()
	y := this.Transformer.To(value)
	s := this.formatValue(value)

	this.ValueGraph.Update(0, y - VALUE_GRAPH_HEIGHT / 2, s)
}

func (this *YDecorator) HideValue() {
	this.ValueGraph.Clear()
}

func (this *YDecorator) formatValue(value float64) string {
	var v = value
	if v < 0 {
		v = -v
	}
	level := math.Log10(v)

	switch {
	case level >= 4:
		return fmt.Sprintf("%.0f", value)
	default:
		return fmt.Sprintf("%.02f", value)
	}
}
