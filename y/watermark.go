package y

import (
	"container/heap"
	"context"
	"sync/atomic"

	"github.com/dgraph-io/ristretto/v2/z"
)

// 这种水印机制常见于：
//
// 流处理系统中跟踪处理进度
//
// 数据库事务的有序提交
//
// 分布式系统中的一致性检查点
//
// 确保操作的顺序性和完整性
type WaterMark struct {
	doneUntil atomic.Uint64
	lastIndex atomic.Uint64
	Name      string
	markCh    chan mark
}

func (w *WaterMark) Init(closer *z.Closer) {
	w.markCh = make(chan mark, 100) //创建长度为100的channel
	go w.process(closer)
}

// TODO: WaterMark 水印的作用
func (w *WaterMark) Begin(index uint64) {
	w.lastIndex.Store(index)
	w.markCh <- mark{index: index, done: false}
}

func (w *WaterMark) BeginMany(indices []uint64) {
	w.lastIndex.Store(indices[len(indices)-1])
	w.markCh <- mark{
		index: indices[len(indices)-1],
		done:  false,
	}
}

// 将index写到WaterMark的markCh中
func (w *WaterMark) Done(index uint64) {
	w.markCh <- mark{
		index: index,
		done:  true,
	}
}

func (w *WaterMark) DoneMany(indices []uint64) {
	w.markCh <- mark{
		index:   0,
		indices: indices,
		done:    true,
	}
}

func (w *WaterMark) DoneUntil() uint64 {
	return w.doneUntil.Load()
}

func (w *WaterMark) LastIndex() uint64 {
	return w.lastIndex.Load()
}

func (w *WaterMark) process(closer *z.Closer) {
	// waitGroup 设置成完成
	defer closer.Done()

	var indices uint64Heap
	pending := make(map[uint64]int)
	waiters := make(map[uint64][]chan struct{})

	heap.Init(&indices)

	processOne := func(index uint64, done bool) {
		prev, present := pending[index]
		if !present {
			heap.Push(&indices, index)
		}
		delta := 1
		if done {
			delta = -1
		}
		pending[index] = prev + delta

		doneUntil := w.DoneUntil()
		if doneUntil > index {
			AssertTruef(false, "Name: %s doneUntil: %d. Index: %d", w.Name, doneUntil, index)
		}

		until := doneUntil
		loops := 0
		for len(indices) > 0 {
			min := indices[0]
			if done := pending[min]; done > 0 {
				break
			}
			heap.Pop(&indices)
			delete(pending, min)
			until = min
			loops++
		}
		if until != doneUntil {
			AssertTrue(w.doneUntil.CompareAndSwap(doneUntil, until))
		}

		notifyAndRemove := func(idx uint64, toNotify []chan struct{}) {
			for _, ch := range toNotify {
				close(ch)
			}
			delete(waiters, idx)
		}

		if until-doneUntil <= uint64(len(waiters)) {
			for idx := doneUntil + 1; idx <= until; idx++ {
				if toNotify, ok := waiters[idx]; ok {
					notifyAndRemove(idx, toNotify)
				}
			}
		} else {
			for idx, toNotify := range waiters {
				if idx <= until {
					notifyAndRemove(idx, toNotify)
				}
			}
		}
	}
	for {
		select {
		case <-closer.HasBeenClosed():
			return
		case mark := <-w.markCh:
			if mark.waiter != nil {
				doneUtil := w.doneUntil.Load()
				if doneUtil >= mark.index {
					close(mark.waiter)
				} else {
					ws, ok := waiters[mark.index]
					if !ok {
						waiters[mark.index] = []chan struct{}{
							mark.waiter,
						}
					} else {
						waiters[mark.index] = append(ws, mark.waiter)
					}
				}
			} else {
				if mark.index > 0 || (mark.index == 0 && len(mark.indices) == 0) {
					processOne(mark.index, mark.done)
				}
				for _, index := range mark.indices {
					processOne(index, mark.done)
				}
			}
		}
	}
}

func (w *WaterMark) SetDoneUtil(val uint64) {
	w.doneUntil.Store(val)
}

func (w *WaterMark) WaitForMark(ctx context.Context, index uint64) error {
	if w.DoneUntil() >= index {
		return nil
	}
	waitCh := make(chan struct{})
	w.markCh <- mark{
		index:  index,
		waiter: waitCh,
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-waitCh:
		return nil
	}
}
