package table

import (
	"wiscdb/y"
)

type MergeIterator struct {
	left    node
	right   node
	small   *node
	curKey  []byte
	reverse bool
}

type node struct {
	valid  bool
	key    []byte
	iter   y.Iterator
	merge  *MergeIterator
	concat *concatIterator
}

type concatIterator struct {
	idx     int
	cur     *y.Iterator
	iters   []*Iterator
	tables  []*Table
	options int
}

func NewMergeIterator(iters []y.Iterator, reverse bool) y.Iterator {
	switch len(iters) {
	case 0:
		return nil
	case 1:
		return iters[0]
	case 2:
		return nil
	}
	return nil
}
