package graphs

import "github.com/therecipe/qt/gui"

type InfoDisplay interface {
	Add(text string, color *gui.QColor)
}

type Graph interface {
	GetValueRange(startIndex int, endIndex int) (float64, float64)
	Update(startIndex int, endIndex int)
	Clear()

	ShowInfo(index int, display InfoDisplay)
}
