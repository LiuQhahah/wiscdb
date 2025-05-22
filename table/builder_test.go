package table

import (
	"github.com/dgraph-io/ristretto/v2/z"
	flatbuffers "github.com/google/flatbuffers/go"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"wiscdb/options"
	"wiscdb/pb"
	"wiscdb/y"
)

func TestBuilder_Add(t *testing.T) {
	type fields struct {
		alloc            *z.Allocator
		curBlock         *bblock
		compressedSize   atomic.Uint32
		uncompressedSize atomic.Uint32
		lenOffsets       uint32
		keyHashes        []uint32
		opts             *Options
		maxVersion       uint64
		onDiskSize       uint32
		staleDataSize    int
		wg               sync.WaitGroup
		blockChan        chan *bblock
		blockList        []*bblock
	}
	type args struct {
		key      []byte
		value    y.ValueStruct
		valueLen uint32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				alloc:            tt.fields.alloc,
				curBlock:         tt.fields.curBlock,
				compressedSize:   tt.fields.compressedSize,
				uncompressedSize: tt.fields.uncompressedSize,
				lenOffsets:       tt.fields.lenOffsets,
				keyHashes:        tt.fields.keyHashes,
				opts:             tt.fields.opts,
				maxVersion:       tt.fields.maxVersion,
				onDiskSize:       tt.fields.onDiskSize,
				staleDataSize:    tt.fields.staleDataSize,
				wg:               tt.fields.wg,
				blockChan:        tt.fields.blockChan,
				blockList:        tt.fields.blockList,
			}
			b.Add(tt.args.key, tt.args.value, tt.args.valueLen)
		})
	}
}

func TestBuilder_AddStaleKey(t *testing.T) {
	type fields struct {
		alloc            *z.Allocator
		curBlock         *bblock
		compressedSize   atomic.Uint32
		uncompressedSize atomic.Uint32
		lenOffsets       uint32
		keyHashes        []uint32
		opts             *Options
		maxVersion       uint64
		onDiskSize       uint32
		staleDataSize    int
		wg               sync.WaitGroup
		blockChan        chan *bblock
		blockList        []*bblock
	}
	type args struct {
		key      []byte
		v        y.ValueStruct
		valueLen uint32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				alloc:            tt.fields.alloc,
				curBlock:         tt.fields.curBlock,
				compressedSize:   tt.fields.compressedSize,
				uncompressedSize: tt.fields.uncompressedSize,
				lenOffsets:       tt.fields.lenOffsets,
				keyHashes:        tt.fields.keyHashes,
				opts:             tt.fields.opts,
				maxVersion:       tt.fields.maxVersion,
				onDiskSize:       tt.fields.onDiskSize,
				staleDataSize:    tt.fields.staleDataSize,
				wg:               tt.fields.wg,
				blockChan:        tt.fields.blockChan,
				blockList:        tt.fields.blockList,
			}
			b.AddStaleKey(tt.args.key, tt.args.v, tt.args.valueLen)
		})
	}
}

func TestBuilder_Close(t *testing.T) {
	type fields struct {
		alloc            *z.Allocator
		curBlock         *bblock
		compressedSize   atomic.Uint32
		uncompressedSize atomic.Uint32
		lenOffsets       uint32
		keyHashes        []uint32
		opts             *Options
		maxVersion       uint64
		onDiskSize       uint32
		staleDataSize    int
		wg               sync.WaitGroup
		blockChan        chan *bblock
		blockList        []*bblock
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				alloc:            tt.fields.alloc,
				curBlock:         tt.fields.curBlock,
				compressedSize:   tt.fields.compressedSize,
				uncompressedSize: tt.fields.uncompressedSize,
				lenOffsets:       tt.fields.lenOffsets,
				keyHashes:        tt.fields.keyHashes,
				opts:             tt.fields.opts,
				maxVersion:       tt.fields.maxVersion,
				onDiskSize:       tt.fields.onDiskSize,
				staleDataSize:    tt.fields.staleDataSize,
				wg:               tt.fields.wg,
				blockChan:        tt.fields.blockChan,
				blockList:        tt.fields.blockList,
			}
			b.Close()
		})
	}
}

func TestBuilder_CutDoneBuildData(t *testing.T) {
	type fields struct {
		alloc            *z.Allocator
		curBlock         *bblock
		compressedSize   atomic.Uint32
		uncompressedSize atomic.Uint32
		lenOffsets       uint32
		keyHashes        []uint32
		opts             *Options
		maxVersion       uint64
		onDiskSize       uint32
		staleDataSize    int
		wg               sync.WaitGroup
		blockChan        chan *bblock
		blockList        []*bblock
	}
	tests := []struct {
		name   string
		fields fields
		want   buildData
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				alloc:            tt.fields.alloc,
				curBlock:         tt.fields.curBlock,
				compressedSize:   tt.fields.compressedSize,
				uncompressedSize: tt.fields.uncompressedSize,
				lenOffsets:       tt.fields.lenOffsets,
				keyHashes:        tt.fields.keyHashes,
				opts:             tt.fields.opts,
				maxVersion:       tt.fields.maxVersion,
				onDiskSize:       tt.fields.onDiskSize,
				staleDataSize:    tt.fields.staleDataSize,
				wg:               tt.fields.wg,
				blockChan:        tt.fields.blockChan,
				blockList:        tt.fields.blockList,
			}
			if got := b.CutDoneBuildData(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CutDoneBuildData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_CutDownBuilder(t *testing.T) {
	type fields struct {
		alloc            *z.Allocator
		curBlock         *bblock
		compressedSize   atomic.Uint32
		uncompressedSize atomic.Uint32
		lenOffsets       uint32
		keyHashes        []uint32
		opts             *Options
		maxVersion       uint64
		onDiskSize       uint32
		staleDataSize    int
		wg               sync.WaitGroup
		blockChan        chan *bblock
		blockList        []*bblock
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				alloc:            tt.fields.alloc,
				curBlock:         tt.fields.curBlock,
				compressedSize:   tt.fields.compressedSize,
				uncompressedSize: tt.fields.uncompressedSize,
				lenOffsets:       tt.fields.lenOffsets,
				keyHashes:        tt.fields.keyHashes,
				opts:             tt.fields.opts,
				maxVersion:       tt.fields.maxVersion,
				onDiskSize:       tt.fields.onDiskSize,
				staleDataSize:    tt.fields.staleDataSize,
				wg:               tt.fields.wg,
				blockChan:        tt.fields.blockChan,
				blockList:        tt.fields.blockList,
			}
			if got := b.CutDownBuilder(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CutDownBuilder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_DataKey(t *testing.T) {
	type fields struct {
		alloc            *z.Allocator
		curBlock         *bblock
		compressedSize   atomic.Uint32
		uncompressedSize atomic.Uint32
		lenOffsets       uint32
		keyHashes        []uint32
		opts             *Options
		maxVersion       uint64
		onDiskSize       uint32
		staleDataSize    int
		wg               sync.WaitGroup
		blockChan        chan *bblock
		blockList        []*bblock
	}
	tests := []struct {
		name   string
		fields fields
		want   *pb.DataKey
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				alloc:            tt.fields.alloc,
				curBlock:         tt.fields.curBlock,
				compressedSize:   tt.fields.compressedSize,
				uncompressedSize: tt.fields.uncompressedSize,
				lenOffsets:       tt.fields.lenOffsets,
				keyHashes:        tt.fields.keyHashes,
				opts:             tt.fields.opts,
				maxVersion:       tt.fields.maxVersion,
				onDiskSize:       tt.fields.onDiskSize,
				staleDataSize:    tt.fields.staleDataSize,
				wg:               tt.fields.wg,
				blockChan:        tt.fields.blockChan,
				blockList:        tt.fields.blockList,
			}
			if got := b.DataKey(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DataKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_Empty(t *testing.T) {
	type fields struct {
		alloc            *z.Allocator
		curBlock         *bblock
		compressedSize   atomic.Uint32
		uncompressedSize atomic.Uint32
		lenOffsets       uint32
		keyHashes        []uint32
		opts             *Options
		maxVersion       uint64
		onDiskSize       uint32
		staleDataSize    int
		wg               sync.WaitGroup
		blockChan        chan *bblock
		blockList        []*bblock
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
			b := &Builder{
				alloc:            tt.fields.alloc,
				curBlock:         tt.fields.curBlock,
				compressedSize:   tt.fields.compressedSize,
				uncompressedSize: tt.fields.uncompressedSize,
				lenOffsets:       tt.fields.lenOffsets,
				keyHashes:        tt.fields.keyHashes,
				opts:             tt.fields.opts,
				maxVersion:       tt.fields.maxVersion,
				onDiskSize:       tt.fields.onDiskSize,
				staleDataSize:    tt.fields.staleDataSize,
				wg:               tt.fields.wg,
				blockChan:        tt.fields.blockChan,
				blockList:        tt.fields.blockList,
			}
			if got := b.Empty(); got != tt.want {
				t.Errorf("Empty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_Opts(t *testing.T) {
	type fields struct {
		alloc            *z.Allocator
		curBlock         *bblock
		compressedSize   atomic.Uint32
		uncompressedSize atomic.Uint32
		lenOffsets       uint32
		keyHashes        []uint32
		opts             *Options
		maxVersion       uint64
		onDiskSize       uint32
		staleDataSize    int
		wg               sync.WaitGroup
		blockChan        chan *bblock
		blockList        []*bblock
	}
	tests := []struct {
		name   string
		fields fields
		want   *Options
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				alloc:            tt.fields.alloc,
				curBlock:         tt.fields.curBlock,
				compressedSize:   tt.fields.compressedSize,
				uncompressedSize: tt.fields.uncompressedSize,
				lenOffsets:       tt.fields.lenOffsets,
				keyHashes:        tt.fields.keyHashes,
				opts:             tt.fields.opts,
				maxVersion:       tt.fields.maxVersion,
				onDiskSize:       tt.fields.onDiskSize,
				staleDataSize:    tt.fields.staleDataSize,
				wg:               tt.fields.wg,
				blockChan:        tt.fields.blockChan,
				blockList:        tt.fields.blockList,
			}
			if got := b.Opts(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Opts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_ReachedCapacity(t *testing.T) {
	type fields struct {
		alloc            *z.Allocator
		curBlock         *bblock
		compressedSize   atomic.Uint32
		uncompressedSize atomic.Uint32
		lenOffsets       uint32
		keyHashes        []uint32
		opts             *Options
		maxVersion       uint64
		onDiskSize       uint32
		staleDataSize    int
		wg               sync.WaitGroup
		blockChan        chan *bblock
		blockList        []*bblock
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
			b := &Builder{
				alloc:            tt.fields.alloc,
				curBlock:         tt.fields.curBlock,
				compressedSize:   tt.fields.compressedSize,
				uncompressedSize: tt.fields.uncompressedSize,
				lenOffsets:       tt.fields.lenOffsets,
				keyHashes:        tt.fields.keyHashes,
				opts:             tt.fields.opts,
				maxVersion:       tt.fields.maxVersion,
				onDiskSize:       tt.fields.onDiskSize,
				staleDataSize:    tt.fields.staleDataSize,
				wg:               tt.fields.wg,
				blockChan:        tt.fields.blockChan,
				blockList:        tt.fields.blockList,
			}
			if got := b.ReachedCapacity(); got != tt.want {
				t.Errorf("ReachedCapacity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_addHelper(t *testing.T) {
	type fields struct {
		alloc            *z.Allocator
		curBlock         *bblock
		compressedSize   atomic.Uint32
		uncompressedSize atomic.Uint32
		lenOffsets       uint32
		keyHashes        []uint32
		opts             *Options
		maxVersion       uint64
		onDiskSize       uint32
		staleDataSize    int
		wg               sync.WaitGroup
		blockChan        chan *bblock
		blockList        []*bblock
	}
	type args struct {
		key   []byte
		v     y.ValueStruct
		vpLen uint32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				alloc:            tt.fields.alloc,
				curBlock:         tt.fields.curBlock,
				compressedSize:   tt.fields.compressedSize,
				uncompressedSize: tt.fields.uncompressedSize,
				lenOffsets:       tt.fields.lenOffsets,
				keyHashes:        tt.fields.keyHashes,
				opts:             tt.fields.opts,
				maxVersion:       tt.fields.maxVersion,
				onDiskSize:       tt.fields.onDiskSize,
				staleDataSize:    tt.fields.staleDataSize,
				wg:               tt.fields.wg,
				blockChan:        tt.fields.blockChan,
				blockList:        tt.fields.blockList,
			}
			b.addHelper(tt.args.key, tt.args.v, tt.args.vpLen)
		})
	}
}

func TestBuilder_addInternal(t *testing.T) {
	type fields struct {
		alloc            *z.Allocator
		curBlock         *bblock
		compressedSize   atomic.Uint32
		uncompressedSize atomic.Uint32
		lenOffsets       uint32
		keyHashes        []uint32
		opts             *Options
		maxVersion       uint64
		onDiskSize       uint32
		staleDataSize    int
		wg               sync.WaitGroup
		blockChan        chan *bblock
		blockList        []*bblock
	}
	type args struct {
		key      []byte
		value    y.ValueStruct
		valueLen uint32
		isStale  bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				alloc:            tt.fields.alloc,
				curBlock:         tt.fields.curBlock,
				compressedSize:   tt.fields.compressedSize,
				uncompressedSize: tt.fields.uncompressedSize,
				lenOffsets:       tt.fields.lenOffsets,
				keyHashes:        tt.fields.keyHashes,
				opts:             tt.fields.opts,
				maxVersion:       tt.fields.maxVersion,
				onDiskSize:       tt.fields.onDiskSize,
				staleDataSize:    tt.fields.staleDataSize,
				wg:               tt.fields.wg,
				blockChan:        tt.fields.blockChan,
				blockList:        tt.fields.blockList,
			}
			b.addInternal(tt.args.key, tt.args.value, tt.args.valueLen, tt.args.isStale)
		})
	}
}

func TestBuilder_allocate(t *testing.T) {
	type fields struct {
		alloc            *z.Allocator
		curBlock         *bblock
		compressedSize   atomic.Uint32
		uncompressedSize atomic.Uint32
		lenOffsets       uint32
		keyHashes        []uint32
		opts             *Options
		maxVersion       uint64
		onDiskSize       uint32
		staleDataSize    int
		wg               sync.WaitGroup
		blockChan        chan *bblock
		blockList        []*bblock
	}
	type args struct {
		need int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				alloc:            tt.fields.alloc,
				curBlock:         tt.fields.curBlock,
				compressedSize:   tt.fields.compressedSize,
				uncompressedSize: tt.fields.uncompressedSize,
				lenOffsets:       tt.fields.lenOffsets,
				keyHashes:        tt.fields.keyHashes,
				opts:             tt.fields.opts,
				maxVersion:       tt.fields.maxVersion,
				onDiskSize:       tt.fields.onDiskSize,
				staleDataSize:    tt.fields.staleDataSize,
				wg:               tt.fields.wg,
				blockChan:        tt.fields.blockChan,
				blockList:        tt.fields.blockList,
			}
			if got := b.allocate(tt.args.need); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("allocate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_append(t *testing.T) {
	type fields struct {
		alloc            *z.Allocator
		curBlock         *bblock
		compressedSize   atomic.Uint32
		uncompressedSize atomic.Uint32
		lenOffsets       uint32
		keyHashes        []uint32
		opts             *Options
		maxVersion       uint64
		onDiskSize       uint32
		staleDataSize    int
		wg               sync.WaitGroup
		blockChan        chan *bblock
		blockList        []*bblock
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				alloc:            tt.fields.alloc,
				curBlock:         tt.fields.curBlock,
				compressedSize:   tt.fields.compressedSize,
				uncompressedSize: tt.fields.uncompressedSize,
				lenOffsets:       tt.fields.lenOffsets,
				keyHashes:        tt.fields.keyHashes,
				opts:             tt.fields.opts,
				maxVersion:       tt.fields.maxVersion,
				onDiskSize:       tt.fields.onDiskSize,
				staleDataSize:    tt.fields.staleDataSize,
				wg:               tt.fields.wg,
				blockChan:        tt.fields.blockChan,
				blockList:        tt.fields.blockList,
			}
			b.append(tt.args.data)
		})
	}
}

func TestBuilder_buildIndex(t *testing.T) {
	type fields struct {
		alloc            *z.Allocator
		curBlock         *bblock
		compressedSize   atomic.Uint32
		uncompressedSize atomic.Uint32
		lenOffsets       uint32
		keyHashes        []uint32
		opts             *Options
		maxVersion       uint64
		onDiskSize       uint32
		staleDataSize    int
		wg               sync.WaitGroup
		blockChan        chan *bblock
		blockList        []*bblock
	}
	type args struct {
		bloom []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
		want1  uint32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				alloc:            tt.fields.alloc,
				curBlock:         tt.fields.curBlock,
				compressedSize:   tt.fields.compressedSize,
				uncompressedSize: tt.fields.uncompressedSize,
				lenOffsets:       tt.fields.lenOffsets,
				keyHashes:        tt.fields.keyHashes,
				opts:             tt.fields.opts,
				maxVersion:       tt.fields.maxVersion,
				onDiskSize:       tt.fields.onDiskSize,
				staleDataSize:    tt.fields.staleDataSize,
				wg:               tt.fields.wg,
				blockChan:        tt.fields.blockChan,
				blockList:        tt.fields.blockList,
			}
			got, got1 := b.buildIndex(tt.args.bloom)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildIndex() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("buildIndex() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestBuilder_calculateCheckSum(t *testing.T) {
	type fields struct {
		alloc            *z.Allocator
		curBlock         *bblock
		compressedSize   atomic.Uint32
		uncompressedSize atomic.Uint32
		lenOffsets       uint32
		keyHashes        []uint32
		opts             *Options
		maxVersion       uint64
		onDiskSize       uint32
		staleDataSize    int
		wg               sync.WaitGroup
		blockChan        chan *bblock
		blockList        []*bblock
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				alloc:            tt.fields.alloc,
				curBlock:         tt.fields.curBlock,
				compressedSize:   tt.fields.compressedSize,
				uncompressedSize: tt.fields.uncompressedSize,
				lenOffsets:       tt.fields.lenOffsets,
				keyHashes:        tt.fields.keyHashes,
				opts:             tt.fields.opts,
				maxVersion:       tt.fields.maxVersion,
				onDiskSize:       tt.fields.onDiskSize,
				staleDataSize:    tt.fields.staleDataSize,
				wg:               tt.fields.wg,
				blockChan:        tt.fields.blockChan,
				blockList:        tt.fields.blockList,
			}
			if got := b.calculateCheckSum(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("calculateCheckSum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_checkBlockSizeGreaterThreshold(t *testing.T) {
	type fields struct {
		alloc            *z.Allocator
		curBlock         *bblock
		compressedSize   atomic.Uint32
		uncompressedSize atomic.Uint32
		lenOffsets       uint32
		keyHashes        []uint32
		opts             *Options
		maxVersion       uint64
		onDiskSize       uint32
		staleDataSize    int
		wg               sync.WaitGroup
		blockChan        chan *bblock
		blockList        []*bblock
	}
	type args struct {
		key   []byte
		value y.ValueStruct
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				alloc:            tt.fields.alloc,
				curBlock:         tt.fields.curBlock,
				compressedSize:   tt.fields.compressedSize,
				uncompressedSize: tt.fields.uncompressedSize,
				lenOffsets:       tt.fields.lenOffsets,
				keyHashes:        tt.fields.keyHashes,
				opts:             tt.fields.opts,
				maxVersion:       tt.fields.maxVersion,
				onDiskSize:       tt.fields.onDiskSize,
				staleDataSize:    tt.fields.staleDataSize,
				wg:               tt.fields.wg,
				blockChan:        tt.fields.blockChan,
				blockList:        tt.fields.blockList,
			}
			if got := b.checkBlockSizeGreaterThreshold(tt.args.key, tt.args.value); got != tt.want {
				t.Errorf("checkBlockSizeGreaterThreshold() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_compressData(t *testing.T) {
	type fields struct {
		alloc            *z.Allocator
		curBlock         *bblock
		compressedSize   atomic.Uint32
		uncompressedSize atomic.Uint32
		lenOffsets       uint32
		keyHashes        []uint32
		opts             *Options
		maxVersion       uint64
		onDiskSize       uint32
		staleDataSize    int
		wg               sync.WaitGroup
		blockChan        chan *bblock
		blockList        []*bblock
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				alloc:            tt.fields.alloc,
				curBlock:         tt.fields.curBlock,
				compressedSize:   tt.fields.compressedSize,
				uncompressedSize: tt.fields.uncompressedSize,
				lenOffsets:       tt.fields.lenOffsets,
				keyHashes:        tt.fields.keyHashes,
				opts:             tt.fields.opts,
				maxVersion:       tt.fields.maxVersion,
				onDiskSize:       tt.fields.onDiskSize,
				staleDataSize:    tt.fields.staleDataSize,
				wg:               tt.fields.wg,
				blockChan:        tt.fields.blockChan,
				blockList:        tt.fields.blockList,
			}
			got, err := b.compressData(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("compressData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("compressData() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_cutDownBlock(t *testing.T) {
	type fields struct {
		alloc            *z.Allocator
		curBlock         *bblock
		compressedSize   atomic.Uint32
		uncompressedSize atomic.Uint32
		lenOffsets       uint32
		keyHashes        []uint32
		opts             *Options
		maxVersion       uint64
		onDiskSize       uint32
		staleDataSize    int
		wg               sync.WaitGroup
		blockChan        chan *bblock
		blockList        []*bblock
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				alloc:            tt.fields.alloc,
				curBlock:         tt.fields.curBlock,
				compressedSize:   tt.fields.compressedSize,
				uncompressedSize: tt.fields.uncompressedSize,
				lenOffsets:       tt.fields.lenOffsets,
				keyHashes:        tt.fields.keyHashes,
				opts:             tt.fields.opts,
				maxVersion:       tt.fields.maxVersion,
				onDiskSize:       tt.fields.onDiskSize,
				staleDataSize:    tt.fields.staleDataSize,
				wg:               tt.fields.wg,
				blockChan:        tt.fields.blockChan,
				blockList:        tt.fields.blockList,
			}
			b.cutDownBlock()
		})
	}
}

func TestBuilder_encrypt(t *testing.T) {
	type fields struct {
		alloc            *z.Allocator
		curBlock         *bblock
		compressedSize   atomic.Uint32
		uncompressedSize atomic.Uint32
		lenOffsets       uint32
		keyHashes        []uint32
		opts             *Options
		maxVersion       uint64
		onDiskSize       uint32
		staleDataSize    int
		wg               sync.WaitGroup
		blockChan        chan *bblock
		blockList        []*bblock
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				alloc:            tt.fields.alloc,
				curBlock:         tt.fields.curBlock,
				compressedSize:   tt.fields.compressedSize,
				uncompressedSize: tt.fields.uncompressedSize,
				lenOffsets:       tt.fields.lenOffsets,
				keyHashes:        tt.fields.keyHashes,
				opts:             tt.fields.opts,
				maxVersion:       tt.fields.maxVersion,
				onDiskSize:       tt.fields.onDiskSize,
				staleDataSize:    tt.fields.staleDataSize,
				wg:               tt.fields.wg,
				blockChan:        tt.fields.blockChan,
				blockList:        tt.fields.blockList,
			}
			got, err := b.encrypt(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("encrypt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_handleBlock(t *testing.T) {
	type fields struct {
		alloc            *z.Allocator
		curBlock         *bblock
		compressedSize   atomic.Uint32
		uncompressedSize atomic.Uint32
		lenOffsets       uint32
		keyHashes        []uint32
		opts             *Options
		maxVersion       uint64
		onDiskSize       uint32
		staleDataSize    int
		wg               sync.WaitGroup
		blockChan        chan *bblock
		blockList        []*bblock
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				alloc:            tt.fields.alloc,
				curBlock:         tt.fields.curBlock,
				compressedSize:   tt.fields.compressedSize,
				uncompressedSize: tt.fields.uncompressedSize,
				lenOffsets:       tt.fields.lenOffsets,
				keyHashes:        tt.fields.keyHashes,
				opts:             tt.fields.opts,
				maxVersion:       tt.fields.maxVersion,
				onDiskSize:       tt.fields.onDiskSize,
				staleDataSize:    tt.fields.staleDataSize,
				wg:               tt.fields.wg,
				blockChan:        tt.fields.blockChan,
				blockList:        tt.fields.blockList,
			}
			b.handleBlock()
		})
	}
}

func TestBuilder_keyDiff(t *testing.T) {
	type fields struct {
		alloc            *z.Allocator
		curBlock         *bblock
		compressedSize   atomic.Uint32
		uncompressedSize atomic.Uint32
		lenOffsets       uint32
		keyHashes        []uint32
		opts             *Options
		maxVersion       uint64
		onDiskSize       uint32
		staleDataSize    int
		wg               sync.WaitGroup
		blockChan        chan *bblock
		blockList        []*bblock
	}
	type args struct {
		newKey []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{
			name: "test",
			fields: fields{
				curBlock: &bblock{
					baseKey: []byte{1, 2, 3},
				},
			},
			args: args{
				newKey: []byte{1, 2, 3},
			},
			want: []byte{},
		},
		{
			name: "test1",
			fields: fields{
				curBlock: &bblock{
					baseKey: []byte{1, 2, 3},
				},
			},
			args: args{
				newKey: []byte{1, 2, 3, 4},
			},
			want: []byte{4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				alloc:            tt.fields.alloc,
				curBlock:         tt.fields.curBlock,
				compressedSize:   tt.fields.compressedSize,
				uncompressedSize: tt.fields.uncompressedSize,
				lenOffsets:       tt.fields.lenOffsets,
				keyHashes:        tt.fields.keyHashes,
				opts:             tt.fields.opts,
				maxVersion:       tt.fields.maxVersion,
				onDiskSize:       tt.fields.onDiskSize,
				staleDataSize:    tt.fields.staleDataSize,
				wg:               tt.fields.wg,
				blockChan:        tt.fields.blockChan,
				blockList:        tt.fields.blockList,
			}
			if got := b.keyDiff(tt.args.newKey); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("keyDiff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_shouldEncrypt(t *testing.T) {
	type fields struct {
		alloc            *z.Allocator
		curBlock         *bblock
		compressedSize   atomic.Uint32
		uncompressedSize atomic.Uint32
		lenOffsets       uint32
		keyHashes        []uint32
		opts             *Options
		maxVersion       uint64
		onDiskSize       uint32
		staleDataSize    int
		wg               sync.WaitGroup
		blockChan        chan *bblock
		blockList        []*bblock
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
			b := &Builder{
				alloc:            tt.fields.alloc,
				curBlock:         tt.fields.curBlock,
				compressedSize:   tt.fields.compressedSize,
				uncompressedSize: tt.fields.uncompressedSize,
				lenOffsets:       tt.fields.lenOffsets,
				keyHashes:        tt.fields.keyHashes,
				opts:             tt.fields.opts,
				maxVersion:       tt.fields.maxVersion,
				onDiskSize:       tt.fields.onDiskSize,
				staleDataSize:    tt.fields.staleDataSize,
				wg:               tt.fields.wg,
				blockChan:        tt.fields.blockChan,
				blockList:        tt.fields.blockList,
			}
			if got := b.shouldEncrypt(); got != tt.want {
				t.Errorf("shouldEncrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_writeBlockOffset(t *testing.T) {
	type fields struct {
		alloc            *z.Allocator
		curBlock         *bblock
		compressedSize   atomic.Uint32
		uncompressedSize atomic.Uint32
		lenOffsets       uint32
		keyHashes        []uint32
		opts             *Options
		maxVersion       uint64
		onDiskSize       uint32
		staleDataSize    int
		wg               sync.WaitGroup
		blockChan        chan *bblock
		blockList        []*bblock
	}
	type args struct {
		builder     *flatbuffers.Builder
		bl          *bblock
		startOffset uint32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   flatbuffers.UOffsetT
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				alloc:            tt.fields.alloc,
				curBlock:         tt.fields.curBlock,
				compressedSize:   tt.fields.compressedSize,
				uncompressedSize: tt.fields.uncompressedSize,
				lenOffsets:       tt.fields.lenOffsets,
				keyHashes:        tt.fields.keyHashes,
				opts:             tt.fields.opts,
				maxVersion:       tt.fields.maxVersion,
				onDiskSize:       tt.fields.onDiskSize,
				staleDataSize:    tt.fields.staleDataSize,
				wg:               tt.fields.wg,
				blockChan:        tt.fields.blockChan,
				blockList:        tt.fields.blockList,
			}
			if got := b.writeBlockOffset(tt.args.builder, tt.args.bl, tt.args.startOffset); got != tt.want {
				t.Errorf("writeBlockOffset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_writeBlockOffsets(t *testing.T) {
	type fields struct {
		alloc            *z.Allocator
		curBlock         *bblock
		compressedSize   atomic.Uint32
		uncompressedSize atomic.Uint32
		lenOffsets       uint32
		keyHashes        []uint32
		opts             *Options
		maxVersion       uint64
		onDiskSize       uint32
		staleDataSize    int
		wg               sync.WaitGroup
		blockChan        chan *bblock
		blockList        []*bblock
	}
	type args struct {
		builder *flatbuffers.Builder
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []flatbuffers.UOffsetT
		want1  uint32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				alloc:            tt.fields.alloc,
				curBlock:         tt.fields.curBlock,
				compressedSize:   tt.fields.compressedSize,
				uncompressedSize: tt.fields.uncompressedSize,
				lenOffsets:       tt.fields.lenOffsets,
				keyHashes:        tt.fields.keyHashes,
				opts:             tt.fields.opts,
				maxVersion:       tt.fields.maxVersion,
				onDiskSize:       tt.fields.onDiskSize,
				staleDataSize:    tt.fields.staleDataSize,
				wg:               tt.fields.wg,
				blockChan:        tt.fields.blockChan,
				blockList:        tt.fields.blockList,
			}
			got, got1 := b.writeBlockOffsets(tt.args.builder)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("writeBlockOffsets() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("writeBlockOffsets() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNewTableBuilder(t *testing.T) {
	type args struct {
		opts Options
	}
	tests := []struct {
		name string
		args args
		want *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTableBuilder(tt.args.opts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTableBuilder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildData_Copy(t *testing.T) {
	type fields struct {
		blockList []*bblock
		index     []byte
		checksum  []byte
		Size      int
		alloc     *z.Allocator
	}
	type args struct {
		dst []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bd := &buildData{
				blockList: tt.fields.blockList,
				index:     tt.fields.index,
				checksum:  tt.fields.checksum,
				Size:      tt.fields.Size,
				alloc:     tt.fields.alloc,
			}
			if got := bd.Copy(tt.args.dst); got != tt.want {
				t.Errorf("Copy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_header_Encode(t *testing.T) {
	type fields struct {
		overlap uint16
		diff    uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := header{
				overlap: tt.fields.overlap,
				diff:    tt.fields.diff,
			}
			if got := h.Encode(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_maxEncodedLen(t *testing.T) {
	type args struct {
		ctype options.CompressionType
		sz    int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maxEncodedLen(tt.args.ctype, tt.args.sz); got != tt.want {
				t.Errorf("maxEncodedLen() = %v, want %v", got, tt.want)
			}
		})
	}
}
