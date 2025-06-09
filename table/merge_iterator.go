package table

import (
	"sort"
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
	cur     *Iterator
	iters   []*Iterator
	tables  []*Table
	options int
}

func (c *concatIterator) Next() {

	c.cur.Next()
	if c.cur.Valid() {
		return
	}
	for {
		if c.options&REVERSED == 0 {
			c.setIdx(c.idx + 1)
		} else {
			c.setIdx(c.idx - 1)
		}
		if c.cur == nil {
			return
		}
		c.cur.ReWind()
		if c.cur.Valid() {
			break
		}
	}
}

func (c *concatIterator) setIdx(idx int) {
	c.idx = idx
	if idx < 0 || idx >= len(c.iters) {
		c.cur = nil
		return
	}
	if c.iters[idx] == nil {
		c.iters[idx] = c.tables[idx].NewIterator(c.options)
	}
	c.cur = c.iters[c.idx]
}

func (c *concatIterator) ReWind() {

	if len(c.iters) == 0 {
		return
	}
	if c.options&REVERSED == 0 {
		c.setIdx(0)
	} else {
		c.setIdx(len(c.iters) - 1)
	}
	c.cur.ReWind()
}

func (c *concatIterator) Seek(key []byte) {
	var idx int
	if c.options&REVERSED == 0 {
		idx = sort.Search(len(c.tables), func(i int) bool {
			return y.CompareKeys(c.tables[i].Biggest(), key) >= 0
		})
	} else {
		n := len(c.tables)
		idx = n - 1 - sort.Search(n, func(i int) bool {
			return y.CompareKeys(c.tables[n-i-1].Smallest(), key) <= 0
		})
	}
	if idx >= len(c.tables) || idx < 0 {
		c.setIdx(-1)
		return
	}
	c.setIdx(idx)
	c.cur.Seek(key)
}

func (c concatIterator) Value() y.ValueStruct {
	return c.cur.Value()
}

func (c *concatIterator) Valid() bool {
	return c.cur != nil && c.cur.Valid()
}

func (c concatIterator) Close() error {
	for _, t := range c.tables {
		if err := t.DecrRef(); err != nil {
			return err
		}
	}

	for _, it := range c.iters {
		if it == nil {
			continue
		}
		if err := it.Close(); err != nil {
			return y.Wrap(err, "ConcatIterator")
		}
	}
	return nil
}

func (c *concatIterator) Key() []byte {
	return c.cur.Key()
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
