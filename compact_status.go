package main

import "sync"

type compactStatus struct {
	sync.RWMutex
	levels []*levelCompactStatus
	tables map[uint64]struct{}
}

func (cs *compactStatus) overlapsWith(level int, this keyRange) bool {
	return false
}

func (cs *compactStatus) delSize(l int) int64 {
	return 0
}

func (cs *compactStatus) delete(cd compactDef) {

}

func (cs *compactStatus) compareAndAdd(_ thisAndNextLevelRLocked, cd compactDef) bool {
	return false
}
