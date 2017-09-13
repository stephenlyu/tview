package main

import (
	"github.com/therecipe/qt/widgets"
	"os"
	"github.com/stephenlyu/tview/views/mainwindow"
	"github.com/stephenlyu/tview/views/graphview"
	"github.com/therecipe/qt/core"
	"github.com/stephenlyu/tds/datasource/tdx"
	"github.com/stephenlyu/tds/entity"
	"github.com/stephenlyu/tds/period"
	"github.com/Sirupsen/logrus"
)

func main() {
	app := widgets.NewQApplication(len(os.Args), os.Args)
	w := mainwindow.GetMainWindow(nil)
	graphView := graphview.CreateGraphView(true, w.Widget)
	w.Push(graphView)
	w.Widget.Show()

	timer := core.NewQTimer(w.Widget)
	timer.SetSingleShot(true)
	timer.ConnectTimeout(func () {
		ds := tdxdatasource.NewDataSource("../../../tds/datasource/tdx/data", true)
		security, _ := entity.ParseSecurity("000001.SZ")
		_, period := period.PeriodFromString("D1")
		err, data := ds.GetData(security, period)
		if err != nil {
			logrus.Fatalf("加载数据失败，错误：%s", err.Error())
		}
		graphView.SetData(data)
	})
	timer.Start(100)

	os.Exit(app.Exec())
}
