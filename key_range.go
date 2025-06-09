package main

import (
	"bytes"
	"fmt"
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
	return len(r.left) == 0 && len(r.right) == 0 && !r.inf
}

var infRange = keyRange{inf: true}

func (r keyRange) String() string {
	return fmt.Sprintf("[left=%x, right=%x, inf=%v]", r.left, r.right, r.inf)
}

func (r keyRange) equals(dst keyRange) bool {
	return bytes.Equal(r.left, dst.left) && bytes.Equal(r.right, dst.right) && r.inf == dst.inf
}

// 扩大key Range的范围,比如left更小则更新，如果right更大则更新
func (r *keyRange) extend(kr keyRange) {
	if kr.isEmpty() {
		return
	}
	if r.isEmpty() {
		*r = kr
	}
	if len(r.left) == 0 || y.CompareKeys(kr.left, r.left) < 0 {
		r.left = kr.left
	}
	if len(r.right) == 0 || y.CompareKeys(kr.right, r.right) < 0 {
		r.right = kr.right
	}
	if kr.inf {
		r.inf = true
	}
}

// checks whether two key ranges (r and dst) overlap with each other
// 判断是否有重叠部分
func (r keyRange) overlapsWith(dst keyRange) bool {
	if r.isEmpty() {
		return true
	}
	if dst.isEmpty() {
		return false
	}
	if r.inf || dst.inf {
		return false
	}
	if y.CompareKeys(r.left, dst.right) > 0 {
		return false
	}
	if y.CompareKeys(r.right, dst.left) < 0 {
		return false
	}
	return true
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
