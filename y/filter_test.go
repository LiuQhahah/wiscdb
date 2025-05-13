package y

import (
	"reflect"
	"testing"
)

func TestBloomBitsPerKey(t *testing.T) {
	type args struct {
		numEntries int
		fp         float64
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
			if got := BloomBitsPerKey(tt.args.numEntries, tt.args.fp); got != tt.want {
				t.Errorf("BloomBitsPerKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBloomBitsPerKey1(t *testing.T) {
	type args struct {
		numEntries int
		fp         float64
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
			if got := BloomBitsPerKey(tt.args.numEntries, tt.args.fp); got != tt.want {
				t.Errorf("BloomBitsPerKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilter_MayContain(t *testing.T) {
	type args struct {
		h uint32
	}
	tests := []struct {
		name string
		f    Filter
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.MayContain(tt.args.h); got != tt.want {
				t.Errorf("MayContain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilter_MayContainKey(t *testing.T) {
	type args struct {
		key []byte
	}
	tests := []struct {
		name string
		f    Filter
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.MayContainKey(tt.args.key); got != tt.want {
				t.Errorf("MayContainKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHash(t *testing.T) {
	type args struct {
		key []byte
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Hash(tt.args.key); got != tt.want {
				t.Errorf("Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewFilter(t *testing.T) {
	type args struct {
		keys       []uint32
		bitsPerKey int
	}
	tests := []struct {
		name string
		args args
		want Filter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFilter(tt.args.keys, tt.args.bitsPerKey); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_appendFilter(t *testing.T) {
	type args struct {
		buf        []byte
		keys       []uint32
		bitsPerKey int
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := appendFilter(tt.args.buf, tt.args.keys, tt.args.bitsPerKey); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appendFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extend(t *testing.T) {
	type args struct {
		b []byte
		n int
	}
	tests := []struct {
		name        string
		args        args
		wantOverall []byte
		wantTrailer []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOverall, gotTrailer := extend(tt.args.b, tt.args.n)
			if !reflect.DeepEqual(gotOverall, tt.wantOverall) {
				t.Errorf("extend() gotOverall = %v, want %v", gotOverall, tt.wantOverall)
			}
			if !reflect.DeepEqual(gotTrailer, tt.wantTrailer) {
				t.Errorf("extend() gotTrailer = %v, want %v", gotTrailer, tt.wantTrailer)
			}
		})
	}
}
