package main

import (
	otrace "go.opencensus.io/trace"
	"wiscdb/table"
	"wiscdb/y"
)

type compactDef struct {
	span         *otrace.Span
	compactorId  int
	t            targets
	p            compactionPriority
	thisLevel    *levelHandler
	nextLevel    *levelHandler
	top          []*table.Table
	bot          []*table.Table
	thisRange    keyRange
	nextRange    keyRange
	splits       []keyRange
	thisSize     int64
	dropPrefixes [][]byte
}

func (cd *compactDef) lockLevels() {
	cd.thisLevel.RLock()
	cd.nextLevel.RLock()
}

func (cd *compactDef) unlockLevels() {
	cd.nextLevel.RUnlock()
	cd.thisLevel.RUnlock()
}

func (cd *compactDef) allTables() []*table.Table {
	return nil
}

func appendIteratorReversed(out []y.Iterator, th []*table.Table, opt int) []y.Iterator {
	return nil
}
