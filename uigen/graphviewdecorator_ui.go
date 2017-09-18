// WARNING! All changes made in this file will be lost!
package uigen

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type UIGraphViewDecorator struct {
	VerticalLayout *widgets.QVBoxLayout
	Widget *widgets.QWidget
	InfoLayout *widgets.QHBoxLayout
	Widget2 *widgets.QWidget
	HorizontalLayout *widgets.QHBoxLayout
	Widget3 *widgets.QWidget
	ContentLayout *widgets.QVBoxLayout
	Widget4 *widgets.QWidget
	YAxisLayout *widgets.QVBoxLayout
}

func (this *UIGraphViewDecorator) SetupUI(GraphViewDecorator *widgets.QWidget) {
	GraphViewDecorator.SetObjectName("GraphViewDecorator")
	GraphViewDecorator.SetGeometry(core.NewQRect4(0, 0, 993, 580))
	this.VerticalLayout = widgets.NewQVBoxLayout2(GraphViewDecorator)
	this.VerticalLayout.SetObjectName("verticalLayout")
	this.VerticalLayout.SetContentsMargins(0, 0, 0, 0)
	this.VerticalLayout.SetSpacing(0)
	this.Widget = widgets.NewQWidget(GraphViewDecorator, core.Qt__Widget)
	this.Widget.SetObjectName("Widget")
	this.Widget.SetMinimumSize(core.NewQSize2(0, 30))
	this.InfoLayout = widgets.NewQHBoxLayout2(this.Widget)
	this.InfoLayout.SetObjectName("infoLayout")
	this.InfoLayout.SetContentsMargins(0, 0, 0, 0)
	this.InfoLayout.SetSpacing(0)
	this.VerticalLayout.AddWidget(this.Widget, 0, 0)
	this.Widget2 = widgets.NewQWidget(GraphViewDecorator, core.Qt__Widget)
	this.Widget2.SetObjectName("Widget2")
	this.HorizontalLayout = widgets.NewQHBoxLayout2(this.Widget2)
	this.HorizontalLayout.SetObjectName("horizontalLayout")
	this.HorizontalLayout.SetContentsMargins(0, 0, 0, 0)
	this.HorizontalLayout.SetSpacing(0)
	this.Widget3 = widgets.NewQWidget(this.Widget2, core.Qt__Widget)
	this.Widget3.SetObjectName("Widget3")
	this.ContentLayout = widgets.NewQVBoxLayout2(this.Widget3)
	this.ContentLayout.SetObjectName("contentLayout")
	this.ContentLayout.SetContentsMargins(0, 0, 0, 0)
	this.ContentLayout.SetSpacing(0)
	this.HorizontalLayout.AddWidget(this.Widget3, 0, 0)
	this.Widget4 = widgets.NewQWidget(this.Widget2, core.Qt__Widget)
	this.Widget4.SetObjectName("Widget4")
	this.Widget4.SetMinimumSize(core.NewQSize2(60, 0))
	this.YAxisLayout = widgets.NewQVBoxLayout2(this.Widget4)
	this.YAxisLayout.SetObjectName("yAxisLayout")
	this.YAxisLayout.SetContentsMargins(0, 0, 0, 0)
	this.YAxisLayout.SetSpacing(0)
	this.HorizontalLayout.AddWidget(this.Widget4, 0, 0)
	this.HorizontalLayout.SetStretch(0, 1)
	this.VerticalLayout.AddWidget(this.Widget2, 0, 0)
	this.VerticalLayout.SetStretch(1, 1)


    this.RetranslateUi(GraphViewDecorator)

}

func (this *UIGraphViewDecorator) RetranslateUi(GraphViewDecorator *widgets.QWidget) {
    _translate := core.QCoreApplication_Translate
	GraphViewDecorator.SetWindowTitle(_translate("GraphViewDecorator", "Form", "", -1))
}
