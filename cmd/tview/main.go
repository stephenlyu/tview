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
	"path/filepath"
	"github.com/z-ray/log"
	"runtime"
	"bytes"
	"github.com/stephenlyu/goformula/formulalibrary"
)

const DATA_DIR = "data"
const FORMULA_DIR = "formulas"
const DEBUG = true

var library = formulalibrary.GlobalLibrary

func initFormulaLibrary() {
	formulaDir := FindDir(FORMULA_DIR)
	library.SetDebug(DEBUG)
	library.LoadEasyLangFormulas(formulaDir)
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
	container := graphview.CreateGraphViewContainer(w.Widget)
	w.Push(container)

	rect := app.Desktop().AvailableGeometry(-1)
	w.Widget.Resize2(int(float64(rect.Width()) * 0.8), int(float64(rect.Height()) * 0.8))
	w.Widget.Move2(int(float64(rect.Width()) * 0.1), int(float64(rect.Height()) * 0.1))
	w.Widget.ShowMaximized()

	//container.AddGraphFormula(0, "MA", []float64{5, 10, 20, 60})
	container.AddGraphFormula(0, "EMA513", []float64{5, 13})
	container.AddGraphFormula(1, "MACD", []float64{12, 26, 9})
	container.AddGraphFormula(2, "VOL", []float64{5, 10})
	container.AddGraphFormula(3, "DDGS", []float64{0.02})

	//container.ShowGraph(2)

	timer := core.NewQTimer(w.Widget)
	timer.SetSingleShot(true)
	timer.ConnectTimeout(func () {
		ds := tdxdatasource.NewDataSource(FindDir(DATA_DIR), true)
		security, _ := entity.ParseSecurity("000001.SZ")
		_, period := period.PeriodFromString("D1")
		err, data := ds.GetForwardAdjustedData(security, period)
		if err != nil {
			logrus.Fatalf("加载数据失败，错误：%s", err.Error())
		}
		container.SetData(data, period)
	})
	timer.Start(100)

	os.Exit(app.Exec())

	writer.Close()
}
