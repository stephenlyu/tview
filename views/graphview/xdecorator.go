package graphview

import (
	"github.com/therecipe/qt/widgets"
	"github.com/stephenlyu/tview/transform"
	"github.com/stephenlyu/tview/model"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/core"
	"github.com/stephenlyu/tview/constants"
	"github.com/stephenlyu/tview/graphs/valuegraph"
	"github.com/cznic/mathutil"
	"github.com/stephenlyu/tds/period"
	"github.com/stephenlyu/tds/date"
	"fmt"
	"strconv"
	"github.com/stephenlyu/tds/entity"
)

const X_TICK_WIDTH_MIN = 100

var weekMap = map[int]string {
	1: "一",
	2: "二",
	3: "三",
	4: "四",
	5: "五",
	6: "六",
	0: "日",
}

type Ticker struct {
	Index int
	Text string
}

//go:generate qtmoc
type XDecorator struct {
	widgets.QGraphicsView

	Pen *gui.QPen

	Transformer transform.ScaleTransformer
	Data *model.Data
	Period period.Period

	Items []widgets.QGraphicsItem_ITF
	ValueGraph *valuegraph.ValueGraph

	FirstVisibleIndex, LastVisibleIndex int
	ItemWidth float64									// 每个数据占用的屏幕宽度
}

func CreateXDecorator(parent widgets.QWidget_ITF) *XDecorator {
	ret := NewXDecorator(parent)
	ret.Transformer = transform.NewLogicTransformer(1)
	ret.init()
	return ret
}

func (this *XDecorator) init() {
	this.SetRenderHint(gui.QPainter__TextAntialiasing, true)
	this.SetCacheMode(widgets.QGraphicsView__CacheBackground)
	this.SetHorizontalScrollBarPolicy(core.Qt__ScrollBarAlwaysOff)
	this.SetVerticalScrollBarPolicy(core.Qt__ScrollBarAlwaysOff)
	this.SetDragMode(widgets.QGraphicsView__NoDrag)
	this.SetAlignment(core.Qt__AlignLeft)
	this.Scale(1, -1)

	// 设置scene
	scene := widgets.NewQGraphicsScene(this)
	scene.SetBackgroundBrush(gui.NewQBrush4(core.Qt__black, core.Qt__SolidPattern))
	this.SetScene(scene)

	this.Pen = gui.NewQPen3(constants.DECORATOR_TEXT_COLOR)
	this.ValueGraph = valuegraph.NewValueGraph(this.Scene(), float64(this.Width()), VALUE_GRAPH_HEIGHT)

	this.ConnectWheelEvent(this.WheelEvent)
	this.ConnectResizeEvent(this.ResizeEvent)
}

func (this *XDecorator) SetData(data []entity.Record, p period.Period) {
	this.Data = model.NewData(data)
	this.Period = p
}

func (this *XDecorator) getMonthTickers() []Ticker {
	unitWidth := this.Transformer.To(1)

	var result []Ticker

	var prevIndex int = -1
	var prevMonth int = -1

	for i := this.FirstVisibleIndex; i < this.LastVisibleIndex; i++ {
		date := this.Data.GetDate(i)
		if prevIndex != -1 {
			if float64(i - prevIndex) * unitWidth < X_TICK_WIDTH_MIN {
				continue
			}
		}
		year, _ := strconv.Atoi(date[:4])
		month, _ := strconv.Atoi(date[4:6])

		if i == this.FirstVisibleIndex {
			result = append(result, Ticker{Index: i, Text: fmt.Sprintf("%d年", year)})
			prevMonth = month
			prevIndex = i
		} else if month != prevMonth {
			result = append(result, Ticker{Index: i, Text: fmt.Sprintf("%d", month)})
			prevMonth = month
			prevIndex = i
		}
	}
	return result
}

func (this *XDecorator) getMinuteTickers() []Ticker {
	unitWidth := this.Transformer.To(1)

	var result []Ticker
	var prevIndex int = -1

	for i := this.FirstVisibleIndex; i < this.LastVisibleIndex; i++ {
		if prevIndex != -1 {
			if float64(i - prevIndex) * unitWidth < X_TICK_WIDTH_MIN {
				continue
			}
		}
		date := this.Data.GetDate(i)

		if i == this.FirstVisibleIndex {
			month, _ := strconv.Atoi(date[4:6])
			day, _ := strconv.Atoi(date[6:8])
			result = append(result, Ticker{Index: i, Text: fmt.Sprintf("%02d月%02d日", month, day)})
		} else {
			result = append(result, Ticker{Index: i, Text: fmt.Sprintf("%s", date[9:14])})
		}
	}
	return result
}

func (this *XDecorator) getTickers() []Ticker {
	if this.Period.Unit() == period.PERIOD_UNIT_MINUTE {
		return this.getMinuteTickers()
	}
	return this.getMonthTickers()
}

func (this *XDecorator) drawUI() {
	tickers := this.getTickers()

	trans := gui.QTransform_FromScale(1.0, -1.0)
	for _, ticker := range tickers {
		x := this.Transformer.To(float64(ticker.Index))
		tick := this.Scene().AddLine2(x, 0, x, float64(this.Height()), this.Pen)
		this.Items = append(this.Items, tick)

		text := ticker.Text
		ti := this.Scene().AddText(text, this.Font())
		ti.SetDefaultTextColor(constants.DECORATOR_TEXT_COLOR)
		ti.AdjustSize()
		ti.SetTransform(trans, true)
		r := ti.BoundingRect()
		ti.SetPos2(x + 2, (float64(this.Height()) + r.Height()) / 2)
		r = ti.BoundingRect()
		this.Items = append(this.Items, ti)
	}
}

func (this *XDecorator) UpdateUI() {
	this.Clear()

	width := float64(this.Data.Count()) * this.ItemWidth
	usableWidth := this.Width() - 2 * H_MARGIN
	fullMode := width < float64(usableWidth)
	if fullMode {
		width = float64(usableWidth)
	}
	this.Scene().SetSceneRect2(-H_MARGIN, 0, width + 2 * H_MARGIN, float64(this.Height()))

	yCenter := float64(this.Height()) / 2
	var xCenter float64
	if fullMode {
		xCenter = float64(this.Scene().Width() - 2 * H_MARGIN) / 2
	} else {
		xCenter = float64(this.LastVisibleIndex + 1) * this.ItemWidth + H_MARGIN - float64(this.Width()) / 2
	}

	this.CenterOn2(xCenter, yCenter)

	this.drawUI()

	this.Scene().Update(this.Scene().SceneRect())
}

func (this *XDecorator) Layout() {
	if this.Data == nil {
		return
	}
	if this.Data.Count() == 0 {
		return
	}

	if this.ItemWidth <= 0 {
		this.ItemWidth = BEST_ITEM_WIDTH
	}

	// 计算ItemWidth
	width := float64(this.Width()) - 2 * H_MARGIN
	n := int(width / this.ItemWidth)
	if n < VISIBLE_KLINES_MIN {
		n = VISIBLE_KLINES_MIN
	}

	// 计算FirstVisibleIndex
	this.FirstVisibleIndex = this.LastVisibleIndex - n + 1
	if this.FirstVisibleIndex < 0 {
		this.FirstVisibleIndex = 0
	}

	// 设置X transformer Scale
	this.Transformer.SetScale(1 / this.ItemWidth)

	this.UpdateUI()
}

func (this *XDecorator) SetVisibleRange(lastVisibleIndex int, visibleCount int) {
	if this.Data.Count() == 0 {
		return
	}
	if lastVisibleIndex < 0 || lastVisibleIndex >= this.Data.Count() {
		return
	}

	if visibleCount <= 0 {
		return
	}

	firstVisibleIndex := int(mathutil.MaxInt32(0, int32(lastVisibleIndex - visibleCount + 1)))
	visibleCount = lastVisibleIndex - firstVisibleIndex + 1

	this.LastVisibleIndex = lastVisibleIndex
	usableWidth := float64(this.Width()) - 2 * H_MARGIN
	this.ItemWidth = usableWidth / float64(visibleCount)

	this.Layout()
}

func (this *XDecorator) TrackPoint(currentIndex int, x float64, y float64) {
	if currentIndex < this.FirstVisibleIndex {
		this.LastVisibleIndex -= (this.FirstVisibleIndex - currentIndex)
		this.Layout()
	} else if currentIndex > this.LastVisibleIndex {
		this.LastVisibleIndex = currentIndex
		this.Layout()
	}

	this.ShowValue(currentIndex)
}

func (this *XDecorator) TrackXY(globalX float64, globalY float64) {
	pt := core.NewQPoint2(int(globalX), int(globalY))
	pt = this.MapFromGlobal(pt)
	ptScene := this.MapToScene(pt)
	currentIndex := int(this.Transformer.From(ptScene.X()))
	if currentIndex >= this.Data.Count() {
		currentIndex = this.Data.Count() - 1
	}

	if currentIndex < 0 {
		currentIndex = 0
	}
	this.ShowValue(currentIndex)
}

func (this *XDecorator) CompleteTrackPoint() {
	this.HideValue()
}

func (this *XDecorator) Clear() {
	for _, item := range this.Items {
		this.Scene().RemoveItem(item)
	}
	this.Items = nil

	this.ValueGraph.Clear()
}

func (this *XDecorator) WheelEvent(event *gui.QWheelEvent) {
}

// Event Handlers
func (this *XDecorator) ResizeEvent(event *gui.QResizeEvent) {
	this.Layout()
}

func (this *XDecorator) ShowValue(index int) {
	this.ValueGraph.Clear()
	x := this.Transformer.To(float64(index))

	var text string
	dt := this.Data.GetDate(index)
	if this.Period.Unit() == period.PERIOD_UNIT_MINUTE {
		text = fmt.Sprintf("%s/%s %s", dt[4:6], dt[6:8], dt[9:14])
	} else {
		ts, _ := date.SecondString2Timestamp(dt)
		week := date.GetDateWeekDay(ts)
		text = fmt.Sprintf("%s/%s/%s/%s", dt[0:4], dt[4:6], dt[6:8], weekMap[week])
	}

	width := this.Transformer.To(float64(this.LastVisibleIndex + 1))
	r := this.ValueGraph.MeasureText(text)
	if x + r.Width() > width {
		x = width - r.Width()
	}

	this.ValueGraph.Update(x, 0, text)
}

func (this *XDecorator) HideValue() {
	this.ValueGraph.Clear()
}
