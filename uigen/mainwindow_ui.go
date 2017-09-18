// WARNING! All changes made in this file will be lost!
package uigen

import (
	"github.com/therecipe/qt/widgets"
	"github.com/therecipe/qt/core"
)

type UIMainwindowMainWindow struct {
	Centralwidget *widgets.QWidget
	VerticalLayout *widgets.QVBoxLayout
	Menubar *widgets.QMenuBar
	Statusbar *widgets.QStatusBar
}

func (this *UIMainwindowMainWindow) SetupUI(MainWindow *widgets.QMainWindow) {
	MainWindow.SetObjectName("MainWindow")
	MainWindow.SetGeometry(core.NewQRect4(0, 0, 800, 600))
	this.Centralwidget = widgets.NewQWidget(MainWindow, core.Qt__Widget)
	this.Centralwidget.SetObjectName("Centralwidget")
	this.VerticalLayout = widgets.NewQVBoxLayout2(this.Centralwidget)
	this.VerticalLayout.SetObjectName("verticalLayout")
	this.VerticalLayout.SetContentsMargins(0, 0, 0, 0)
	this.VerticalLayout.SetSpacing(0)
	MainWindow.SetCentralWidget(this.Centralwidget)
	this.Menubar = widgets.NewQMenuBar(MainWindow)
	this.Menubar.SetObjectName("Menubar")
	this.Menubar.SetGeometry(core.NewQRect4(0, 0, 800, 22))
	MainWindow.SetMenuBar(this.Menubar)
	this.Statusbar = widgets.NewQStatusBar(MainWindow)
	this.Statusbar.SetObjectName("Statusbar")
	MainWindow.SetStatusBar(this.Statusbar)


    this.RetranslateUi(MainWindow)

}

func (this *UIMainwindowMainWindow) RetranslateUi(MainWindow *widgets.QMainWindow) {
    _translate := core.QCoreApplication_Translate
	MainWindow.SetWindowTitle(_translate("MainWindow", "MainWindow", "", -1))
}
