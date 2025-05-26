package main

import (
	"reflect"
	"testing"
)

func Test_valuePointer_Decode(t *testing.T) {
	type fields struct {
		Fid    uint32
		Len    uint32
		Offset uint32
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "test",
			fields: fields{
				Fid:    0x1,
				Len:    0x10,
				Offset: 0xAA,
			},
			args: args{
				b: []byte{0x01, 0x10, 0xAA, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &valuePointer{
				Fid:    tt.fields.Fid,
				Len:    tt.fields.Len,
				Offset: tt.fields.Offset,
			}
			p.Decode(tt.args.b)
			if !reflect.DeepEqual(p.Fid, tt.fields.Fid) {
				t.Errorf("valuePointer.Decode() = %v, want %v", p.Fid, tt.fields.Fid)
			}
			if !reflect.DeepEqual(p.Len, tt.fields.Len) {
				t.Errorf("valuePointer.Decode() = %v, want %v", p.Len, tt.fields.Len)
			}
			if !reflect.DeepEqual(p.Offset, tt.fields.Offset) {
				t.Errorf("valuePointer.Decode() = %v, want %v", p.Offset, tt.fields.Offset)
			}
		})
	}
}

func Test_valuePointer_Encode(t *testing.T) {
	type fields struct {
		Fid    uint32
		Len    uint32
		Offset uint32
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "test",
			fields: fields{
				Fid:    0x1,
				Len:    0x10,
				Offset: 0xAA,
			},
			want: []byte{0x01, 0x10, 0xAA, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := valuePointer{
				Fid:    tt.fields.Fid,
				Len:    tt.fields.Len,
				Offset: tt.fields.Offset,
			}
			if got := p.Encode(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_valuePointer_IsZero(t *testing.T) {
	type fields struct {
		Fid    uint32
		Len    uint32
		Offset uint32
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
			p := valuePointer{
				Fid:    tt.fields.Fid,
				Len:    tt.fields.Len,
				Offset: tt.fields.Offset,
			}
			if got := p.IsZero(); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_valuePointer_Less(t *testing.T) {
	type fields struct {
		Fid    uint32
		Len    uint32
		Offset uint32
	}
	type args struct {
		o valuePointer
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
			p := valuePointer{
				Fid:    tt.fields.Fid,
				Len:    tt.fields.Len,
				Offset: tt.fields.Offset,
			}
			if got := p.Less(tt.args.o); got != tt.want {
				t.Errorf("Less() = %v, want %v", got, tt.want)
			}
		})
	}
}
