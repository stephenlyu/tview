package graphview

import (
	"github.com/therecipe/qt/widgets"
	"github.com/therecipe/qt/core"
	"github.com/stephenlyu/tview/uigen"
	"github.com/stephenlyu/tds/period"
	"github.com/stephenlyu/tds/entity"
)

//go:generate qtmoc
type XBar struct {
	widgets.QWidget
	uigen.UIXbarForm

	*XDecorator
}

// Life Cycle Routines

func CreateXBar(parent widgets.QWidget_ITF) *XBar {
	this := NewXBar(parent, core.Qt__Widget)
	this.SetupUI(&this.QWidget)

	this.XDecorator = CreateXDecorator(this)
	this.XDecoratorLayout.AddWidget(this.XDecorator, 0, 0)

	return this
}

func (this *XBar) SetData(data []entity.Record, p period.Period) {
	this.BtnPeriod.SetText(p.DisplayName())
	this.XDecorator.SetData(data, p)
}
