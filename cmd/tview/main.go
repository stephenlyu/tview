package main

import (
	"github.com/therecipe/qt/widgets"
	"os"
	"strings"
	"github.com/stephenlyu/tview/views/mainwindow"
	"github.com/stephenlyu/tview/views/graphview"
	"github.com/therecipe/qt/core"
	"github.com/stephenlyu/tds/datasource/tdx"
	"github.com/stephenlyu/tds/entity"
	"github.com/stephenlyu/tds/period"
	"github.com/Sirupsen/logrus"
	"github.com/stephenlyu/tview/model"
	"path/filepath"
	"github.com/z-ray/log"
	"runtime"
	"bytes"
)

const DATA_DIR = "data"
const FORMULA_DIR = "formulas"

func initFormulaLibrary() {
	formulaDir := FindDir(FORMULA_DIR)
	model.GlobalLibrary.Register("MA", model.NewEasyLangFormulaCreatorFactory(filepath.Join(formulaDir, "MA.d")))
	model.GlobalLibrary.Register("MACD", model.NewEasyLangFormulaCreatorFactory(filepath.Join(formulaDir, "MACD.d")))
	model.GlobalLibrary.Register("VOL", model.NewEasyLangFormulaCreatorFactory(filepath.Join(formulaDir, "VOL.d")))
}

func FindDir(dirName string) string {
	_, err := os.Stat(dirName)
	if err == nil {
		return dirName
	}

	if !os.IsNotExist(err) {
		panic(err)
	}

	cwd, _ := os.Getwd()

	cwd, _ = filepath.Abs(cwd)
	separator := string([]byte{os.PathSeparator})
	parts := strings.Split(filepath.Clean(cwd), separator)
	for i := len(parts) - 1; i > 0; i-- {
		filePath := filepath.Join(append(parts[:i], dirName)...)
		if cwd[0] == os.PathSeparator {
			filePath = separator + filePath
		}

		_, err := os.Stat(filePath)
		if err == nil {
			return filePath
		}
	}
	return dirName
}

func PanicTrace(kb int) []byte {
	s := []byte("/src/runtime/panic.go")
	e := []byte("\ngoroutine ")
	line := []byte("\n")
	stack := make([]byte, kb<<10) //4KB
	length := runtime.Stack(stack, true)
	start := bytes.Index(stack, s)
	stack = stack[start:length]
	start = bytes.Index(stack, line) + 1
	stack = stack[start:]
	end := bytes.LastIndex(stack, line)
	if end != -1 {
		stack = stack[:end]
	}
	end = bytes.Index(stack, e)
	if end != -1 {
		stack = stack[:end]
	}
	stack = bytes.TrimRight(stack, "\n")
	return stack
}

func main() {
	writer, _ := os.Create("tview.log")

	log.SetOutput(writer)
	defer func() {
		if err := recover(); err != nil {
			log.Println(string(PanicTrace(2)))
			log.Error(err)
			writer.Close()
			os.Exit(-2)
		}
	}()

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
		ds := tdxdatasource.NewDataSource(FindDir(DATA_DIR), true)
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

	writer.Close()
}
