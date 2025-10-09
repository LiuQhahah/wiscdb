package main

import (
	"sync"
	"wiscdb/y"
)

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
	cs.Lock()
	defer cs.Unlock()
	tl := cd.thisLevel.level
	y.AssertTruef(tl < len(cs.levels), "Got level %d. Max Levels: %d", tl, len(cs.levels))
	thisLevel := cs.levels[cd.thisLevel.level]
	nextLevel := cs.levels[cd.nextLevel.level]

	if thisLevel.overlapsWith(cd.thisRange) {
		return false
	}
	if nextLevel.overlapsWith(cd.nextRange) {
		return false
	}

	thisLevel.ranges = append(thisLevel.ranges, cd.thisRange)
	nextLevel.ranges = append(nextLevel.ranges, cd.nextRange)

	for _, t := range append(cd.top, cd.bot...) {
		cs.tables[t.ID()] = struct{}{}
	}
	return true
}
