package internal

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/dgraph-io/ristretto/v2/z"
	"path/filepath"
	"sync"
	"wiscdb/skl"
	"wiscdb/y"
)

type DB struct {
	testOnlyExtensions
	lock          sync.RWMutex
	dirLockGuard  *directoryLockGuard
	valueDisGuard *directoryLockGuard
	closers       closer
	mt            *memTable
	opt           Options
	registry      *KeyRegistry
}

type closer struct {
	updateSize  *z.Closer
	compactors  *z.Closer
	memTable    *z.Closer
	writes      *z.Closer
	valueGC     *z.Closer
	pub         *z.Closer
	cacheHealth *z.Closer
}

func (db *DB) openMemTables(opt Options) error {
	return nil
}

func (db *DB) newMemTable() (*memTable, error) {
	return nil, nil
}

// create/create memtable with file id.
func (db *DB) openMemTable(fid, flags int) (*memTable, error) {
	filePath := db.getFilePathWithFid(fid)
	s := skl.NewSkipList(arenaSize(db.opt))
	memtable := &memTable{
		sl:  s,
		opt: db.opt,
		buf: &bytes.Buffer{},
	}
	memtable.wal = &valueLogFile{
		path:     filePath,
		fid:      uint32(fid),
		opt:      db.opt,
		writeAt:  vlogHeaderSize,
		registry: db.registry,
	}
	//open write ahead log
	walError := memtable.wal.open(filePath, flags, 2*db.opt.MemTableSize)
	if walError != nil && !errors.Is(walError, z.NewFile) {
		return nil, y.Wrapf(walError, "While opening memtable: %s", filePath)
	}

	if errors.Is(walError, z.NewFile) {
		return memtable, walError
	}
	err := memtable.UpdateSkipList()
	return memtable, y.Wrapf(err, "while updating skiplist")
}

// return file name with file if
// 1: filepath package
// 2. fmt package
// 3. const memTableExt:=".mem"
func (db *DB) getFilePathWithFid(fid int) string {
	return filepath.Join(db.opt.Dir, fmt.Sprintf("%05d%s", fid, memFileExt))
}

func (db *DB) View(fb func(txn *Txn) error) error {
	return nil
}

func (db *DB) NewTransaction(update bool) *Txn {
	return &Txn{}
}

func (db *DB) newTransaction(update, isManaged bool) *Txn {
	return &Txn{}
}

func (db *DB) Update(fn func(txn *Txn) error) error {
	return nil
}

func (db *DB) initBannedNamespaces() error {
	return nil
}

func (db *DB) IsClosed() bool {
	return false
}

func (db *DB) isBanned(key []byte) error {
	return nil
}

func (db *DB) BanNamespace(ns uint64) error {
	return nil
}

func (db *DB) sendToWriteCh(entries []*Entry) (*request, error) {
	return &request{}, nil
}

func (db *DB) valueThreshold() int64 {
	return 0
}

func (db *DB) newStream() *Stream {
	return &Stream{}
}

func (db *DB) NewStream() *Stream {
	return db.newStream()
}

func (db *DB) NewStreamAt(readTs uint64) *Stream {
	stream := db.newStream()
	stream.readTs = readTs
	return stream
}

/*
*

	arena Size: MemTableSize 指的是内存中的memtable

默认大小： 64<<20(64MiB) + 10 + 88*100

TODO: maxBatchSize 的作用
TODO: MaxNodeSize指的是什么
default size: MemTableSize:64MiB.
*/
func arenaSize(opt Options) int64 {
	return opt.MemTableSize + opt.maxBatchSize + opt.maxBatchCount*int64(skl.MaxNodeSize)
}
