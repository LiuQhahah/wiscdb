package y

import (
	"reflect"
	"sync"
	"testing"
)

func TestNewThrottle(t *testing.T) {
	type args struct {
		max int
	}
	tests := []struct {
		name string
		args args
		want *Throttle
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewThrottle(tt.args.max); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewThrottle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestThrottle_Do(t1 *testing.T) {
	type fields struct {
		once      sync.Once
		wg        sync.WaitGroup
		ch        chan struct{}
		errCh     chan error
		finishErr error
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Throttle{
				once:      tt.fields.once,
				wg:        tt.fields.wg,
				ch:        tt.fields.ch,
				errCh:     tt.fields.errCh,
				finishErr: tt.fields.finishErr,
			}
			if err := t.Do(); (err != nil) != tt.wantErr {
				t1.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestThrottle_Done(t1 *testing.T) {
	type fields struct {
		once      sync.Once
		wg        sync.WaitGroup
		ch        chan struct{}
		errCh     chan error
		finishErr error
	}
	type args struct {
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Throttle{
				once:      tt.fields.once,
				wg:        tt.fields.wg,
				ch:        tt.fields.ch,
				errCh:     tt.fields.errCh,
				finishErr: tt.fields.finishErr,
			}
			t.Done(tt.args.err)
		})
	}
}

func TestThrottle_Finish(t1 *testing.T) {
	type fields struct {
		once      sync.Once
		wg        sync.WaitGroup
		ch        chan struct{}
		errCh     chan error
		finishErr error
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Throttle{
				once:      tt.fields.once,
				wg:        tt.fields.wg,
				ch:        tt.fields.ch,
				errCh:     tt.fields.errCh,
				finishErr: tt.fields.finishErr,
			}
			if err := t.Finish(); (err != nil) != tt.wantErr {
				t1.Errorf("Finish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
