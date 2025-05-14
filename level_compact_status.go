package main

type levelCompactStatus struct {
	ranges  []keyRange
	delSize int64
}

func (lcs *levelCompactStatus) debug() string {
	return ""
}

func (lcs *levelCompactStatus) overlapsWith(dst keyRange) bool {
	return false
}

func (lcs *levelCompactStatus) remove(dst keyRange) bool {
	return false
}
