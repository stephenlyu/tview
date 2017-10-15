package formulagraph

import (
	"github.com/stephenlyu/tview/model"
	"github.com/stephenlyu/tds/util"
	"github.com/stephenlyu/tview/transform"
	"github.com/therecipe/qt/widgets"
	"github.com/stephenlyu/tview/graphs"
	"github.com/stephenlyu/tview/constants"
	"github.com/stephenlyu/tview/graphs/linegraph"
	"github.com/stephenlyu/tview/graphs/colorstickgraph"
	"github.com/stephenlyu/tview/graphs/volgraph"
	"github.com/therecipe/qt/gui"
	"math"
	"github.com/stephenlyu/tview/graphs/stickgraph"
	"github.com/stephenlyu/tview/graphs/linestickgraph"
)

type FormulaGraph struct {
	Model model.Model
	KLineModel model.Model
	Scene *widgets.QGraphicsScene
	xTransformer transform.ScaleTransformer

	graphs []graphs.Graph
}

func NewFormulaGraph(model model.Model, klineModel model.Model, scene *widgets.QGraphicsScene, xTransformer transform.ScaleTransformer) *FormulaGraph {
	util.Assert(model != nil, "model != nil")

	this := &FormulaGraph{
		Model: model,
		KLineModel: klineModel,
		Scene: scene,
		xTransformer: xTransformer,
	}
	this.init()
	return this
}

func (this *FormulaGraph) init() {
	this.graphs = make([]graphs.Graph, this.Model.VarCount())
	j := 0
	var color *gui.QColor
	for i := 0; i < this.Model.VarCount(); i++ {
		var graph graphs.Graph
		fColor := this.Model.Color(i)
		if fColor != nil {
			color = gui.NewQColor3(fColor.Red, fColor.Green, fColor.Blue, 255)
		} else {
			color = graphs.COLORS[j % len(graphs.COLORS)]
			j++
		}
		switch this.Model.GraphType(i) {
		case constants.GraphTypeLine:
			graph = linegraph.NewLineGraph(this.Model, i, color, this.Scene, this.xTransformer)
		case constants.GraphTypeColorStick:
			graph = colorstickgraph.NewColorStickGraph(this.Model, i, color, this.Scene, this.xTransformer)
		case constants.GraphTypeStick:
			graph = stickgraph.NewStickGraph(this.Model, i, color, this.Scene, this.xTransformer)
		case constants.GraphTypeLineStick:
			graph = linestickgraph.NewLineStickGraph(this.Model, i, color, this.Scene, this.xTransformer)
		case constants.GraphTypeVolStick:
			graph = volgraph.NewVolStickGraph(this.Model, i, color, this.KLineModel, this.Scene, this.xTransformer)
		}
		this.graphs[i] = graph
	}
}

func (this *FormulaGraph) GetValueRange(startIndex int, endIndex int) (float64, float64) {
	high := -math.MaxFloat64
	low := math.MaxFloat64

	for i := 0; i < len(this.graphs); i++ {
		if this.Model.NoDraw(i) {
			continue
		}
		low1, high1 := this.graphs[i].GetValueRange(startIndex, endIndex)
		if low1 < low {
			low = low1
		}
		if high1 > high {
			high = high1
		}
	}

	return low, high
}

// 更新当前显示的K线
func (this *FormulaGraph) Update(startIndex int, endIndex int) {
	for i, graph := range this.graphs {
		if this.Model.NoDraw(i) {
			continue
		}
		graph.Update(startIndex, endIndex)
	}
}

// 清除所有的K线
func (this *FormulaGraph) Clear() {
	for i, graph := range this.graphs {
		if this.Model.NoDraw(i) {
			continue
		}
		graph.Clear()
	}
}

func (this *FormulaGraph) ShowInfo(index int, display graphs.InfoDisplay) {
	formulaModel := this.Model.(*model.FormulaModel)
	display.Add(formulaModel.Formula.Name(), gui.NewQColor3(255, 255, 255, 255))
	for _, graph := range this.graphs {
		graph.ShowInfo(index, display)
	}
}
