package internal

import (
	"github.com/dgraph-io/ristretto/v2/z"
	"sync"
	"wiscdb/y"
)

type oracle struct {
	isManaged       bool
	deleteConflicts bool
	sync.Mutex
	writeChLock   sync.Mutex
	nextTxnTs     uint64
	txnMark       *y.WaterMark
	discardTs     uint64
	readMark      *y.WaterMark
	committedTxns []committedTxn
	lastCleanupTs uint64
	closer        *z.Closer
}

func newOracle(opt Options) *oracle {
	return &oracle{}
}

func (o *oracle) Stop() {

}
func (o *oracle) readts() uint64 {
	return 0
}
func (o *oracle) nextTs() uint64 {
	return 0
}

func (o *oracle) incrementNextTs() {

}

func (o *oracle) setDiscard(ts uint64) {

}

func (o *oracle) cleanupCommittedTransactions() {

}

func (o *oracle) doneCommit(cts uint64) {

}
