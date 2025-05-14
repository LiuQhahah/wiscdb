package main

import (
	"math"
	"wiscdb/table"
	"wiscdb/y"
)

type keyRange struct {
	left  []byte
	right []byte
	inf   bool
	size  int64
}

func (r keyRange) isEmpty() bool {
	return false
}

var infRange = keyRange{inf: true}

func (r keyRange) String() string {
	return ""
}

func (r keyRange) equals(dst keyRange) bool {
	return false
}

func (r *keyRange) extend(kr keyRange) {

}

func (r keyRange) overlapsWith(dst keyRange) bool {
	return false
}

// 返回多张表的key返回，主要是最大值与最小值
func getKeyRange(tables ...*table.Table) keyRange {
	if len(tables) == 0 {
		return keyRange{}
	}
	//初始化最小值和最大值
	smallest := tables[0].Smallest()
	biggest := tables[0].Biggest()
	//遍历N张表
	for i := 1; i < len(tables); i++ {
		if y.CompareKeys(tables[i].Smallest(), smallest) < 0 {
			smallest = tables[i].Smallest()
		}
		if y.CompareKeys(tables[i].Biggest(), biggest) > 0 {
			biggest = tables[i].Biggest()
		}
	}
	//返回一张表中的最小值和最大值
	return keyRange{
		left:  y.KeyWithTs(y.ParseKey(smallest), math.MaxUint64),
		right: y.KeyWithTs(y.ParseKey(biggest), 0),
	}
}
