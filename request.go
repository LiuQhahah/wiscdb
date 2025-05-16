package main

import (
	"sync"
	"sync/atomic"
)

type request struct {
	Entries []*Entry
	Ptrs    []valuePointer
	Wg      sync.WaitGroup
	Err     error
	ref     atomic.Int32
}

func (req *request) reset() {
	req.Entries = req.Entries[:0]
	req.Ptrs = req.Ptrs[:0]
	req.Wg = sync.WaitGroup{}
	req.Err = nil
	req.ref.Store(0)
}
func (req *request) IncrRef() {
	req.ref.Add(1)
}

// 引用减1 如果引用为0,表明请求没有被引用,此时将request还回到pool中
func (req *request) DecrRef() {
	nRef := req.ref.Add(-1)
	if nRef > 0 {
		return
	}
	// 如果引用为0,则调用下方函数
	req.Entries = nil
	//pool的Put函数作用: 把req放回到request pool中
	requestPool.Put(req)
}

func (req *request) Wait() error {
	req.Wg.Wait()
	err := req.Err
	req.DecrRef()
	return err
}

type requests []*request

func (reqs requests) DecrRef() {

}

func (reqs requests) IncrRef() {

}
