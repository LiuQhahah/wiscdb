package internal

import (
	"reflect"
	"testing"
	"time"
)

func TestEntry_WithDiscard(t *testing.T) {
	type fields struct {
		Key          []byte
		Value        []byte
		ExpiresAt    uint64
		version      uint64
		offset       uint32
		UserMeta     byte
		meta         byte
		hlen         int
		valThreshold int64
	}
	tests := []struct {
		name   string
		fields fields
		want   *Entry
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Entry{
				Key:          tt.fields.Key,
				Value:        tt.fields.Value,
				ExpiresAt:    tt.fields.ExpiresAt,
				version:      tt.fields.version,
				offset:       tt.fields.offset,
				UserMeta:     tt.fields.UserMeta,
				meta:         tt.fields.meta,
				hlen:         tt.fields.hlen,
				valThreshold: tt.fields.valThreshold,
			}
			if got := e.WithDiscard(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDiscard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntry_WithTTL(t *testing.T) {
	type fields struct {
		Key          []byte
		Value        []byte
		ExpiresAt    uint64
		version      uint64
		offset       uint32
		UserMeta     byte
		meta         byte
		hlen         int
		valThreshold int64
	}
	type args struct {
		dur time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Entry
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Entry{
				Key:          tt.fields.Key,
				Value:        tt.fields.Value,
				ExpiresAt:    tt.fields.ExpiresAt,
				version:      tt.fields.version,
				offset:       tt.fields.offset,
				UserMeta:     tt.fields.UserMeta,
				meta:         tt.fields.meta,
				hlen:         tt.fields.hlen,
				valThreshold: tt.fields.valThreshold,
			}
			if got := e.WithTTL(tt.args.dur); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithTTL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntry_WithUserMeta(t *testing.T) {
	type fields struct {
		Key          []byte
		Value        []byte
		ExpiresAt    uint64
		version      uint64
		offset       uint32
		UserMeta     byte
		meta         byte
		hlen         int
		valThreshold int64
	}
	type args struct {
		meta byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Entry
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Entry{
				Key:          tt.fields.Key,
				Value:        tt.fields.Value,
				ExpiresAt:    tt.fields.ExpiresAt,
				version:      tt.fields.version,
				offset:       tt.fields.offset,
				UserMeta:     tt.fields.UserMeta,
				meta:         tt.fields.meta,
				hlen:         tt.fields.hlen,
				valThreshold: tt.fields.valThreshold,
			}
			if got := e.WithUserMeta(tt.args.meta); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithUserMeta() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntry_checkValueThreshold(t *testing.T) {
	type fields struct {
		Key          []byte
		Value        []byte
		ExpiresAt    uint64
		version      uint64
		offset       uint32
		UserMeta     byte
		meta         byte
		hlen         int
		valThreshold int64
	}
	type args struct {
		threshold int64
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
			e := &Entry{
				Key:          tt.fields.Key,
				Value:        tt.fields.Value,
				ExpiresAt:    tt.fields.ExpiresAt,
				version:      tt.fields.version,
				offset:       tt.fields.offset,
				UserMeta:     tt.fields.UserMeta,
				meta:         tt.fields.meta,
				hlen:         tt.fields.hlen,
				valThreshold: tt.fields.valThreshold,
			}
			if got := e.checkValueThreshold(tt.args.threshold); got != tt.want {
				t.Errorf("checkValueThreshold() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntry_estimateSizeAndSetThreshold(t *testing.T) {
	type fields struct {
		Key          []byte
		Value        []byte
		ExpiresAt    uint64
		version      uint64
		offset       uint32
		UserMeta     byte
		meta         byte
		hlen         int
		valThreshold int64
	}
	type args struct {
		threshold int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Entry{
				Key:          tt.fields.Key,
				Value:        tt.fields.Value,
				ExpiresAt:    tt.fields.ExpiresAt,
				version:      tt.fields.version,
				offset:       tt.fields.offset,
				UserMeta:     tt.fields.UserMeta,
				meta:         tt.fields.meta,
				hlen:         tt.fields.hlen,
				valThreshold: tt.fields.valThreshold,
			}
			if got := e.estimateSizeAndSetThreshold(tt.args.threshold); got != tt.want {
				t.Errorf("estimateSizeAndSetThreshold() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntry_isEmpty(t *testing.T) {
	type fields struct {
		Key          []byte
		Value        []byte
		ExpiresAt    uint64
		version      uint64
		offset       uint32
		UserMeta     byte
		meta         byte
		hlen         int
		valThreshold int64
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
			e := &Entry{
				Key:          tt.fields.Key,
				Value:        tt.fields.Value,
				ExpiresAt:    tt.fields.ExpiresAt,
				version:      tt.fields.version,
				offset:       tt.fields.offset,
				UserMeta:     tt.fields.UserMeta,
				meta:         tt.fields.meta,
				hlen:         tt.fields.hlen,
				valThreshold: tt.fields.valThreshold,
			}
			if got := e.isEmpty(); got != tt.want {
				t.Errorf("isEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntry_print(t *testing.T) {
	type fields struct {
		Key          []byte
		Value        []byte
		ExpiresAt    uint64
		version      uint64
		offset       uint32
		UserMeta     byte
		meta         byte
		hlen         int
		valThreshold int64
	}
	type args struct {
		prefix string
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
			e := Entry{
				Key:          tt.fields.Key,
				Value:        tt.fields.Value,
				ExpiresAt:    tt.fields.ExpiresAt,
				version:      tt.fields.version,
				offset:       tt.fields.offset,
				UserMeta:     tt.fields.UserMeta,
				meta:         tt.fields.meta,
				hlen:         tt.fields.hlen,
				valThreshold: tt.fields.valThreshold,
			}
			e.print(tt.args.prefix)
		})
	}
}

func TestEntry_skipVlogAndSetThreshold(t *testing.T) {
	type fields struct {
		Key          []byte
		Value        []byte
		ExpiresAt    uint64
		version      uint64
		offset       uint32
		UserMeta     byte
		meta         byte
		hlen         int
		valThreshold int64
	}
	type args struct {
		threshold int64
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
			e := &Entry{
				Key:          tt.fields.Key,
				Value:        tt.fields.Value,
				ExpiresAt:    tt.fields.ExpiresAt,
				version:      tt.fields.version,
				offset:       tt.fields.offset,
				UserMeta:     tt.fields.UserMeta,
				meta:         tt.fields.meta,
				hlen:         tt.fields.hlen,
				valThreshold: tt.fields.valThreshold,
			}
			if got := e.skipVlogAndSetThreshold(tt.args.threshold); got != tt.want {
				t.Errorf("skipVlogAndSetThreshold() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntry_withMergeBit(t *testing.T) {
	type fields struct {
		Key          []byte
		Value        []byte
		ExpiresAt    uint64
		version      uint64
		offset       uint32
		UserMeta     byte
		meta         byte
		hlen         int
		valThreshold int64
	}
	tests := []struct {
		name   string
		fields fields
		want   *Entry
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Entry{
				Key:          tt.fields.Key,
				Value:        tt.fields.Value,
				ExpiresAt:    tt.fields.ExpiresAt,
				version:      tt.fields.version,
				offset:       tt.fields.offset,
				UserMeta:     tt.fields.UserMeta,
				meta:         tt.fields.meta,
				hlen:         tt.fields.hlen,
				valThreshold: tt.fields.valThreshold,
			}
			if got := e.withMergeBit(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("withMergeBit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewEntry(t *testing.T) {
	type args struct {
		key   []byte
		value []byte
	}
	tests := []struct {
		name string
		args args
		want *Entry
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEntry(tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEntry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestXORBlock(t *testing.T) {
	type args struct {
		dst []byte
		src []byte
		key []byte
		iv  []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				dst: []byte{0x91, 0xaa},
				src: []byte{0x91, 0xaa},
				key: []byte{0x91, 0xaa},
				iv:  []byte{0x91, 0xaa},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := XORBlock(tt.args.dst, tt.args.src, tt.args.key, tt.args.iv); (err != nil) != tt.wantErr {
				t.Errorf("XORBlock() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
