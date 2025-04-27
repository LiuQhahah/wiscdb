package internal

import (
	"bytes"
	"github.com/dgraph-io/ristretto/v2/z"
	"github.com/pkg/errors"
	"hash/crc32"
	"os"
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
	return vlog.writableLogOffset.Load()
}

// 遍历每一个request中的entryies如果当前value log的偏移量超过了设定的值比如1GB
// 那么就会重置偏移量
func (vlog *valueLog) validateWrites(reqs []*request) error {
	vlogOffset := uint64(vlog.wOffset())
	for _, req := range reqs {
		size := estimateRequestSize(req)
		estimatedVlogOffset := vlogOffset + size
		//MaxUint32 1<<32 = 4GB-1byte
		if estimatedVlogOffset > uint64(maxVlogFileSize) {
			return errors.Errorf("Request size offset %d is bigger than maximum offset %d", estimatedVlogOffset, maxVlogFileSize)
		}
		//default size:ValueLogFileSize=1GB
		if estimatedVlogOffset >= uint64(vlog.opt.ValueLogFileSize) {
			vlogOffset = 0
			continue
		}
		vlogOffset = estimatedVlogOffset
	}
	return nil
}

// 将请求写到disk和
func (vlog *valueLog) write(reqs []*request) error {
	if vlog.db.opt.InMemory {
		return nil
	}
	if err := vlog.validateWrites(reqs); err != nil {
		return y.Wrapf(err, "while validating writes")
	}
	vlog.filesLock.RLock()
	maxFid := vlog.maxFid
	//key value map 存储的是value log file
	curlf := vlog.filesMap[maxFid]
	vlog.filesLock.RUnlock()
	defer func() {
		if vlog.opt.SyncWrites {
			if err := curlf.Sync(); err != nil {
				vlog.opt.Errorf("Error while curlf sync: %v\n", err)
			}
		}
	}()

	write := func(buf *bytes.Buffer) error {
		if buf.Len() == 0 {
			return nil
		}
		n := uint32(buf.Len())
		endOffset := vlog.writableLogOffset.Add(n)
		if int(endOffset) >= len(curlf.Data) {
			if err := curlf.Truncate(int64(endOffset)); err != nil {
				return err
			}
		}

		start := int(endOffset - n)
		//将buf存储到curlf的data中
		y.AssertTrue(copy(curlf.Data[start:], buf.Bytes()) == int(n))

		curlf.size.Store(endOffset)
		return nil
	}
	//将value log写到磁盘中
	toDisk := func() error {
		//判断当前value log file的文件偏移量如果超过设置的value log则写道log file中
		if vlog.wOffset() > uint32(vlog.opt.ValueLogFileSize) || vlog.numEntriesWritten > vlog.opt.ValueLogMaxEntries {
			if err := curlf.doneWriting(vlog.wOffset()); err != nil {
				return err
			}
			newlf, err := vlog.createVlogFile()
			if err != nil {
				return err
			}
			curlf = newlf
		}
		return nil
	}

	buf := new(bytes.Buffer)
	for i := range reqs {
		b := reqs[i]
		b.Ptrs = b.Ptrs[:0]
		var written, bytesWritten int
		valueSizes := make([]int64, 0, len(b.Entries))
		for j := range b.Entries {
			buf.Reset()
			e := b.Entries[j]
			valueSizes = append(valueSizes, int64(len(e.Value)))
			if e.skipVlogAndSetThreshold(vlog.db.valueThreshold()) {
				b.Ptrs = append(b.Ptrs, valuePointer{})
				continue
			}
			var p valuePointer
			p.Fid = curlf.fid
			p.Offset = vlog.wOffset()
			tmpMeta := e.meta
			e.meta = e.meta &^ (bitTxn | bitFinTxn)
			pLen, err := curlf.encodeEntry(buf, e, p.Offset)
			if err != nil {
				return err
			}
			e.meta = tmpMeta
			p.Len = uint32(pLen)
			b.Ptrs = append(b.Ptrs, p)
			if err := write(buf); err != nil {
				return err
			}
			written++
			bytesWritten += buf.Len()
		}
		y.NumWritesVlogAdd(vlog.opt.MetricsEnabled, int64(written))
		y.NumBytesWrittenVlogAdd(vlog.opt.MetricsEnabled, int64(bytesWritten))

		vlog.numEntriesWritten += uint32(written)
		vlog.db.threshold.update(valueSizes)

		//将value log写到磁盘中
		if err := toDisk(); err != nil {
			return err
		}
	}

	return toDisk()
}

// 创建一个新的value log file
func (vlog *valueLog) createVlogFile() (*valueLogFile, error) {
	fid := vlog.maxFid + 1
	path := vlog.fPath(fid)
	vlogFile := &valueLogFile{
		fid:      fid,
		path:     path,
		registry: vlog.db.registry,
		writeAt:  vlogHeaderSize,
		opt:      vlog.opt,
	}
	err := vlogFile.open(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 2*vlog.opt.ValueLogFileSize)
	if err != z.NewFile && err != nil {
		return nil, err
	}
	vlog.filesLock.Lock()
	vlog.filesMap[fid] = vlogFile
	y.AssertTrue(vlog.maxFid < fid)
	vlog.maxFid = fid
	vlog.writableLogOffset.Store(vlogHeaderSize)
	vlog.numEntriesWritten = 0
	vlog.filesLock.Unlock()
	return vlogFile, nil
}

// 一个request中包含多个entry即多个key-value
// 一个entry的size包含header+key的尺寸、value的尺寸以及校验和的大小
func estimateRequestSize(req *request) uint64 {
	size := uint64(0)
	for _, e := range req.Entries {
		size += uint64(maxHeaderSize + len(e.Key) + len(e.Value) + crc32.Size)
	}
	return size
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
