package model

import (
	"github.com/stephenlyu/tds/util"
	"github.com/stephenlyu/tview/constants"
	"math"
	"sort"
)

type SeparatorModelListener interface {
	OnModelChanged()
}

type SeparatorModel struct {
	values []float64
	listeners []SeparatorModelListener
}

func NewSeparatorModel() *SeparatorModel {
	return &SeparatorModel{}
}

func (this *SeparatorModel) Count() int {
	return len(this.values)
}

func (this *SeparatorModel) Get(index int) float64 {
	util.Assert(index >= 0 && index < len(this.values), "")
	return this.values[index]
}

func (this *SeparatorModel) Update(min float64, max float64, viewHeight float64) {
	if min == max {
		this.values = []float64{min}
		return
	}

	n := int(viewHeight / constants.SEPARATOR_GAP_MIN)
	if n > constants.SEPARATOR_MAX {
		n = constants.SEPARATOR_MAX
	}

	this.values = calculateValues(min, max, n)
}

func (this *SeparatorModel) AddListener(listener SeparatorModelListener) {
	for _, l := range this.listeners {
		if l == listener {
			return
		}
	}
	this.listeners = append(this.listeners, listener)
}

func (this *SeparatorModel) RemoveListener(listener SeparatorModelListener) {
	for i, l := range this.listeners {
		if l == listener {
			this.listeners = append(this.listeners[:i], this.listeners[i+1:]...)
			return
		}
	}
}

func (this *SeparatorModel) NotifyDataChanged() {
	for _, listener := range this.listeners {
		listener.OnModelChanged()
	}
}

func calculateValues(min float64, max float64, n int) []float64 {
	diff := max - min
	per := diff / float64(n)
	level := int(math.Log10(per))
	unit := math.Pow(10, float64(level))
	per = float64(int(per / unit)) * unit
	if per < 0 {
		per -= unit
	} else {
		per += unit
	}

	var values []float64
	if min > 0 || max < 0 {
		v := float64(int((min + per - 1) / per)) * per
		for ; v < max; v += per {
			values = append(values, v)
		}
	} else {
		values = append(values, 0)
		for v := 0.; v > min; v-=per {
			values = append(values, v)
		}
		for v := 0.; v < max; v += per {
			values = append(values, v)
		}
	}
	sort.SliceStable(values, func (i, j int) bool {
		return values[i] < values[j]
	})
	return values
}