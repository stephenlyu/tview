package main

import (
	"github.com/therecipe/qt/widgets"
	"os"
	"github.com/stephenlyu/tview/views/mainwindow"
	"github.com/stephenlyu/tview/views/graphview"
)

func main() {
	app := widgets.NewQApplication(len(os.Args), os.Args)
	w := mainwindow.GetMainWindow(nil)
	w.Push(graphview.CreateGraphView(w.Widget))
	w.Widget.Show()

	os.Exit(app.Exec())
}
