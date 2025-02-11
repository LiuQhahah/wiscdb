package internal

import (
	"sync"
	"wiscdb/y"
)

type Item struct {
	key       []byte
	vptr      []byte
	val       []byte
	version   uint64
	expiresAt uint64
	slice     *y.Slice
	next      *Item
	txn       *Txn
	err       error
	wg        sync.WaitGroup
	status    int
	meta      byte
	userMeta  byte
}

func (item *Item) String() string {
	return ""
}

func (item *Item) Key() []byte {
	return nil
}

func (item *Item) KeyCopy(dst []byte) []byte {
	return nil
}

func (item *Item) Version() uint64 {
	return 0
}

func (item *Item) Value(fn func(val []byte) error) error {
	return nil
}

func (item *Item) yieldItemValue() ([]byte, func(), error) {
	return nil, func() {

	}, nil
}
