package y

import (
	"context"
	"github.com/dgraph-io/ristretto/v2/z"
	"sync/atomic"
)

type WaterMark struct {
	doneUntil atomic.Uint64
	lastIndex atomic.Uint64
	Name      string
	markCh    chan mark
}

func (w *WaterMark) Init(closer *z.Closer) {

}

func (w *WaterMark) Begin(index uint64) {

}

func (w *WaterMark) BeginMany(indices []uint64) {

}

func (w *WaterMark) Done(index uint64) {

}
func (w *WaterMark) DoneMany(indices []uint64) {

}

func (w *WaterMark) DoneUntil() uint64 {
	return 0
}

func (w *WaterMark) LastIndex() uint64 {
	return 0
}

func (w *WaterMark) process(closer *z.Closer) {

}

func (w *WaterMark) SetDoneUtil(val uint64) {

}
func (w *WaterMark) WaitForMark(ctx context.Context, index uint64) error {
	return nil
}
