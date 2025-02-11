package internal

import (
	"log"
	"testing"
)

func Test_defaultLog_Debugf(t *testing.T) {
	type fields struct {
		Logger *log.Logger
		level  loggingLevel
	}
	type args struct {
		f string
		v []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "TestDebug",
			fields: fields{
				level: 0,
			},
			args: args{
				f: "testForDEBUG",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &defaultLog{
				Logger: tt.fields.Logger,
				level:  tt.fields.level,
			}
			l.Debugf(tt.args.f, tt.args.v...)
		})
	}
}
