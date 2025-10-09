package main

import "testing"

func Test_keyRange_overlapsWith(t *testing.T) {
	type fields struct {
		left  []byte
		right []byte
		inf   bool
		size  int64
	}
	type args struct {
		dst keyRange
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "test",
			fields: fields{

				//后8位作为timestamp
				left:  []byte{0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
				right: []byte{0x10, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
			},
			args: args{
				dst: keyRange{
					left:  []byte{0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
					right: []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
				},
			},
			// 00-10 肯定在00-FF之间
			want: true,
		},
		{
			name: "test2",
			fields: fields{

				//后8位作为timestamp
				left:  []byte{0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
				right: []byte{0x10, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
			},
			args: args{
				dst: keyRange{
					left:  []byte{0x20, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
					right: []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
				},
			},
			// 00-10 不在20-FF之间
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := keyRange{
				left:  tt.fields.left,
				right: tt.fields.right,
				inf:   tt.fields.inf,
				size:  tt.fields.size,
			}
			if got := r.overlapsWith(tt.args.dst); got != tt.want {
				t.Errorf("overlapsWith() = %v, want %v", got, tt.want)
			}
		})
	}
}
