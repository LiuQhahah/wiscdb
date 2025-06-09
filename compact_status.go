package main

import "sync"

type compactStatus struct {
	sync.RWMutex
	levels []*levelCompactStatus
	tables map[uint64]struct{}
}

func (cs *compactStatus) overlapsWith(level int, this keyRange) bool {
	cs.RLock()
	defer cs.RUnlock()

	// 根据level获取对应的 levelCompactStatus
	thisLevel := cs.levels[level]
	// 根据传入的key判断 是否存在重叠的key
	return thisLevel.overlapsWith(this)
}

func (cs *compactStatus) delSize(l int) int64 {
	return 0
}

func (cs *compactStatus) delete(cd compactDef) {

}

func (cs *compactStatus) compareAndAdd(_ thisAndNextLevelRLocked, cd compactDef) bool {
	return false
}
