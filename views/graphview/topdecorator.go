package graphview

import (
	"github.com/therecipe/qt/widgets"
	"github.com/therecipe/qt/gui"
	"fmt"
	"github.com/therecipe/qt/core"
)

type TopDecorator struct {
	InfoLayout *widgets.QHBoxLayout

	Labels []*widgets.QLabel
	currentIndex int
}

func NewTopDecorator(layout *widgets.QHBoxLayout) *TopDecorator {
	return &TopDecorator{InfoLayout: layout}
}

func (this *TopDecorator) Clear() {
	for _, label := range this.Labels {
		label.SetText("")
		label.Hide()
	}
	this.currentIndex = -1
}

func (this *TopDecorator) ensureLabel() *widgets.QLabel {
	this.currentIndex++
	if this.currentIndex >= len(this.Labels) {
		label := widgets.NewQLabel2("", nil, core.Qt__Widget)
		this.Labels = append(this.Labels, label)
		this.InfoLayout.InsertWidget(this.currentIndex, label, 0, 0)
	}
	return this.Labels[this.currentIndex]
}

func (this *TopDecorator) Add(text string, color *gui.QColor) {
	label := this.ensureLabel()
	label.SetText(fmt.Sprintf("<font color='%s'>%s</font>", color.Name(), text))
	label.Show()
}
