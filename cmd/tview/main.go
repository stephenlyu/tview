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
	"github.com/stephenlyu/tview/model"
	"path/filepath"
)

const ROOT = "/Users/stephenlv/go/src/github.com/stephenlyu/tview/cmd/tview"

func initFormulaLibrary() {
	model.GlobalLibrary.Register("MA", model.NewEasyLangFormulaCreatorFactory(filepath.Join(ROOT, "formulas/MA.d")))
	model.GlobalLibrary.Register("MACD", model.NewEasyLangFormulaCreatorFactory(filepath.Join(ROOT, "formulas/MACD.d")))
	model.GlobalLibrary.Register("VOL", model.NewEasyLangFormulaCreatorFactory(filepath.Join(ROOT, "formulas/VOL.d")))
}

func main() {
	initFormulaLibrary()

	app := widgets.NewQApplication(len(os.Args), os.Args)
	w := mainwindow.GetMainWindow(nil)
	//graphView := graphview.CreateGraphView(true, w.Widget)
	container := graphview.CreateGraphViewContainer(w.Widget)
	w.Push(container)
	w.Widget.Show()

	container.AddGraphFormula(0, "MA", []float64{5, 10, 20, 60})
	container.AddGraphFormula(1, "MACD", []float64{12, 26, 9})
	container.AddGraphFormula(2, "VOL", []float64{5, 10})

	timer := core.NewQTimer(w.Widget)
	timer.SetSingleShot(true)
	timer.ConnectTimeout(func () {
		ds := tdxdatasource.NewDataSource(filepath.Join(ROOT, "data"), true)
		security, _ := entity.ParseSecurity("000001.SZ")
		_, period := period.PeriodFromString("D1")
		err, data := ds.GetData(security, period)
		if err != nil {
			logrus.Fatalf("加载数据失败，错误：%s", err.Error())
		}
		container.SetData(data, period)
	})
	timer.Start(100)

	os.Exit(app.Exec())
}
