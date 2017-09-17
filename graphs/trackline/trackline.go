package trackline

import (
	"github.com/therecipe/qt/widgets"
	"github.com/stephenlyu/tview/constants"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/core"
)


type TrackLine struct {
	Scene *widgets.QGraphicsScene
	pen *gui.QPen

	XLine *widgets.QGraphicsLineItem
	YLine *widgets.QGraphicsLineItem
}

func NewTrackLine(scene *widgets.QGraphicsScene) *TrackLine {
	pen := gui.NewQPen3(constants.TRACK_LINE_COLOR)
	pen.SetWidth(1)
	pen.SetStyle(core.Qt__SolidLine)
	return &TrackLine{
		Scene: scene,
		pen: pen,
	}
}

func (this *TrackLine) UpdateTrackLine(x float64, y float64) {
	this.Clear()

	r := this.Scene.SceneRect()

	this.YLine = this.Scene.AddLine2(x, r.Y(), x, r.Y() + r.Height(), this.pen)

	this.XLine = this.Scene.AddLine2(r.X(), y, r.X() + r.Width(), y, this.pen)
}

func (this *TrackLine) Clear() {
	if this.XLine != nil {
		this.Scene.RemoveItem(this.XLine)
		this.XLine = nil
	}
	if this.YLine != nil {
		this.Scene.RemoveItem(this.YLine)
		this.YLine = nil
	}
}
