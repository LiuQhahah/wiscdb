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

}
func (req *request) IncrRef() {

}

func (req *request) DecrRef() {

}

func (req *request) Wait() error {
	return nil
}

type requests []*request

func (reqs requests) DecrRef() {

}

func (reqs requests) IncrRef() {

}
