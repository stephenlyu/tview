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
	"github.com/stephenlyu/goformula/stockfunc/formula"
	"github.com/stephenlyu/tview/graphs/klinegraph"
	"github.com/stephenlyu/tview/graphs/ploylinegraph"
	"fmt"
	"github.com/stephenlyu/tview/graphs/drawlinegraph"
	"github.com/stephenlyu/tview/graphs/sticklinegraph"
)

type FormulaGraph struct {
	Model        model.Model
	KLineModel   model.Model
	Scene        *widgets.QGraphicsScene
	xTransformer transform.ScaleTransformer

	normalGraphs []graphs.Graph					// Normal graphs is for output variables

	actionGraphs []graphs.Graph					// Action graphs is for all draw ations
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

func (this *FormulaGraph) initNormalGraph() {
	this.normalGraphs = make([]graphs.Graph, this.Model.VarCount())
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
		case constants.GraphTypeNone:
		// DrawLine & PloyLine case
		}
		this.normalGraphs[i] = graph
	}
}

func (this *FormulaGraph) initActionGraph() {
	this.actionGraphs = make([]graphs.Graph, this.Model.DrawActionCount())
	j := 0
	var color *gui.QColor

	for i := 0; i < this.Model.DrawActionCount(); i++ {
		drawAction := this.Model.DrawAction(i)

		var graph graphs.Graph
		fColor := drawAction.GetColor()
		if fColor != nil {
			color = gui.NewQColor3(fColor.Red, fColor.Green, fColor.Blue, 255)
		} else if drawAction.GetVarIndex() != - 1 {
			color = graphs.COLORS[drawAction.GetVarIndex() % len(graphs.COLORS)]
		} else {
			color = graphs.COLORS[j % len(graphs.COLORS)]
			j++
		}

		if !drawAction.IsNoDraw() {
			switch drawAction.GetActionType() {
			case formula.FORMULA_DRAW_ACTION_DRAWKLINE:
				graph = klinegraph.NewDrawKLineGraph(this.Model, drawAction.(formula.DrawKLine), this.Scene, this.xTransformer)
			case formula.FORMULA_DRAW_ACTION_DRAWTEXT:
			case formula.FORMULA_DRAW_ACTION_DRAWICON:
			case formula.FORMULA_DRAW_ACTION_DRAWLINE:
				graph = drawlinegraph.NewDrawLineGraph(this.Model, drawAction.(formula.DrawLine), color, this.Scene, this.xTransformer)
			case formula.FORMULA_DRAW_ACTION_STICKLINE:
				graph = sticklinegraph.NewStickLineGraph(this.Model, drawAction.(formula.StickLine), color, this.Scene, this.xTransformer)
			case formula.FORMULA_DRAW_ACTION_PLOYLINE:
				graph = ploylinegraph.NewPloyLineGraph(this.Model, drawAction.(formula.PloyLine), color, this.Scene, this.xTransformer)
			}
		}
		this.actionGraphs[i] = graph
	}
}

func (this *FormulaGraph) init() {
	this.initNormalGraph()
	this.initActionGraph()
}

func (this *FormulaGraph) GetValueRange(startIndex int, endIndex int) (float64, float64) {
	high := -math.MaxFloat64
	low := math.MaxFloat64

	// Normal Graph
	for i := 0; i < len(this.normalGraphs); i++ {
		if this.normalGraphs[i] == nil {
			continue
		}
		if this.Model.NoDraw(i) {
			continue
		}
		low1, high1 := this.normalGraphs[i].GetValueRange(startIndex, endIndex)
		if low1 < low {
			low = low1
		}
		if high1 > high {
			high = high1
		}
	}

	// Action Graph
	for i := 0; i < len(this.actionGraphs); i++ {
		if this.actionGraphs[i] == nil {
			continue
		}
		action := this.Model.DrawAction(i)
		if action.IsNoDraw() {
			continue
		}
		low1, high1 := this.actionGraphs[i].GetValueRange(startIndex, endIndex)
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
	// Normal Graphs
	for i, graph := range this.normalGraphs {
		if graph == nil {
			continue
		}
		if this.Model.NoDraw(i) {
			continue
		}
		graph.Update(startIndex, endIndex)
	}

	// Action Graphs
	for i, graph := range this.actionGraphs {
		if graph == nil {
			continue
		}
		if this.Model.DrawAction(i).IsNoDraw() {
			continue
		}
		graph.Update(startIndex, endIndex)
	}
}

// 清除所有的K线
func (this *FormulaGraph) Clear() {
	// Normal Graphs
	for i, graph := range this.normalGraphs {
		if graph == nil {
			continue
		}
		if this.Model.NoDraw(i) {
			continue
		}
		graph.Clear()
	}

	// Action Graphs
	for i, graph := range this.actionGraphs {
		if graph == nil {
			continue
		}
		if this.Model.DrawAction(i).IsNoDraw() {
			continue
		}
		graph.Clear()
	}
}

func (this *FormulaGraph) ShowSubInfo(valueIndex int, index int, display graphs.InfoDisplay) {
	if this.Model.NoText(valueIndex) {
		return
	}
	if index < 0 || index >= this.Model.Count() {
		return
	}

	j := 0
	var color *gui.QColor
	fColor := this.Model.Color(valueIndex)
	if fColor != nil {
		color = gui.NewQColor3(fColor.Red, fColor.Green, fColor.Blue, 255)
	} else {
		color = graphs.COLORS[j % len(graphs.COLORS)]
		j++
	}

	name := this.Model.GetNames()[valueIndex]
	v := this.Model.GetRaw(index)[valueIndex]
	display.Add(fmt.Sprintf("%s: %s", name, graphs.FormatValue(v, 2)), color)
}

func (this *FormulaGraph) ShowInfo(index int, display graphs.InfoDisplay) {
	formulaModel := this.Model.(*model.FormulaModel)
	display.Add(formulaModel.Formula.Name(), gui.NewQColor3(255, 255, 255, 255))
	for i, graph := range this.normalGraphs {
		if graph == nil {
			this.ShowSubInfo(i, index, display)
			continue
		}
		graph.ShowInfo(index, display)
	}
}
