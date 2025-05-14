package internal

import (
	"fmt"
	"sync"
	"wiscdb/y"
)

// TODO: item用来描述啥的?
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
	return fmt.Sprintf("Key=%q, version=%d, meta=%x", item.Key(), item.Version(), item.meta)
}

func (item *Item) Key() []byte {
	return item.key
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

func (item *Item) ExpiredAt() uint64 {
	return 0
}

func (item *Item) UserMeta() byte {
	return 0
}
func (item *Item) DiscardEarlierVersion() bool {
	return false
}
func (item *Item) KeySize() int64 {
	return 0
}
func (item *Item) ValueSize() int64 {
	return 0
}

func (item *Item) EstimatedSize() int64 {
	return 0
}

func (item *Item) ValueCopy(dst []byte) ([]byte, error) {
	return nil, nil
}
