package main

import (
	"encoding/binary"
	"github.com/dgraph-io/ristretto/v2/z"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"wiscdb/y"
)

// discardStats 保持追踪数据的数量,被删除的logFile的数量.
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
	//读取Discard文件
	fName := filepath.Join(opt.ValueDir, discardFName)
	// 1<<20 有1MB
	mf, err := z.OpenMmapFile(fName, os.O_CREATE|os.O_RDWR, 1<<20)
	lf := &discardStats{
		MmapFile: mf,
		opt:      opt,
	}
	if err == z.NewFile {
		lf.zeroOut()
	} else if err != nil {
		return nil, y.Wrapf(err, "while opening file: %s\n", discardFName)
	}

	for slot := 0; slot < lf.maxSlot(); slot++ {
		if lf.get(16*slot) == 0 {
			lf.nextEmptySlot = slot
			break
		}
	}
	sort.Sort(lf)
	opt.Infof("Discard stats nextEmptySlot:%d\n", lf.nextEmptySlot)
	return lf, nil
}

func (lf *discardStats) zeroOut() {

}

func (lf *discardStats) maxSlot() int {
	return 0
}
func (lf *discardStats) set(offset int, val uint64) {
	binary.BigEndian.PutUint64(lf.Data[offset:offset+8], val)
}

func (lf *discardStats) get(offset int) uint64 {
	return binary.BigEndian.Uint64(lf.Data[offset : offset+8])
}

// 将文件写到被丢弃的状态
func (lf *discardStats) Update(fidu uint32, discard int64) int64 {
	fid := uint64(fidu)
	lf.Lock()
	defer lf.Unlock()

	idx := sort.Search(lf.nextEmptySlot, func(slot int) bool {
		return lf.get(slot*16) >= fid
	})
	if idx < lf.nextEmptySlot && lf.get(idx*16) == fid {
		off := idx*16 + 8
		curDisc := lf.get(off)
		if discard == 0 {
			return int64(curDisc)
		}
		if discard < 0 {
			lf.set(off, 0)
			return 0
		}
		lf.set(off, curDisc+uint64(discard))
		return int64(curDisc + uint64(discard))
	}
	if discard <= 0 {
		return 0
	}
	idx = lf.nextEmptySlot
	lf.set(idx*16, fid)
	lf.set(idx*16+8, uint64(discard))

	lf.nextEmptySlot++
	for lf.nextEmptySlot >= lf.maxSlot() {
		y.Check(lf.Truncate(2 * int64(len(lf.Data))))
	}
	lf.zeroOut()
	sort.Sort(lf)
	return discard
}

func (lf *discardStats) MaxDiscard() (uint32, int64) {
	return 0, 0
}
