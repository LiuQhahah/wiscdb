package y

import (
	"bytes"
	"reflect"
	"testing"
)

func TestXORBlockAllocate(t *testing.T) {
	type args struct {
		src []byte
		key []byte
		iv  []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				src: []byte{0x01, 0x02, 0x03},
				key: []byte{0x01, 0x02, 0x01, 0x02, 0x01, 0x02, 0x01, 0x02, 0x01, 0x02, 0x01, 0x02, 0x01, 0x02, 0x01, 0x02},
				iv:  []byte{0x01, 0x02, 0x01, 0x02, 0x01, 0x02, 0x01, 0x02, 0x01, 0x02, 0x01, 0x02, 0x01, 0x02, 0x01, 0x02},
			},
			want:    []byte{0x84, 0x78, 0x4C},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := XORBlockAllocate(tt.args.src, tt.args.key, tt.args.iv)
			if (err != nil) != tt.wantErr {
				t.Errorf("XORBlockAllocate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("XORBlockAllocate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateIV(t *testing.T) {
	tests := []struct {
		name    string
		want    []byte
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateIV()
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateIV() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateIV() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestXORBlockAllocate1(t *testing.T) {
	type args struct {
		src []byte
		key []byte
		iv  []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := XORBlockAllocate(tt.args.src, tt.args.key, tt.args.iv)
			if (err != nil) != tt.wantErr {
				t.Errorf("XORBlockAllocate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("XORBlockAllocate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestXORBlockStream(t *testing.T) {
	type args struct {
		src []byte
		key []byte
		iv  []byte
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			err := XORBlockStream(w, tt.args.src, tt.args.key, tt.args.iv)
			if (err != nil) != tt.wantErr {
				t.Errorf("XORBlockStream() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("XORBlockStream() gotW = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
