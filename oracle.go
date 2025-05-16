package main

import (
	"context"
	"github.com/dgraph-io/ristretto/v2/z"
	"sync"
	"wiscdb/y"
)

type oracle struct {
	isManaged       bool
	detectConflicts bool
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

// oracle 可以做冲突检测
func newOracle(opt Options) *oracle {
	orc := &oracle{
		isManaged:       opt.MetricsEnabled,
		detectConflicts: opt.DetectConflicts,
		readMark:        &y.WaterMark{Name: "wisc.PendingReads"},
		txnMark:         &y.WaterMark{Name: "wisc.TxnTimeStamp"},
		closer:          z.NewCloser(2),
	}
	orc.readMark.Init(orc.closer)
	orc.txnMark.Init(orc.closer)
	return orc
}

func (o *oracle) Stop() {
	o.closer.SignalAndWait()
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
	o.Lock()
	defer o.Unlock()
	return o.nextTxnTs
}

func (o *oracle) incrementNextTs() {
	o.Lock()
	defer o.Unlock()
	o.nextTxnTs++
}

func (o *oracle) setDiscardTs(ts uint64) {
	o.Lock()
	defer o.Unlock()
	o.discardTs = ts
	o.cleanupCommittedTransactions()
}

func (o *oracle) cleanupCommittedTransactions() {
	if !o.detectConflicts {
		return
	}
	var maxReadTs uint64
	if o.isManaged {
		maxReadTs = o.discardTs
	} else {
		maxReadTs = o.readMark.DoneUntil()
	}
	y.AssertTrue(maxReadTs >= o.lastCleanupTs)
	if maxReadTs == o.lastCleanupTs {
		return
	}
	o.lastCleanupTs = maxReadTs
	tmp := o.committedTxns[:0]
	for _, txn := range o.committedTxns {
		if txn.ts <= maxReadTs {
			continue
		}
		tmp = append(tmp, txn)
	}
	o.committedTxns = tmp
}

func (o *oracle) doneCommit(cts uint64) {
	if o.isManaged {
		return
	}
	o.txnMark.Done(cts)
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
	if o.detectConflicts {
		o.committedTxns = append(o.committedTxns, committedTxn{
			ts:           ts,
			conflictKeys: txn.conflictKeys,
		})
	}
	return ts, true
}
