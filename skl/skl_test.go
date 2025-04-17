package skl

import (
	"reflect"
	"sync/atomic"
	"testing"
	"wiscdb/y"
)

func TestArena_getKey(t *testing.T) {
	type fields struct {
		n   atomic.Uint32
		buf []byte
	}
	type args struct {
		offset uint32
		size   uint16
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
				buf: []byte{1, 2, 3, 4},
			},
			args: args{
				offset: 1,
				size:   2,
			},
			want: []byte{2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Arena{
				n:   tt.fields.n,
				buf: tt.fields.buf,
			}
			if got := s.getKey(tt.args.offset, tt.args.size); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArena_getNode(t *testing.T) {
	type fields struct {
		n   atomic.Uint32
		buf []byte
	}
	var value1 atomic.Uint32
	value1.Store(100)
	var value2 atomic.Uint64
	value2.Store(2)

	type args struct {
		offset uint32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *node
	}{
		{
			name: "test",
			fields: fields{
				n:   value1,
				buf: []byte{1, 1, 11, 1, 1, 1, 11, 1},
			},
			args: args{
				offset: 1,
			},
			want: &node{
				value: value2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Arena{
				n:   tt.fields.n,
				buf: tt.fields.buf,
			}
			if got := s.getNode(tt.args.offset); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArena_getNodeOffset(t *testing.T) {
	type fields struct {
		n   atomic.Uint32
		buf []byte
	}
	type args struct {
		nd *node
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   uint32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Arena{
				n:   tt.fields.n,
				buf: tt.fields.buf,
			}
			if got := s.getNodeOffset(tt.args.nd); got != tt.want {
				t.Errorf("getNodeOffset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArena_getVal(t *testing.T) {
	type fields struct {
		n   atomic.Uint32
		buf []byte
	}
	type args struct {
		offset uint32
		size   uint32
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRet y.ValueStruct
	}{
		{
			name: "test",
			fields: fields{
				buf: []byte{0, 0, 1, 2, 3, 4, 5},
			},
			args: args{
				offset: 2,
				size:   4,
			},
			wantRet: y.ValueStruct{
				Meta:      1,
				UserMeta:  2,
				ExpiresAt: 3,
				Value:     []byte{4},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Arena{
				n:   tt.fields.n,
				buf: tt.fields.buf,
			}
			if gotRet := s.getVal(tt.args.offset, tt.args.size); !reflect.DeepEqual(gotRet, tt.wantRet) {
				t.Errorf("getVal() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}

func TestArena_putKey(t *testing.T) {
	type fields struct {
		n   atomic.Uint32
		buf []byte
	}
	type args struct {
		key []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   uint32
	}{
		{
			name: "test",
			fields: fields{
				buf: make([]byte, 10),
			},
			args: args{
				key: []byte{129},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Arena{
				n:   tt.fields.n,
				buf: tt.fields.buf,
			}
			if got := s.putKey(tt.args.key); got != tt.want {
				t.Errorf("putKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArena_putNode(t *testing.T) {

	var value atomic.Uint32
	value.Store(100) // Initializing the value (optional, default is 0)
	var value2 atomic.Uint32
	value2.Store(1000) // Initializing the value (optional, default is 0)
	type fields struct {
		n   atomic.Uint32
		buf []byte
	}
	type args struct {
		height int
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   uint32
	}{
		{
			name: "test",
			fields: fields{
				buf: []byte{0, 0, 0},
			},
			args: args{
				height: 1,
			},
			want: 0,
		},
		{
			name: "test",
			fields: fields{
				n:   value,
				buf: []byte{0, 0, 0},
			},
			args: args{
				height: 1,
			},
			want: 104,
		},

		{
			name: "test2",
			fields: fields{
				n:   value2,
				buf: make([]byte, 1000),
			},
			args: args{
				maxHeight,
			},
			want: 1000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Arena{
				n:   tt.fields.n,
				buf: tt.fields.buf,
			}
			if got := s.putNode(tt.args.height); got != tt.want {
				t.Errorf("putNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArena_putVal(t *testing.T) {
	type fields struct {
		n   atomic.Uint32
		buf []byte
	}
	type args struct {
		v y.ValueStruct
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   uint32
	}{
		{
			name: "test",
			fields: fields{
				buf: make([]byte, 10),
			},
			args: args{
				v: y.ValueStruct{
					Value: []byte{1, 1, 1, 1, 11, 1},
				},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Arena{
				n:   tt.fields.n,
				buf: tt.fields.buf,
			}
			if got := s.putVal(tt.args.v); got != tt.want {
				t.Errorf("putVal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArena_size(t *testing.T) {
	type fields struct {
		n   atomic.Uint32
		buf []byte
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
			s := &Arena{
				n:   tt.fields.n,
				buf: tt.fields.buf,
			}
			if got := s.size(); got != tt.want {
				t.Errorf("size() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterator_Close(t *testing.T) {
	type fields struct {
		list *SkipList
		n    *node
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
			i := &Iterator{
				list: tt.fields.list,
				n:    tt.fields.n,
			}
			if err := i.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIterator_Key(t *testing.T) {
	type fields struct {
		list *SkipList
		n    *node
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
			i := &Iterator{
				list: tt.fields.list,
				n:    tt.fields.n,
			}
			if got := i.Key(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Key() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterator_NewIterator(t *testing.T) {
	type fields struct {
		list *SkipList
		n    *node
	}
	tests := []struct {
		name   string
		fields fields
		want   *Iterator
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Iterator{
				list: tt.fields.list,
				n:    tt.fields.n,
			}
			if got := i.NewIterator(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIterator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterator_Next(t *testing.T) {
	type fields struct {
		list *SkipList
		n    *node
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Iterator{
				list: tt.fields.list,
				n:    tt.fields.n,
			}
			i.Next()
		})
	}
}

func TestIterator_Prev(t *testing.T) {
	type fields struct {
		list *SkipList
		n    *node
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Iterator{
				list: tt.fields.list,
				n:    tt.fields.n,
			}
			i.Prev()
		})
	}
}

func TestIterator_Seek(t *testing.T) {
	type fields struct {
		list *SkipList
		n    *node
	}
	type args struct {
		target []byte
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
			i := &Iterator{
				list: tt.fields.list,
				n:    tt.fields.n,
			}
			i.Seek(tt.args.target)
		})
	}
}

func TestIterator_SeekForPrev(t *testing.T) {
	type fields struct {
		list *SkipList
		n    *node
	}
	type args struct {
		target []byte
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
			i := &Iterator{
				list: tt.fields.list,
				n:    tt.fields.n,
			}
			i.SeekForPrev(tt.args.target)
		})
	}
}

func TestIterator_SeekToFirst(t *testing.T) {
	type fields struct {
		list *SkipList
		n    *node
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Iterator{
				list: tt.fields.list,
				n:    tt.fields.n,
			}
			i.SeekToFirst()
		})
	}
}

func TestIterator_SeekToLast(t *testing.T) {
	type fields struct {
		list *SkipList
		n    *node
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Iterator{
				list: tt.fields.list,
				n:    tt.fields.n,
			}
			i.SeekToLast()
		})
	}
}

func TestIterator_Valid(t *testing.T) {
	type fields struct {
		list *SkipList
		n    *node
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
			i := &Iterator{
				list: tt.fields.list,
				n:    tt.fields.n,
			}
			if got := i.Valid(); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterator_Value(t *testing.T) {
	type fields struct {
		list *SkipList
		n    *node
	}
	tests := []struct {
		name   string
		fields fields
		want   y.ValueStruct
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Iterator{
				list: tt.fields.list,
				n:    tt.fields.n,
			}
			if got := i.Value(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterator_ValueUint64(t *testing.T) {
	type fields struct {
		list *SkipList
		n    *node
	}
	tests := []struct {
		name   string
		fields fields
		want   uint64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Iterator{
				list: tt.fields.list,
				n:    tt.fields.n,
			}
			if got := i.ValueUint64(); got != tt.want {
				t.Errorf("ValueUint64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSkipList(t *testing.T) {
	type args struct {
		arenaSize int64
	}
	tests := []struct {
		name string
		args args
		want *SkipList
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSkipList(tt.args.arenaSize); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSkipList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSkipList_DecrRef(t *testing.T) {
	type fields struct {
		height  atomic.Int32
		head    *node
		ref     atomic.Int32
		arena   *Arena
		OnClose func()
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SkipList{
				height:  tt.fields.height,
				head:    tt.fields.head,
				ref:     tt.fields.ref,
				arena:   tt.fields.arena,
				OnClose: tt.fields.OnClose,
			}
			s.DecrRef()
		})
	}
}

func TestSkipList_Empty(t *testing.T) {
	type fields struct {
		height  atomic.Int32
		head    *node
		ref     atomic.Int32
		arena   *Arena
		OnClose func()
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
			s := &SkipList{
				height:  tt.fields.height,
				head:    tt.fields.head,
				ref:     tt.fields.ref,
				arena:   tt.fields.arena,
				OnClose: tt.fields.OnClose,
			}
			if got := s.Empty(); got != tt.want {
				t.Errorf("Empty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSkipList_Get(t *testing.T) {
	type fields struct {
		height  atomic.Int32
		head    *node
		ref     atomic.Int32
		arena   *Arena
		OnClose func()
	}
	type args struct {
		key []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   y.ValueStruct
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SkipList{
				height:  tt.fields.height,
				head:    tt.fields.head,
				ref:     tt.fields.ref,
				arena:   tt.fields.arena,
				OnClose: tt.fields.OnClose,
			}
			if got := s.Get(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSkipList_IncrRef(t *testing.T) {
	type fields struct {
		height  atomic.Int32
		head    *node
		ref     atomic.Int32
		arena   *Arena
		OnClose func()
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SkipList{
				height:  tt.fields.height,
				head:    tt.fields.head,
				ref:     tt.fields.ref,
				arena:   tt.fields.arena,
				OnClose: tt.fields.OnClose,
			}
			s.IncrRef()
		})
	}
}

func TestSkipList_MemSize(t *testing.T) {
	type fields struct {
		height  atomic.Int32
		head    *node
		ref     atomic.Int32
		arena   *Arena
		OnClose func()
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
			s := &SkipList{
				height:  tt.fields.height,
				head:    tt.fields.head,
				ref:     tt.fields.ref,
				arena:   tt.fields.arena,
				OnClose: tt.fields.OnClose,
			}
			if got := s.MemSize(); got != tt.want {
				t.Errorf("MemSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSkipList_NewUniIterator(t *testing.T) {
	type fields struct {
		height  atomic.Int32
		head    *node
		ref     atomic.Int32
		arena   *Arena
		OnClose func()
	}
	type args struct {
		reversed bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *UniIterator
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SkipList{
				height:  tt.fields.height,
				head:    tt.fields.head,
				ref:     tt.fields.ref,
				arena:   tt.fields.arena,
				OnClose: tt.fields.OnClose,
			}
			if got := s.NewUniIterator(tt.args.reversed); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUniIterator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSkipList_Put(t *testing.T) {
	type fields struct {
		height  atomic.Int32
		head    *node
		ref     atomic.Int32
		arena   *Arena
		OnClose func()
	}
	type args struct {
		key []byte
		v   y.ValueStruct
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
			s := &SkipList{
				height:  tt.fields.height,
				head:    tt.fields.head,
				ref:     tt.fields.ref,
				arena:   tt.fields.arena,
				OnClose: tt.fields.OnClose,
			}
			s.Put(tt.args.key, tt.args.v)
		})
	}
}

func TestSkipList_findLast(t *testing.T) {
	type fields struct {
		height  atomic.Int32
		head    *node
		ref     atomic.Int32
		arena   *Arena
		OnClose func()
	}
	tests := []struct {
		name   string
		fields fields
		want   *node
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SkipList{
				height:  tt.fields.height,
				head:    tt.fields.head,
				ref:     tt.fields.ref,
				arena:   tt.fields.arena,
				OnClose: tt.fields.OnClose,
			}
			if got := s.findLast(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findLast() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSkipList_findNear(t *testing.T) {
	type fields struct {
		height  atomic.Int32
		head    *node
		ref     atomic.Int32
		arena   *Arena
		OnClose func()
	}
	type args struct {
		key        []byte
		less       bool
		allowEqual bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *node
		want1  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SkipList{
				height:  tt.fields.height,
				head:    tt.fields.head,
				ref:     tt.fields.ref,
				arena:   tt.fields.arena,
				OnClose: tt.fields.OnClose,
			}
			got, got1 := s.findNear(tt.args.key, tt.args.less, tt.args.allowEqual)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findNear() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("findNear() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestSkipList_findSplitForLevel(t *testing.T) {
	type fields struct {
		height  atomic.Int32
		head    *node
		ref     atomic.Int32
		arena   *Arena
		OnClose func()
	}
	type args struct {
		key    []byte
		before *node
		less   int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *node
		want1  *node
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SkipList{
				height:  tt.fields.height,
				head:    tt.fields.head,
				ref:     tt.fields.ref,
				arena:   tt.fields.arena,
				OnClose: tt.fields.OnClose,
			}
			got, got1 := s.findSpliceForLevel(tt.args.key, tt.args.before, tt.args.less)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findSpliceForLevel() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("findSpliceForLevel() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestSkipList_getHeight(t *testing.T) {
	type fields struct {
		height  atomic.Int32
		head    *node
		ref     atomic.Int32
		arena   *Arena
		OnClose func()
	}
	tests := []struct {
		name   string
		fields fields
		want   int32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SkipList{
				height:  tt.fields.height,
				head:    tt.fields.head,
				ref:     tt.fields.ref,
				arena:   tt.fields.arena,
				OnClose: tt.fields.OnClose,
			}
			if got := s.getHeight(); got != tt.want {
				t.Errorf("getHeight() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSkipList_getNext(t *testing.T) {
	type fields struct {
		height  atomic.Int32
		head    *node
		ref     atomic.Int32
		arena   *Arena
		OnClose func()
	}
	type args struct {
		nd     *node
		height int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *node
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SkipList{
				height:  tt.fields.height,
				head:    tt.fields.head,
				ref:     tt.fields.ref,
				arena:   tt.fields.arena,
				OnClose: tt.fields.OnClose,
			}
			if got := s.getNext(tt.args.nd, tt.args.height); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getNext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSkipList_randomHeight(t *testing.T) {
	type fields struct {
		height  atomic.Int32
		head    *node
		ref     atomic.Int32
		arena   *Arena
		OnClose func()
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SkipList{
				height:  tt.fields.height,
				head:    tt.fields.head,
				ref:     tt.fields.ref,
				arena:   tt.fields.arena,
				OnClose: tt.fields.OnClose,
			}
			if got := s.randomHeight(); got != tt.want {
				t.Errorf("randomHeight() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUniIterator_Close(t *testing.T) {
	type fields struct {
		iter     *Iterator
		reversed bool
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
			s := &UniIterator{
				iter:     tt.fields.iter,
				reversed: tt.fields.reversed,
			}
			if err := s.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUniIterator_Key(t *testing.T) {
	type fields struct {
		iter     *Iterator
		reversed bool
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
			s := &UniIterator{
				iter:     tt.fields.iter,
				reversed: tt.fields.reversed,
			}
			if got := s.Key(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Key() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUniIterator_Next(t *testing.T) {
	type fields struct {
		iter     *Iterator
		reversed bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &UniIterator{
				iter:     tt.fields.iter,
				reversed: tt.fields.reversed,
			}
			s.Next()
		})
	}
}

func TestUniIterator_ReWind(t *testing.T) {
	type fields struct {
		iter     *Iterator
		reversed bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &UniIterator{
				iter:     tt.fields.iter,
				reversed: tt.fields.reversed,
			}
			s.ReWind()
		})
	}
}

func TestUniIterator_Seek(t *testing.T) {
	type fields struct {
		iter     *Iterator
		reversed bool
	}
	type args struct {
		key []byte
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
			s := &UniIterator{
				iter:     tt.fields.iter,
				reversed: tt.fields.reversed,
			}
			s.Seek(tt.args.key)
		})
	}
}

func TestUniIterator_Valid(t *testing.T) {
	type fields struct {
		iter     *Iterator
		reversed bool
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
			s := &UniIterator{
				iter:     tt.fields.iter,
				reversed: tt.fields.reversed,
			}
			if got := s.Valid(); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUniIterator_Value(t *testing.T) {
	type fields struct {
		iter     *Iterator
		reversed bool
	}
	tests := []struct {
		name   string
		fields fields
		want   y.ValueStruct
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &UniIterator{
				iter:     tt.fields.iter,
				reversed: tt.fields.reversed,
			}
			if got := s.Value(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_decodeValue(t *testing.T) {
	type args struct {
		value uint64
	}
	tests := []struct {
		name          string
		args          args
		wantValOffset uint32
		wantValSize   uint32
	}{
		{
			name: "test",
			args: args{
				value: 0x01234567ABCDEFAB,
			},
			wantValOffset: 0x01234567,
			wantValSize:   0xABCDEFAB,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValOffset, gotValSize := decodeValue(tt.args.value)
			if gotValOffset != tt.wantValOffset {
				t.Errorf("decodeValue() gotValOffset = %v, want %v", gotValOffset, tt.wantValOffset)
			}
			if gotValSize != tt.wantValSize {
				t.Errorf("decodeValue() gotValSize = %v, want %v", gotValSize, tt.wantValSize)
			}
		})
	}
}

func Test_encodeValue(t *testing.T) {
	type args struct {
		valOffset uint32
		valSize   uint32
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{
			name: "test",
			args: args{
				valOffset: 0x01234567,
				valSize:   0xABCDEFAB,
			},
			want: 0x01234567ABCDEFAB,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := encodeValue(tt.args.valOffset, tt.args.valSize); got != tt.want {
				t.Errorf("encodeValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newArena(t *testing.T) {
	type args struct {
		n int64
	}
	nn := atomic.Uint32{}
	nn.Store(1)
	tests := []struct {
		name string
		args args
		want *Arena
	}{
		{
			name: "Test",
			args: args{
				n: 64 << 10,
			},
			want: &Arena{
				n:   nn,
				buf: make([]byte, 64<<10),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newArena(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newArena() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newNode(t *testing.T) {
	type args struct {
		arena  *Arena
		key    []byte
		v      y.ValueStruct
		height int
	}
	tests := []struct {
		name string
		args args
		want *node
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newNode(tt.args.arena, tt.args.key, tt.args.v, tt.args.height); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_node_casNextOffset(t *testing.T) {
	type fields struct {
		value     atomic.Uint64
		ketOffset uint32
		keySize   uint16
		height    uint16
		tower     [18]atomic.Uint32
	}
	type args struct {
		h   int
		old uint32
		val uint32
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
			s := &node{
				value:     tt.fields.value,
				keyOffset: tt.fields.ketOffset,
				keySize:   tt.fields.keySize,
				height:    tt.fields.height,
				tower:     tt.fields.tower,
			}
			if got := s.casNextOffset(tt.args.h, tt.args.old, tt.args.val); got != tt.want {
				t.Errorf("casNextOffset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_node_getNextOffset(t *testing.T) {
	type fields struct {
		value     atomic.Uint64
		ketOffset uint32
		keySize   uint16
		height    uint16
		tower     [18]atomic.Uint32
	}
	type args struct {
		h int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   uint32
	}{
		{
			name: "test",
			fields: fields{
				value:     atomic.Uint64{},
				ketOffset: 0,
				keySize:   0,
				height:    0,
				tower:     [18]atomic.Uint32{},
			},
			args: args{
				h: 10,
			},
			want: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.value.Store(10)
			tt.fields.tower[10].Store(100)
			s := &node{
				value:     tt.fields.value,
				keyOffset: tt.fields.ketOffset,
				keySize:   tt.fields.keySize,
				height:    tt.fields.height,
				tower:     tt.fields.tower,
			}
			if got := s.getNextOffset(tt.args.h); got != tt.want {
				t.Errorf("getNextOffset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_node_setValue(t *testing.T) {
	type fields struct {
		value     atomic.Uint64
		ketOffset uint32
		keySize   uint16
		height    uint16
		tower     [18]atomic.Uint32
	}
	type args struct {
		arena *Arena
		v     y.ValueStruct
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
			s := &node{
				value:     tt.fields.value,
				keyOffset: tt.fields.ketOffset,
				keySize:   tt.fields.keySize,
				height:    tt.fields.height,
				tower:     tt.fields.tower,
			}
			s.setValue(tt.args.arena, tt.args.v)
		})
	}
}

func Test_height_Swap(t *testing.T) {
	height := 10
	listHeight := 8
	s := SkipList{}
	s.height.Store(8)
	for height > listHeight {
		//更新下当前跳表的属性：高度
		//如果listHeight不等于s.height就会执行失败，此时就需要调用listHeight = int(s.getHeight()) 刷新一下
		if s.height.CompareAndSwap(int32(listHeight), int32(height)) {
			break
		}
		//重新获取跳表高度
		listHeight = int(s.getHeight())
	}
	//if height > listHeight {
	//	s.height.Store(int32(height))
	//}

}
