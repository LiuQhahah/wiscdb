package internal

import (
	"github.com/dgraph-io/ristretto/v2/z"
	"sync"
)

type DB struct {
	testOnlyExtensions
	lock          sync.RWMutex
	dirLockGuard  *directoryLockGuard
	valueDisGuard *directoryLockGuard
	closers       closer
	mt            *memTable
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

func (db *DB) openMemTable(fid, flags int) (*memTable, error) {
	return nil, nil
}

func (db *DB) mtFilePath(fid int) string {
	return ""
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
