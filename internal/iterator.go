package internal

import "wiscdb/y"

type Iterator struct {
	iitr   y.Iterator
	txn    *Txn
	readTs uint64
	opt    IteratorOptions
	item   *Item
	data   list
}

func (it *Iterator) Seek(key []byte) {

}

func (it *Iterator) Rewind() {

}

func (it *Iterator) Close() {

}

func (it *Iterator) prefetch() {

}

func hasPrefix(it *Iterator) bool {
	return false
}
func (it *Iterator) ValidForPrefix(prefix []byte) bool {
	return false
}

func (it *Iterator) parseItem() bool {
	return false
}

func (it *Iterator) newItem() *Item {
	return &Item{}
}

func (it *Iterator) fill(item *Item) {

}

func (it *Iterator) Valid() bool {
	return false
}

func (it *Iterator) Next() {

}
