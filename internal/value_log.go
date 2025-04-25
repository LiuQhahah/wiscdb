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
	filesMap           map[uint32]*valueLogFile
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

func (vlog *valueLog) createVlogFile() (*valueLogFile, error) {
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

func (vlog *valueLog) getFileRLocked(vp valuePointer) (*valueLogFile, error) {
	return &valueLogFile{}, nil
}

func (vlog *valueLog) readValueBytes(vp valuePointer) ([]byte, *valueLogFile, error) {
	return nil, &valueLogFile{}, nil
}

func (vlog *valueLog) Read(vp valuePointer, _ *y.Slice) ([]byte, func(), error) {
	return nil, func() {

	}, nil
}

func runCallback(cb func()) {

}
func (vlog *valueLog) getUnlockCallback(lf *valueLogFile) func() {
	return func() {

	}
}

// check file count and delete file
func (vlog *valueLog) decrIteratorCount() error {
	//将活跃迭代数量减1
	num := vlog.numActiveIterators.Add(-1)
	if num != 0 {
		return nil
	}
	//加锁
	vlog.filesLock.Lock()
	//创建数组,数组长度是valuelog中将要被删除的文件数量
	//数组中存储的是要被删除的文件的id
	lfs := make([]*valueLogFile, 0, len(vlog.filesToBeDeleted))
	for _, id := range vlog.filesToBeDeleted {
		lfs = append(lfs, vlog.filesMap[id])
		//使用go自带的删除函数
		delete(vlog.filesMap, id)
	}
	//删除完后，将valuelog的字段设置成空
	vlog.filesToBeDeleted = nil
	//解锁
	vlog.filesLock.Unlock()
	for _, lf := range lfs {
		//逐个删除logfile
		if err := vlog.deleteLogFile(lf); err != nil {
			return err
		}
	}
	return nil
}

// 将value log文件写discardStats中
func (vlog *valueLog) deleteLogFile(lf *valueLogFile) error {
	if lf == nil {
		return nil
	}
	lf.lock.Lock()
	defer lf.lock.Lock()
	vlog.discardStats.Update(lf.fid, -1)
	return lf.Delete()
}

func (vlog *valueLog) incrIteratorCount() {
	vlog.numActiveIterators.Add(1)
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

func (vlog *valueLog) doRunGC(lf *valueLogFile) error {
	return nil
}

func (vlog *valueLog) rewrite(f *valueLogFile) error {
	return nil
}

func (vlog *valueLog) iteratorCount() int {
	return 0
}

func discardEntry(e Entry, vs y.ValueStruct, db *DB) bool {
	return false
}

func (vlog *valueLog) pickLog(discardRatio float64) *valueLogFile {
	return &valueLogFile{}
}
