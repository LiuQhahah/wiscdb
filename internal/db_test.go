package internal

import (
	"reflect"
	"sync"
	"testing"
)

func TestDB_BanNamespace(t *testing.T) {
	type fields struct {
		testOnlyExtensions testOnlyExtensions
		lock               sync.RWMutex
		dirLockGuard       *directoryLockGuard
		valueDisGuard      *directoryLockGuard
		closers            closer
		mt                 *memTable
	}
	type args struct {
		ns uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &DB{
				testOnlyExtensions: tt.fields.testOnlyExtensions,
				lock:               tt.fields.lock,
				dirLockGuard:       tt.fields.dirLockGuard,
				valueDisGuard:      tt.fields.valueDisGuard,
				closers:            tt.fields.closers,
				mt:                 tt.fields.mt,
			}
			if err := db.BanNamespace(tt.args.ns); (err != nil) != tt.wantErr {
				t.Errorf("BanNamespace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDB_IsClosed(t *testing.T) {
	type fields struct {
		testOnlyExtensions testOnlyExtensions
		lock               sync.RWMutex
		dirLockGuard       *directoryLockGuard
		valueDisGuard      *directoryLockGuard
		closers            closer
		mt                 *memTable
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &DB{
				testOnlyExtensions: tt.fields.testOnlyExtensions,
				lock:               tt.fields.lock,
				dirLockGuard:       tt.fields.dirLockGuard,
				valueDisGuard:      tt.fields.valueDisGuard,
				closers:            tt.fields.closers,
				mt:                 tt.fields.mt,
			}
			if got := db.IsClosed(); got != tt.want {
				t.Errorf("IsClosed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDB_NewStream(t *testing.T) {
	type fields struct {
		testOnlyExtensions testOnlyExtensions
		lock               sync.RWMutex
		dirLockGuard       *directoryLockGuard
		valueDisGuard      *directoryLockGuard
		closers            closer
		mt                 *memTable
	}
	tests := []struct {
		name   string
		fields fields
		want   *Stream
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &DB{
				testOnlyExtensions: tt.fields.testOnlyExtensions,
				lock:               tt.fields.lock,
				dirLockGuard:       tt.fields.dirLockGuard,
				valueDisGuard:      tt.fields.valueDisGuard,
				closers:            tt.fields.closers,
				mt:                 tt.fields.mt,
			}
			if got := db.NewStream(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStream() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDB_NewStreamAt(t *testing.T) {
	type fields struct {
		testOnlyExtensions testOnlyExtensions
		lock               sync.RWMutex
		dirLockGuard       *directoryLockGuard
		valueDisGuard      *directoryLockGuard
		closers            closer
		mt                 *memTable
	}
	type args struct {
		readTs uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Stream
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &DB{
				testOnlyExtensions: tt.fields.testOnlyExtensions,
				lock:               tt.fields.lock,
				dirLockGuard:       tt.fields.dirLockGuard,
				valueDisGuard:      tt.fields.valueDisGuard,
				closers:            tt.fields.closers,
				mt:                 tt.fields.mt,
			}
			if got := db.NewStreamAt(tt.args.readTs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStreamAt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDB_NewTransaction(t *testing.T) {
	type fields struct {
		testOnlyExtensions testOnlyExtensions
		lock               sync.RWMutex
		dirLockGuard       *directoryLockGuard
		valueDisGuard      *directoryLockGuard
		closers            closer
		mt                 *memTable
	}
	type args struct {
		update bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Txn
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &DB{
				testOnlyExtensions: tt.fields.testOnlyExtensions,
				lock:               tt.fields.lock,
				dirLockGuard:       tt.fields.dirLockGuard,
				valueDisGuard:      tt.fields.valueDisGuard,
				closers:            tt.fields.closers,
				mt:                 tt.fields.mt,
			}
			if got := db.NewTransaction(tt.args.update); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDB_Update(t *testing.T) {
	type fields struct {
		testOnlyExtensions testOnlyExtensions
		lock               sync.RWMutex
		dirLockGuard       *directoryLockGuard
		valueDisGuard      *directoryLockGuard
		closers            closer
		mt                 *memTable
	}
	type args struct {
		fn func(txn *Txn) error
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &DB{
				testOnlyExtensions: tt.fields.testOnlyExtensions,
				lock:               tt.fields.lock,
				dirLockGuard:       tt.fields.dirLockGuard,
				valueDisGuard:      tt.fields.valueDisGuard,
				closers:            tt.fields.closers,
				mt:                 tt.fields.mt,
			}
			if err := db.Update(tt.args.fn); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDB_View(t *testing.T) {
	type fields struct {
		testOnlyExtensions testOnlyExtensions
		lock               sync.RWMutex
		dirLockGuard       *directoryLockGuard
		valueDisGuard      *directoryLockGuard
		closers            closer
		mt                 *memTable
	}
	type args struct {
		fb func(txn *Txn) error
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &DB{
				testOnlyExtensions: tt.fields.testOnlyExtensions,
				lock:               tt.fields.lock,
				dirLockGuard:       tt.fields.dirLockGuard,
				valueDisGuard:      tt.fields.valueDisGuard,
				closers:            tt.fields.closers,
				mt:                 tt.fields.mt,
			}
			if err := db.View(tt.args.fb); (err != nil) != tt.wantErr {
				t.Errorf("View() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDB_initBannedNamespaces(t *testing.T) {
	type fields struct {
		testOnlyExtensions testOnlyExtensions
		lock               sync.RWMutex
		dirLockGuard       *directoryLockGuard
		valueDisGuard      *directoryLockGuard
		closers            closer
		mt                 *memTable
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &DB{
				testOnlyExtensions: tt.fields.testOnlyExtensions,
				lock:               tt.fields.lock,
				dirLockGuard:       tt.fields.dirLockGuard,
				valueDisGuard:      tt.fields.valueDisGuard,
				closers:            tt.fields.closers,
				mt:                 tt.fields.mt,
			}
			if err := db.initBannedNamespaces(); (err != nil) != tt.wantErr {
				t.Errorf("initBannedNamespaces() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDB_isBanned(t *testing.T) {
	type fields struct {
		testOnlyExtensions testOnlyExtensions
		lock               sync.RWMutex
		dirLockGuard       *directoryLockGuard
		valueDisGuard      *directoryLockGuard
		closers            closer
		mt                 *memTable
	}
	type args struct {
		key []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &DB{
				testOnlyExtensions: tt.fields.testOnlyExtensions,
				lock:               tt.fields.lock,
				dirLockGuard:       tt.fields.dirLockGuard,
				valueDisGuard:      tt.fields.valueDisGuard,
				closers:            tt.fields.closers,
				mt:                 tt.fields.mt,
			}
			if err := db.isBanned(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("isBanned() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDB_mtFilePath(t *testing.T) {
	type fields struct {
		testOnlyExtensions testOnlyExtensions
		lock               sync.RWMutex
		dirLockGuard       *directoryLockGuard
		valueDisGuard      *directoryLockGuard
		closers            closer
		mt                 *memTable
	}
	type args struct {
		fid int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &DB{
				testOnlyExtensions: tt.fields.testOnlyExtensions,
				lock:               tt.fields.lock,
				dirLockGuard:       tt.fields.dirLockGuard,
				valueDisGuard:      tt.fields.valueDisGuard,
				closers:            tt.fields.closers,
				mt:                 tt.fields.mt,
			}
			if got := db.getFilePathWithFid(tt.args.fid); got != tt.want {
				t.Errorf("getFilePathWithFid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDB_newMemTable(t *testing.T) {
	type fields struct {
		testOnlyExtensions testOnlyExtensions
		lock               sync.RWMutex
		dirLockGuard       *directoryLockGuard
		valueDisGuard      *directoryLockGuard
		closers            closer
		mt                 *memTable
	}
	tests := []struct {
		name    string
		fields  fields
		want    *memTable
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &DB{
				testOnlyExtensions: tt.fields.testOnlyExtensions,
				lock:               tt.fields.lock,
				dirLockGuard:       tt.fields.dirLockGuard,
				valueDisGuard:      tt.fields.valueDisGuard,
				closers:            tt.fields.closers,
				mt:                 tt.fields.mt,
			}
			got, err := db.newMemTable()
			if (err != nil) != tt.wantErr {
				t.Errorf("newMemTable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newMemTable() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDB_newStream(t *testing.T) {
	type fields struct {
		testOnlyExtensions testOnlyExtensions
		lock               sync.RWMutex
		dirLockGuard       *directoryLockGuard
		valueDisGuard      *directoryLockGuard
		closers            closer
		mt                 *memTable
	}
	tests := []struct {
		name   string
		fields fields
		want   *Stream
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &DB{
				testOnlyExtensions: tt.fields.testOnlyExtensions,
				lock:               tt.fields.lock,
				dirLockGuard:       tt.fields.dirLockGuard,
				valueDisGuard:      tt.fields.valueDisGuard,
				closers:            tt.fields.closers,
				mt:                 tt.fields.mt,
			}
			if got := db.newStream(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newStream() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDB_newTransaction(t *testing.T) {
	type fields struct {
		testOnlyExtensions testOnlyExtensions
		lock               sync.RWMutex
		dirLockGuard       *directoryLockGuard
		valueDisGuard      *directoryLockGuard
		closers            closer
		mt                 *memTable
	}
	type args struct {
		update    bool
		isManaged bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Txn
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &DB{
				testOnlyExtensions: tt.fields.testOnlyExtensions,
				lock:               tt.fields.lock,
				dirLockGuard:       tt.fields.dirLockGuard,
				valueDisGuard:      tt.fields.valueDisGuard,
				closers:            tt.fields.closers,
				mt:                 tt.fields.mt,
			}
			if got := db.newTransaction(tt.args.update, tt.args.isManaged); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDB_openMemTable(t *testing.T) {
	type fields struct {
		testOnlyExtensions testOnlyExtensions
		lock               sync.RWMutex
		dirLockGuard       *directoryLockGuard
		valueDisGuard      *directoryLockGuard
		closers            closer
		mt                 *memTable
	}
	type args struct {
		fid   int
		flags int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *memTable
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &DB{
				testOnlyExtensions: tt.fields.testOnlyExtensions,
				lock:               tt.fields.lock,
				dirLockGuard:       tt.fields.dirLockGuard,
				valueDisGuard:      tt.fields.valueDisGuard,
				closers:            tt.fields.closers,
				mt:                 tt.fields.mt,
			}
			got, err := db.openMemTable(tt.args.fid, tt.args.flags)
			if (err != nil) != tt.wantErr {
				t.Errorf("openMemTable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("openMemTable() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDB_openMemTables(t *testing.T) {
	type fields struct {
		testOnlyExtensions testOnlyExtensions
		lock               sync.RWMutex
		dirLockGuard       *directoryLockGuard
		valueDisGuard      *directoryLockGuard
		closers            closer
		mt                 *memTable
	}
	type args struct {
		opt Options
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &DB{
				testOnlyExtensions: tt.fields.testOnlyExtensions,
				lock:               tt.fields.lock,
				dirLockGuard:       tt.fields.dirLockGuard,
				valueDisGuard:      tt.fields.valueDisGuard,
				closers:            tt.fields.closers,
				mt:                 tt.fields.mt,
			}
			if err := db.openMemTables(tt.args.opt); (err != nil) != tt.wantErr {
				t.Errorf("openMemTables() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDB_sendToWriteCh(t *testing.T) {
	type fields struct {
		testOnlyExtensions testOnlyExtensions
		lock               sync.RWMutex
		dirLockGuard       *directoryLockGuard
		valueDisGuard      *directoryLockGuard
		closers            closer
		mt                 *memTable
	}
	type args struct {
		entries []*Entry
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *request
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &DB{
				testOnlyExtensions: tt.fields.testOnlyExtensions,
				lock:               tt.fields.lock,
				dirLockGuard:       tt.fields.dirLockGuard,
				valueDisGuard:      tt.fields.valueDisGuard,
				closers:            tt.fields.closers,
				mt:                 tt.fields.mt,
			}
			got, err := db.sendToWriteCh(tt.args.entries)
			if (err != nil) != tt.wantErr {
				t.Errorf("sendToWriteCh() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sendToWriteCh() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDB_valueThreshold(t *testing.T) {
	type fields struct {
		testOnlyExtensions testOnlyExtensions
		lock               sync.RWMutex
		dirLockGuard       *directoryLockGuard
		valueDisGuard      *directoryLockGuard
		closers            closer
		mt                 *memTable
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &DB{
				testOnlyExtensions: tt.fields.testOnlyExtensions,
				lock:               tt.fields.lock,
				dirLockGuard:       tt.fields.dirLockGuard,
				valueDisGuard:      tt.fields.valueDisGuard,
				closers:            tt.fields.closers,
				mt:                 tt.fields.mt,
			}
			if got := db.valueThreshold(); got != tt.want {
				t.Errorf("valueThreshold() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_arenaSize(t *testing.T) {
	type args struct {
		opt Options
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "test",
			args: args{DefaultOptions("../")},
			want: 64 << 20,
		},
		{
			name: "test2",
			args: args{opt: Options{
				MemTableSize:  64 << 20,
				maxBatchSize:  10,
				maxBatchCount: 100,
			}},
			want: 64<<20 + 10 + 88*100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := arenaSize(tt.args.opt); got != tt.want {
				t.Errorf("arenaSize() = %v, want %v", got, tt.want)
			}
		})
	}
}
