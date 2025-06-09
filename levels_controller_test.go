package main

import (
	"github.com/dgraph-io/ristretto/v2/z"
	"reflect"
	"sync/atomic"
	"testing"
	"wiscdb/pb"
	"wiscdb/table"
	"wiscdb/y"
)

func TestHasAnyPrefixes(t *testing.T) {
	type args struct {
		s              []byte
		listOfPrefixes [][]byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasAnyPrefixes(tt.args.s, tt.args.listOfPrefixes); got != tt.want {
				t.Errorf("HasAnyPrefixes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevelsController_AddLevel0Table(t *testing.T) {
	type fields struct {
		nextFileID    atomic.Uint64
		l0stallMs     atomic.Int64
		levels        []*levelHandler
		kv            *DB
		compactStatus compactStatus
	}
	type args struct {
		t *table.Table
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &LevelsController{
				nextFileID:    tt.fields.nextFileID,
				l0stallMs:     tt.fields.l0stallMs,
				levels:        tt.fields.levels,
				kv:            tt.fields.kv,
				compactStatus: tt.fields.compactStatus,
			}
			if err := s.AddLevel0Table(tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("AddLevel0Table() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLevelsController_AppendIterators(t *testing.T) {
	type fields struct {
		nextFileID    atomic.Uint64
		l0stallMs     atomic.Int64
		levels        []*levelHandler
		kv            *DB
		compactStatus compactStatus
	}
	type args struct {
		iters []y.Iterator
		opt   *IteratorOptions
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []y.Iterator
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &LevelsController{
				nextFileID:    tt.fields.nextFileID,
				l0stallMs:     tt.fields.l0stallMs,
				levels:        tt.fields.levels,
				kv:            tt.fields.kv,
				compactStatus: tt.fields.compactStatus,
			}
			if got := s.AppendIterators(tt.args.iters, tt.args.opt); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppendIterators() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevelsController_Get(t *testing.T) {
	type fields struct {
		nextFileID    atomic.Uint64
		l0stallMs     atomic.Int64
		levels        []*levelHandler
		kv            *DB
		compactStatus compactStatus
	}
	type args struct {
		key        []byte
		maxVs      y.ValueStruct
		startLevel int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    y.ValueStruct
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &LevelsController{
				nextFileID:    tt.fields.nextFileID,
				l0stallMs:     tt.fields.l0stallMs,
				levels:        tt.fields.levels,
				kv:            tt.fields.kv,
				compactStatus: tt.fields.compactStatus,
			}
			got, err := s.Get(tt.args.key, tt.args.maxVs, tt.args.startLevel)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevelsController_ReserveFileID(t *testing.T) {
	type fields struct {
		nextFileID    atomic.Uint64
		l0stallMs     atomic.Int64
		levels        []*levelHandler
		kv            *DB
		compactStatus compactStatus
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
			s := &LevelsController{
				nextFileID:    tt.fields.nextFileID,
				l0stallMs:     tt.fields.l0stallMs,
				levels:        tt.fields.levels,
				kv:            tt.fields.kv,
				compactStatus: tt.fields.compactStatus,
			}
			if got := s.ReserveFileID(); got != tt.want {
				t.Errorf("ReserveFileID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevelsController_addSplits(t *testing.T) {
	type fields struct {
		nextFileID    atomic.Uint64
		l0stallMs     atomic.Int64
		levels        []*levelHandler
		kv            *DB
		compactStatus compactStatus
	}
	type args struct {
		cd *compactDef
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
			s := &LevelsController{
				nextFileID:    tt.fields.nextFileID,
				l0stallMs:     tt.fields.l0stallMs,
				levels:        tt.fields.levels,
				kv:            tt.fields.kv,
				compactStatus: tt.fields.compactStatus,
			}
			s.addSplits(tt.args.cd)
		})
	}
}

func TestLevelsController_appendIterator(t *testing.T) {
	type fields struct {
		nextFileID    atomic.Uint64
		l0stallMs     atomic.Int64
		levels        []*levelHandler
		kv            *DB
		compactStatus compactStatus
	}
	type args struct {
		iters []y.Iterator
		opt   *IteratorOptions
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []y.Iterator
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &LevelsController{
				nextFileID:    tt.fields.nextFileID,
				l0stallMs:     tt.fields.l0stallMs,
				levels:        tt.fields.levels,
				kv:            tt.fields.kv,
				compactStatus: tt.fields.compactStatus,
			}
			if got := s.appendIterator(tt.args.iters, tt.args.opt); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appendIterator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevelsController_checkOverlap(t *testing.T) {
	type fields struct {
		nextFileID    atomic.Uint64
		l0stallMs     atomic.Int64
		levels        []*levelHandler
		kv            *DB
		compactStatus compactStatus
	}
	type args struct {
		tables []*table.Table
		lev    int
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
			s := &LevelsController{
				nextFileID:    tt.fields.nextFileID,
				l0stallMs:     tt.fields.l0stallMs,
				levels:        tt.fields.levels,
				kv:            tt.fields.kv,
				compactStatus: tt.fields.compactStatus,
			}
			if got := s.checkOverlap(tt.args.tables, tt.args.lev); got != tt.want {
				t.Errorf("checkOverlap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevelsController_cleanupLevels(t *testing.T) {
	type fields struct {
		nextFileID    atomic.Uint64
		l0stallMs     atomic.Int64
		levels        []*levelHandler
		kv            *DB
		compactStatus compactStatus
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
			s := &LevelsController{
				nextFileID:    tt.fields.nextFileID,
				l0stallMs:     tt.fields.l0stallMs,
				levels:        tt.fields.levels,
				kv:            tt.fields.kv,
				compactStatus: tt.fields.compactStatus,
			}
			if err := s.cleanupLevels(); (err != nil) != tt.wantErr {
				t.Errorf("cleanupLevels() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLevelsController_close(t *testing.T) {
	type fields struct {
		nextFileID    atomic.Uint64
		l0stallMs     atomic.Int64
		levels        []*levelHandler
		kv            *DB
		compactStatus compactStatus
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
			s := &LevelsController{
				nextFileID:    tt.fields.nextFileID,
				l0stallMs:     tt.fields.l0stallMs,
				levels:        tt.fields.levels,
				kv:            tt.fields.kv,
				compactStatus: tt.fields.compactStatus,
			}
			if err := s.close(); (err != nil) != tt.wantErr {
				t.Errorf("close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLevelsController_compactBuildTables(t *testing.T) {
	type fields struct {
		nextFileID    atomic.Uint64
		l0stallMs     atomic.Int64
		levels        []*levelHandler
		kv            *DB
		compactStatus compactStatus
	}
	type args struct {
		lev int
		cd  compactDef
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*table.Table
		want1   func() error
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &LevelsController{
				nextFileID:    tt.fields.nextFileID,
				l0stallMs:     tt.fields.l0stallMs,
				levels:        tt.fields.levels,
				kv:            tt.fields.kv,
				compactStatus: tt.fields.compactStatus,
			}
			got, got1, err := s.compactBuildTables(tt.args.lev, tt.args.cd)
			if (err != nil) != tt.wantErr {
				t.Errorf("compactBuildTables() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("compactBuildTables() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("compactBuildTables() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestLevelsController_doCompact(t *testing.T) {
	type fields struct {
		nextFileID    atomic.Uint64
		l0stallMs     atomic.Int64
		levels        []*levelHandler
		kv            *DB
		compactStatus compactStatus
	}
	type args struct {
		id int
		p  compactionPriority
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &LevelsController{
				nextFileID:    tt.fields.nextFileID,
				l0stallMs:     tt.fields.l0stallMs,
				levels:        tt.fields.levels,
				kv:            tt.fields.kv,
				compactStatus: tt.fields.compactStatus,
			}
			if err := s.doCompact(tt.args.id, tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("doCompact() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLevelsController_fillMaxLevelTables(t *testing.T) {
	type fields struct {
		nextFileID    atomic.Uint64
		l0stallMs     atomic.Int64
		levels        []*levelHandler
		kv            *DB
		compactStatus compactStatus
	}
	type args struct {
		tables []*table.Table
		cd     *compactDef
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
			s := &LevelsController{
				nextFileID:    tt.fields.nextFileID,
				l0stallMs:     tt.fields.l0stallMs,
				levels:        tt.fields.levels,
				kv:            tt.fields.kv,
				compactStatus: tt.fields.compactStatus,
			}
			if got := s.fillMaxLevelTables(tt.args.tables, tt.args.cd); got != tt.want {
				t.Errorf("fillMaxLevelTables() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevelsController_fillTables(t *testing.T) {
	type fields struct {
		nextFileID    atomic.Uint64
		l0stallMs     atomic.Int64
		levels        []*levelHandler
		kv            *DB
		compactStatus compactStatus
	}
	type args struct {
		cd *compactDef
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
			s := &LevelsController{
				nextFileID:    tt.fields.nextFileID,
				l0stallMs:     tt.fields.l0stallMs,
				levels:        tt.fields.levels,
				kv:            tt.fields.kv,
				compactStatus: tt.fields.compactStatus,
			}
			if got := s.fillTables(tt.args.cd); got != tt.want {
				t.Errorf("fillTables() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevelsController_fillTablesL0(t *testing.T) {
	type fields struct {
		nextFileID    atomic.Uint64
		l0stallMs     atomic.Int64
		levels        []*levelHandler
		kv            *DB
		compactStatus compactStatus
	}
	type args struct {
		cd *compactDef
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
			s := &LevelsController{
				nextFileID:    tt.fields.nextFileID,
				l0stallMs:     tt.fields.l0stallMs,
				levels:        tt.fields.levels,
				kv:            tt.fields.kv,
				compactStatus: tt.fields.compactStatus,
			}
			if got := s.fillTablesL0(tt.args.cd); got != tt.want {
				t.Errorf("fillTablesL0() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevelsController_fillTablesL0ToL0(t *testing.T) {
	type fields struct {
		nextFileID    atomic.Uint64
		l0stallMs     atomic.Int64
		levels        []*levelHandler
		kv            *DB
		compactStatus compactStatus
	}
	type args struct {
		cd *compactDef
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
			s := &LevelsController{
				nextFileID:    tt.fields.nextFileID,
				l0stallMs:     tt.fields.l0stallMs,
				levels:        tt.fields.levels,
				kv:            tt.fields.kv,
				compactStatus: tt.fields.compactStatus,
			}
			if got := s.fillTablesL0ToL0(tt.args.cd); got != tt.want {
				t.Errorf("fillTablesL0ToL0() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevelsController_fillTablesL0ToLBase(t *testing.T) {
	type fields struct {
		nextFileID    atomic.Uint64
		l0stallMs     atomic.Int64
		levels        []*levelHandler
		kv            *DB
		compactStatus compactStatus
	}
	type args struct {
		cd *compactDef
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
			s := &LevelsController{
				nextFileID:    tt.fields.nextFileID,
				l0stallMs:     tt.fields.l0stallMs,
				levels:        tt.fields.levels,
				kv:            tt.fields.kv,
				compactStatus: tt.fields.compactStatus,
			}
			if got := s.fillTablesL0ToLBase(tt.args.cd); got != tt.want {
				t.Errorf("fillTablesL0ToLBase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevelsController_getLevelInfo(t *testing.T) {
	type fields struct {
		nextFileID    atomic.Uint64
		l0stallMs     atomic.Int64
		levels        []*levelHandler
		kv            *DB
		compactStatus compactStatus
	}
	tests := []struct {
		name   string
		fields fields
		want   []LevelInfo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &LevelsController{
				nextFileID:    tt.fields.nextFileID,
				l0stallMs:     tt.fields.l0stallMs,
				levels:        tt.fields.levels,
				kv:            tt.fields.kv,
				compactStatus: tt.fields.compactStatus,
			}
			if got := s.getLevelInfo(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getLevelInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevelsController_keySplits(t *testing.T) {
	type fields struct {
		nextFileID    atomic.Uint64
		l0stallMs     atomic.Int64
		levels        []*levelHandler
		kv            *DB
		compactStatus compactStatus
	}
	type args struct {
		numPerTable int
		prefix      []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &LevelsController{
				nextFileID:    tt.fields.nextFileID,
				l0stallMs:     tt.fields.l0stallMs,
				levels:        tt.fields.levels,
				kv:            tt.fields.kv,
				compactStatus: tt.fields.compactStatus,
			}
			if got := s.keySplits(tt.args.numPerTable, tt.args.prefix); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("keySplits() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevelsController_lastLevel(t *testing.T) {
	type fields struct {
		nextFileID    atomic.Uint64
		l0stallMs     atomic.Int64
		levels        []*levelHandler
		kv            *DB
		compactStatus compactStatus
	}
	tests := []struct {
		name   string
		fields fields
		want   *levelHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &LevelsController{
				nextFileID:    tt.fields.nextFileID,
				l0stallMs:     tt.fields.l0stallMs,
				levels:        tt.fields.levels,
				kv:            tt.fields.kv,
				compactStatus: tt.fields.compactStatus,
			}
			if got := s.lastLevel(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("lastLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevelsController_levelTargets(t *testing.T) {
	type fields struct {
		nextFileID    atomic.Uint64
		l0stallMs     atomic.Int64
		levels        []*levelHandler
		kv            *DB
		compactStatus compactStatus
	}
	tests := []struct {
		name   string
		fields fields
		want   targets
	}{
		{
			name:   "test",
			fields: fields{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &LevelsController{
				nextFileID:    tt.fields.nextFileID,
				l0stallMs:     tt.fields.l0stallMs,
				levels:        tt.fields.levels,
				kv:            tt.fields.kv,
				compactStatus: tt.fields.compactStatus,
			}
			if got := s.levelTargets(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("levelTargets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevelsController_pickCompactLevels(t *testing.T) {
	type fields struct {
		nextFileID    atomic.Uint64
		l0stallMs     atomic.Int64
		levels        []*levelHandler
		kv            *DB
		compactStatus compactStatus
	}
	type args struct {
		priosBuffer []compactionPriority
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantPrios []compactionPriority
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &LevelsController{
				nextFileID:    tt.fields.nextFileID,
				l0stallMs:     tt.fields.l0stallMs,
				levels:        tt.fields.levels,
				kv:            tt.fields.kv,
				compactStatus: tt.fields.compactStatus,
			}
			if gotPrios := s.pickCompactLevels(tt.args.priosBuffer); !reflect.DeepEqual(gotPrios, tt.wantPrios) {
				t.Errorf("pickCompactLevels() = %v, want %v", gotPrios, tt.wantPrios)
			}
		})
	}
}

func TestLevelsController_runCompactDef(t *testing.T) {
	type fields struct {
		nextFileID    atomic.Uint64
		l0stallMs     atomic.Int64
		levels        []*levelHandler
		kv            *DB
		compactStatus compactStatus
	}
	type args struct {
		id int
		l  int
		cd compactDef
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &LevelsController{
				nextFileID:    tt.fields.nextFileID,
				l0stallMs:     tt.fields.l0stallMs,
				levels:        tt.fields.levels,
				kv:            tt.fields.kv,
				compactStatus: tt.fields.compactStatus,
			}
			if err := s.runCompactDef(tt.args.id, tt.args.l, tt.args.cd); (err != nil) != tt.wantErr {
				t.Errorf("runCompactDef() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLevelsController_runCompactor(t *testing.T) {
	type fields struct {
		nextFileID    atomic.Uint64
		l0stallMs     atomic.Int64
		levels        []*levelHandler
		kv            *DB
		compactStatus compactStatus
	}
	type args struct {
		id int
		lc *z.Closer
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
			s := &LevelsController{
				nextFileID:    tt.fields.nextFileID,
				l0stallMs:     tt.fields.l0stallMs,
				levels:        tt.fields.levels,
				kv:            tt.fields.kv,
				compactStatus: tt.fields.compactStatus,
			}
			s.runCompactor(tt.args.id, tt.args.lc)
		})
	}
}

func TestLevelsController_sortByHeuristic(t *testing.T) {
	type fields struct {
		nextFileID    atomic.Uint64
		l0stallMs     atomic.Int64
		levels        []*levelHandler
		kv            *DB
		compactStatus compactStatus
	}
	type args struct {
		tables []*table.Table
		cd     *compactDef
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
			s := &LevelsController{
				nextFileID:    tt.fields.nextFileID,
				l0stallMs:     tt.fields.l0stallMs,
				levels:        tt.fields.levels,
				kv:            tt.fields.kv,
				compactStatus: tt.fields.compactStatus,
			}
			s.sortByHeuristic(tt.args.tables, tt.args.cd)
		})
	}
}

func TestLevelsController_sortByStaleDataSize(t *testing.T) {
	type fields struct {
		nextFileID    atomic.Uint64
		l0stallMs     atomic.Int64
		levels        []*levelHandler
		kv            *DB
		compactStatus compactStatus
	}
	type args struct {
		tables []*table.Table
		cd     *compactDef
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
			s := &LevelsController{
				nextFileID:    tt.fields.nextFileID,
				l0stallMs:     tt.fields.l0stallMs,
				levels:        tt.fields.levels,
				kv:            tt.fields.kv,
				compactStatus: tt.fields.compactStatus,
			}
			s.sortByStaleDataSize(tt.args.tables, tt.args.cd)
		})
	}
}

func TestLevelsController_startCompact(t *testing.T) {
	type fields struct {
		nextFileID    atomic.Uint64
		l0stallMs     atomic.Int64
		levels        []*levelHandler
		kv            *DB
		compactStatus compactStatus
	}
	type args struct {
		lc *z.Closer
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
			s := &LevelsController{
				nextFileID:    tt.fields.nextFileID,
				l0stallMs:     tt.fields.l0stallMs,
				levels:        tt.fields.levels,
				kv:            tt.fields.kv,
				compactStatus: tt.fields.compactStatus,
			}
			s.startCompact(tt.args.lc)
		})
	}
}

func TestLevelsController_subCompact(t *testing.T) {
	type fields struct {
		nextFileID    atomic.Uint64
		l0stallMs     atomic.Int64
		levels        []*levelHandler
		kv            *DB
		compactStatus compactStatus
	}
	type args struct {
		it               y.Iterator
		kr               keyRange
		cd               compactDef
		inflightBuilders *y.Throttle
		res              chan<- *table.Table
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
			s := &LevelsController{
				nextFileID:    tt.fields.nextFileID,
				l0stallMs:     tt.fields.l0stallMs,
				levels:        tt.fields.levels,
				kv:            tt.fields.kv,
				compactStatus: tt.fields.compactStatus,
			}
			s.subCompact(tt.args.it, tt.args.kr, tt.args.cd, tt.args.inflightBuilders, tt.args.res)
		})
	}
}

func TestLevelsController_validate(t *testing.T) {
	type fields struct {
		nextFileID    atomic.Uint64
		l0stallMs     atomic.Int64
		levels        []*levelHandler
		kv            *DB
		compactStatus compactStatus
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
			s := &LevelsController{
				nextFileID:    tt.fields.nextFileID,
				l0stallMs:     tt.fields.l0stallMs,
				levels:        tt.fields.levels,
				kv:            tt.fields.kv,
				compactStatus: tt.fields.compactStatus,
			}
			if err := s.validate(); (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLevelsController_verifyChecksum(t *testing.T) {
	type fields struct {
		nextFileID    atomic.Uint64
		l0stallMs     atomic.Int64
		levels        []*levelHandler
		kv            *DB
		compactStatus compactStatus
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
			s := &LevelsController{
				nextFileID:    tt.fields.nextFileID,
				l0stallMs:     tt.fields.l0stallMs,
				levels:        tt.fields.levels,
				kv:            tt.fields.kv,
				compactStatus: tt.fields.compactStatus,
			}
			if err := s.verifyChecksum(); (err != nil) != tt.wantErr {
				t.Errorf("verifyChecksum() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_buildChangeSet(t *testing.T) {
	type args struct {
		cd        *compactDef
		newTables []*table.Table
	}
	tests := []struct {
		name string
		args args
		want pb.ManifestChangeSet
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildChangeSet(tt.args.cd, tt.args.newTables); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildChangeSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_closeAllTables(t *testing.T) {
	type args struct {
		tables [][]*table.Table
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			closeAllTables(tt.args.tables)
		})
	}
}

func Test_getIDMap(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name string
		args args
		want map[uint64]struct{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getIDMap(tt.args.dir); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getIDMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hasAnyPrefixes(t *testing.T) {
	type args struct {
		s              []byte
		listOfPrefixes [][]byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasAnyPrefixes(tt.args.s, tt.args.listOfPrefixes); got != tt.want {
				t.Errorf("hasAnyPrefixes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newLevelController(t *testing.T) {
	type args struct {
		db *DB
		mf *Manifest
	}
	tests := []struct {
		name    string
		args    args
		want    *LevelsController
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newLevelController(tt.args.db, tt.args.mf)
			if (err != nil) != tt.wantErr {
				t.Errorf("newLevelController() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newLevelController() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newLevelsController(t *testing.T) {
	type args struct {
		db *DB
		mf *Manifest
	}
	tests := []struct {
		name    string
		args    args
		want    *LevelsController
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newLevelsController(tt.args.db, tt.args.mf)
			if (err != nil) != tt.wantErr {
				t.Errorf("newLevelsController() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newLevelsController() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_revertToManifest(t *testing.T) {
	type args struct {
		kv    *DB
		mf    *Manifest
		idMap map[uint64]struct{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := revertToManifest(tt.args.kv, tt.args.mf, tt.args.idMap); (err != nil) != tt.wantErr {
				t.Errorf("revertToManifest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_tableToString(t *testing.T) {
	type args struct {
		tables []*table.Table
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tableToString(tt.args.tables); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tableToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
