// WARNING! All changes made in this file will be lost!
package uigen

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type UIXbarForm struct {
	HorizontalLayout *widgets.QHBoxLayout
	Widget *widgets.QWidget
	XDecoratorLayout *widgets.QHBoxLayout
	Widget2 *widgets.QWidget
	PlaceHolderLayout *widgets.QHBoxLayout
	BtnPeriod *widgets.QPushButton
}

func (this *UIXbarForm) SetupUI(Form *widgets.QWidget) {
	Form.SetObjectName("Form")
	Form.SetGeometry(core.NewQRect4(0, 0, 731, 20))
	Form.SetMinimumSize(core.NewQSize2(0, 20))
	Form.SetMaximumSize(core.NewQSize2(16777215, 20))
	Form.SetStyleSheet("border-top: 1px solid red;\nborder-bottom: 1px solid red;")
	this.HorizontalLayout = widgets.NewQHBoxLayout2(Form)
	this.HorizontalLayout.SetObjectName("horizontalLayout")
	this.HorizontalLayout.SetContentsMargins(0, 0, 0, 0)
	this.HorizontalLayout.SetSpacing(0)
	this.Widget = widgets.NewQWidget(Form, core.Qt__Widget)
	this.Widget.SetObjectName("Widget")
	this.Widget.SetMinimumSize(core.NewQSize2(60, 0))
	this.Widget.SetStyleSheet("background-color:black;")
	this.XDecoratorLayout = widgets.NewQHBoxLayout2(this.Widget)
	this.XDecoratorLayout.SetObjectName("XDecoratorLayout")
	this.XDecoratorLayout.SetContentsMargins(0, 0, 0, 0)
	this.XDecoratorLayout.SetSpacing(0)
	this.HorizontalLayout.AddWidget(this.Widget, 0, 0)
	this.Widget2 = widgets.NewQWidget(Form, core.Qt__Widget)
	this.Widget2.SetObjectName("Widget2")
	this.Widget2.SetMinimumSize(core.NewQSize2(60, 0))
	this.Widget2.SetMaximumSize(core.NewQSize2(60, 16777215))
	this.Widget2.SetStyleSheet("background-color:black;")
	this.PlaceHolderLayout = widgets.NewQHBoxLayout2(this.Widget2)
	this.PlaceHolderLayout.SetObjectName("placeHolderLayout")
	this.PlaceHolderLayout.SetContentsMargins(0, 0, 0, 0)
	this.PlaceHolderLayout.SetSpacing(0)
	this.BtnPeriod = widgets.NewQPushButton(this.Widget2)
	this.BtnPeriod.SetObjectName("BtnPeriod")
	this.BtnPeriod.SetMinimumSize(core.NewQSize2(60, 0))
	this.BtnPeriod.SetMaximumSize(core.NewQSize2(60, 20))
	this.BtnPeriod.SetStyleSheet("color:white;")
	this.PlaceHolderLayout.AddWidget(this.BtnPeriod, 0, 0)
	this.HorizontalLayout.AddWidget(this.Widget2, 0, 0)
	this.HorizontalLayout.SetStretch(0, 1)


    this.RetranslateUi(Form)

}

func (this *UIXbarForm) RetranslateUi(Form *widgets.QWidget) {
    _translate := core.QCoreApplication_Translate
	Form.SetWindowTitle(_translate("Form", "Form", "", -1))
	this.BtnPeriod.SetText(_translate("Form", "1分钟", "", -1))
}
