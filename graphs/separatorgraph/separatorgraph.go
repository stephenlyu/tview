package separatorgraph

import (
	"github.com/therecipe/qt/widgets"
	"github.com/therecipe/qt/gui"
	"github.com/stephenlyu/tview/constants"
	"github.com/stephenlyu/tview/model"
	"github.com/therecipe/qt/core"
	"github.com/stephenlyu/tview/transform"
)

type SeparatorGraph struct {
	Scene *widgets.QGraphicsScene
	transformer transform.ScaleTransformer
	pen *gui.QPen

	model *model.SeparatorModel
	items []*widgets.QGraphicsLineItem
}

func NewSeparatorGraph(scene *widgets.QGraphicsScene, transformer transform.ScaleTransformer, model *model.SeparatorModel) *SeparatorGraph {
	pen := gui.NewQPen3(constants.SEPARATOR_LINE_COLOR)
	pen.SetWidth(1)
	pen.SetStyle(core.Qt__DotLine)
	ret := &SeparatorGraph{
		Scene: scene,
		transformer: transformer,
		pen: pen,
		model: model,
	}
	model.AddListener(ret)
	return ret
}

func (this *SeparatorGraph) buildLines() {
	this.Clear()

	r := this.Scene.SceneRect()
	for i := 0; i < this.model.Count(); i++ {
		v := this.transformer.To(this.model.Get(i))
		item := this.Scene.AddLine2(r.X(), v, r.X() + r.Width(), v, this.pen)
		item.SetZValue(-1)
		this.items = append(this.items, item)
	}
}

func (this *SeparatorGraph) Clear() {
	for _, item := range this.items {
		this.Scene.RemoveItem(item)
	}
	this.items = nil
}

func (this *SeparatorGraph) OnModelChanged(yMin float64, yMax float64) {
	this.buildLines()
}
