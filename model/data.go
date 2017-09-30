package model

import (
	"github.com/stephenlyu/tds/entity"
	"github.com/stephenlyu/goformula/stockfunc/function"
)

type Data struct {
	rawData []entity.Record
	transData []function.Record
}

func NewData(data []entity.Record) *Data {
	ret := &Data {
		rawData: data,
	}
	ret.init()
	return ret
}

func (this *Data) init() {
	this.transData = make([]function.Record, len(this.rawData))
	for i := range this.rawData {
		this.transData[i] = &this.rawData[i]
	}
}

func (this *Data) Append(record *entity.Record) {
	this.rawData = append(this.rawData, *record)
	this.transData = append(this.transData, &this.rawData[len(this.rawData) - 1])
}

func (this *Data) Get(index int) function.Record {
	if index < 0 || index >= len(this.transData) {
		panic("index out of range")
	}
	return this.transData[index]
}

func (this *Data) Count() int {
	return len(this.transData)
}

func (this *Data) Records() []function.Record {
	return this.transData
}

func (this *Data) GetDate(index int) string {
	if index < 0 || index >= this.Count() {
		panic("bad index")
	}

	return this.Get(index).GetDate()
}
