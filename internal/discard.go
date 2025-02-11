package internal

import (
	"github.com/dgraph-io/ristretto/v2/z"
	"sync"
)

type discardStats struct {
	sync.Mutex
	*z.MmapFile
	opt           Options
	nextEmptySlot int
}

func (lf *discardStats) Iterate(f func(fid, stats uint64)) {

}

func (lf *discardStats) Len() int {
	return 0
}

func (lf *discardStats) Less(i, j int) bool {
	return false
}

func (lf *discardStats) Swap(i, j int) {

}

const discardFName string = "DISCARD"

func InitDiscardStats(opt Options) (*discardStats, error) {
	return &discardStats{}, nil
}

func (lf *discardStats) zeroOut() {

}

func (lf *discardStats) maxSlot() int {
	return 0
}
func (lf *discardStats) set(offset int, val uint64) {

}

func (lf *discardStats) get(offset int) uint64 {
	return 0
}

func (lf *discardStats) Update(fidu uint32, discard int64) int64 {
	return 0
}

func (lf *discardStats) MaxDiscard() (uint32, int64) {
	return 0, 0
}
