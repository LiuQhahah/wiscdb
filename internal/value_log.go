package internal

import (
	"github.com/dgraph-io/ristretto/v2/z"
	"sync"
	"sync/atomic"
	"wiscdb/y"
)

type valueLog struct {
	dirPath            string
	filesLock          sync.RWMutex
	filesMap           map[uint32]*logFile
	maxFid             uint32
	filesToBeDeleted   []uint32
	numActiveIterators atomic.Int32
	db                 *DB
	writableLogOffset  atomic.Uint32
	numEntriesWritten  uint32
	opt                Options
	garbageCh          chan struct{}
	discardStats       *discardStats
}

func (vlog *valueLog) sync() error {
	return nil
}
func (vlog *valueLog) dropAll() (int, error) {
	return 0, nil
}
func (vlog *valueLog) Close() error {
	return nil
}

func (vlog *valueLog) wOffset() uint32 {
	return 0
}

func (vlog *valueLog) validateWrites(reqs []*request) error {
	return nil
}

func (vlog *valueLog) write(reqs []*request) error {
	return nil
}

func (vlog *valueLog) createVlogFile() (*logFile, error) {
	return nil, nil
}

func estimateRequestSize(req *request) uint64 {
	return 0
}
func (vlog *valueLog) fPath(fid uint32) string {
	return ""
}

func vlogFilePath(dirPath string, fid uint32) string {
	return ""
}

func errFile(err error, path string, msg string) error {
	return nil
}

func (vlog *valueLog) getFileRLocked(vp valuePointer) (*logFile, error) {
	return &logFile{}, nil
}

func (vlog *valueLog) readValueBytes(vp valuePointer) ([]byte, *logFile, error) {
	return nil, &logFile{}, nil
}

func (vlog *valueLog) Read(vp valuePointer, _ *y.Slice) ([]byte, func(), error) {
	return nil, func() {

	}, nil
}

func runCallback(cb func()) {

}
func (vlog *valueLog) getUnlockCallback(lf *logFile) func() {
	return func() {

	}
}

func (vlog *valueLog) decrIteratorCount() error {
	return nil
}

func (vlog *valueLog) deleteLogFile(lf *logFile) error {
	return nil
}

func (vlog *valueLog) incrIteratorCount() {

}

func (vlog *valueLog) init(db *DB) {

}

func (vlog *valueLog) open(db *DB) error {
	return nil
}

func (vlog *valueLog) sortedFids() []uint32 {
	return nil
}

func (vlog *valueLog) populateFilesMap() error {
	return nil
}

func (vlog *valueLog) updateDiscardsStats(stats map[uint32]int64) {

}

func (vlog *valueLog) waitOnGC(lc *z.Closer) {

}

func (vlog *valueLog) runGC(discardRatio float64) error {
	return nil
}

func (vlog *valueLog) doRunGC(lf *logFile) error {
	return nil
}

func (vlog *valueLog) rewrite(f *logFile) error {
	return nil
}

func (vlog *valueLog) iteratorCount() int {
	return 0
}

func discardEntry(e Entry, vs y.ValueStruct, db *DB) bool {
	return false
}

func (vlog *valueLog) pickLog(discardRatio float64) *logFile {
	return &logFile{}
}
