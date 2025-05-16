package main

import (
	"context"
	"github.com/dgraph-io/ristretto/v2/z"
	"sync"
	"wiscdb/y"
)

type oracle struct {
	isManaged       bool // TODO: isManaged的含义
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
	// 遍历提交的事务
	for _, txn := range o.committedTxns {
		// 如果事务的ts比最大的ReadTs大,则将事务更新到tmp中
		if txn.ts <= maxReadTs {
			continue
		}
		tmp = append(tmp, txn)
	}
	// 将事务更新到committedTxns中
	o.committedTxns = tmp
}

// 如果oracle 的isManaged是true,则直接返回
func (o *oracle) doneCommit(cts uint64) {
	if o.isManaged {
		return
	}
	// 否则更新下txnMark
	o.txnMark.Done(cts)
}

func (o *oracle) doneRead(txn *Txn) {
	// 判断当前事务的doneRead是否为ture
	if !txn.doneRead {
		txn.doneRead = true
		// 如果为false则设置成true
		// 同时更新oracle的readMark状态,传入的参数是事务的readTs
		o.readMark.Done(txn.readTs)
	}
}

func (o *oracle) hasConflict(txn *Txn) bool {
	// 遍历事务中reads的值，如果为0表明没有冲突
	if len(txn.reads) == 0 {
		return false
	}
	// 遍历提交的事务
	for _, committedTxn := range o.committedTxns {
		if committedTxn.ts <= txn.readTs {
			continue
		}

		// 遍历事务中的被hash的key
		for _, ro := range txn.reads {
			// 如果被hash的key存储在conflictKeys中则返回true
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
	// 检测是否有冲突
	if o.hasConflict(txn) {
		return 0, true
	}
	var ts uint64

	//如果isManaged是false
	if !o.isManaged {
		o.doneRead(txn)
		o.cleanupCommittedTransactions()
		ts = o.nextTxnTs
		o.nextTxnTs++
		o.txnMark.Begin(ts)
	} else {
		// 如果isManaged是true的,直接将事务的commitTs复制给Ts
		ts = txn.commitTs
	}
	y.AssertTrue(ts >= o.lastCleanupTs)
	if o.detectConflicts {
		// 将需要提交的事务追加到oracle的commitedTxns中
		o.committedTxns = append(o.committedTxns, committedTxn{
			ts:           ts,
			conflictKeys: txn.conflictKeys,
		})
	}
	return ts, true
}
