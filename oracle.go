package main

import (
	"context"
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
func (o *oracle) readTs() uint64 {
	if o.isManaged {
		panic("ReadTs should not be retrieved for managed DB")
	}
	var readTs uint64
	o.Lock()
	readTs = o.nextTxnTs - 1
	o.readMark.Begin(readTs)
	o.Unlock()
	y.Check(o.txnMark.WaitForMark(context.Background(), readTs))
	return readTs
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

func (o *oracle) doneRead(txn *Txn) {
	if !txn.doneRead {
		txn.doneRead = true
		o.readMark.Done(txn.readTs)
	}
}

func (o *oracle) hasConflict(txn *Txn) bool {
	if len(txn.reads) == 0 {
		return false
	}
	for _, committedTxn := range o.committedTxns {
		if committedTxn.ts <= txn.readTs {
			continue
		}

		for _, ro := range txn.reads {
			if _, has := committedTxn.conflictKeys[ro]; has {
				return true
			}
		}
	}
	return false
}
func (o *oracle) newCommitTs(txn *Txn) (uint64, bool) {
	o.Lock()
	defer o.Unlock()
	if o.hasConflict(txn) {
		return 0, true
	}
	var ts uint64

	if !o.isManaged {
		o.doneRead(txn)
		o.cleanupCommittedTransactions()
		ts = o.nextTxnTs
		o.nextTxnTs++
		o.txnMark.Begin(ts)
	} else {
		ts = txn.commitTs
	}
	y.AssertTrue(ts >= o.lastCleanupTs)
	if o.deleteConflicts {
		o.committedTxns = append(o.committedTxns, committedTxn{
			ts:           ts,
			conflictKeys: txn.conflictKeys,
		})
	}
	return ts, true
}
