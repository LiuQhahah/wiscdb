package internal

import (
	"hash"
	"io"
	"reflect"
	"testing"
)

func Test_hashReader_Read(t1 *testing.T) {
	type fields struct {
		r         io.Reader
		h         hash.Hash32
		bytesRead int
	}
	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "test",
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &hashReader{
				r:         tt.fields.r,
				h:         tt.fields.h,
				bytesRead: tt.fields.bytesRead,
			}
			got, err := t.Read(tt.args.p)
			if (err != nil) != tt.wantErr {
				t1.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t1.Errorf("Read() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hashReader_ReadByte(t1 *testing.T) {
	type fields struct {
		r         io.Reader
		h         hash.Hash32
		bytesRead int
	}
	tests := []struct {
		name    string
		fields  fields
		want    byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &hashReader{
				r:         tt.fields.r,
				h:         tt.fields.h,
				bytesRead: tt.fields.bytesRead,
			}
			got, err := t.ReadByte()
			if (err != nil) != tt.wantErr {
				t1.Errorf("ReadByte() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t1.Errorf("ReadByte() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hashReader_Sum32(t1 *testing.T) {
	type fields struct {
		r         io.Reader
		h         hash.Hash32
		bytesRead int
	}
	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &hashReader{
				r:         tt.fields.r,
				h:         tt.fields.h,
				bytesRead: tt.fields.bytesRead,
			}
			if got := t.Sum32(); got != tt.want {
				t1.Errorf("Sum32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newHashReader(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name string
		args args
		want *hashReader
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newHashReader(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newHashReader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_safeRead_Entry(t *testing.T) {
	type fields struct {
		k          []byte
		v          []byte
		readOffset uint32
		lf         *valueLogFile
	}
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Entry
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &safeRead{
				k:            tt.fields.k,
				v:            tt.fields.v,
				recordOffset: tt.fields.readOffset,
				vLogFile:     tt.fields.lf,
			}
			got, err := r.Entry(tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("Entry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Entry() got = %v, want %v", got, tt.want)
			}
		})
	}
}
