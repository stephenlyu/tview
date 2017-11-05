package main

import (
	"github.com/therecipe/qt/widgets"
	"os"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/core"
)

func main() {
	app := widgets.NewQApplication(len(os.Args), os.Args)
	w := widgets.NewQGraphicsView(nil)
	scene := widgets.NewQGraphicsScene(w)
	w.SetScene(scene)
	w.Scale(1, -1)

	w.SetGeometry2(0, 0, 600, 400)
	scene.SetSceneRect2(-300, -200, 600, 400)

	pen := gui.NewQPen3(gui.NewQColor3(255, 0, 0, 255))
	scene.AddLine2(-300, 0, 300, 0, pen)
	scene.AddLine2(0, -200, 0, 200, pen)

	trans := gui.QTransform_FromScale(1.0, -1.0)

	ti := scene.AddText("Hello", w.Font())

	ti.SetTransform(trans, false)
	scene.AddRect(ti.BoundingRect(), pen, gui.NewQBrush2(core.Qt__NoBrush))

	ti.SetPos2(0, ti.BoundingRect().Height())

	w.Show()

	os.Exit(app.Exec())
}
