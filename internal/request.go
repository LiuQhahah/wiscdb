package internal

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

func (req *request) DecrRef() {
	nRef := req.ref.Add(-1)
	if nRef > 0 {
		return
	}
	req.Entries = nil
	requestPool.Put(req)
}

func (req *request) Wait() error {
	return nil
}

type requests []*request

func (reqs requests) DecrRef() {

}

func (reqs requests) IncrRef() {

}
