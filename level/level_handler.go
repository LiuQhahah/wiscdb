package level

import (
	"sync"
	"wiscdb/internal"
	"wiscdb/table"
	"wiscdb/y"
)

type levelHandler struct {
	sync.RWMutex

	tables        []*table.Table
	totalSize     int64
	totalStalSize int64
	level         int
	strLevel      string
	db            *internal.DB
}

func newLevelHandler(db *internal.DB, level int) *levelHandler {
	return &levelHandler{}
}
func (s *levelHandler) isLastLevel() bool {
	return false
}

func (s *levelHandler) getTotalSize() int64 {
	return 0
}

func (s *levelHandler) initTables(tables []*table.Table) {

}

func (s *levelHandler) deleteTables(toDel []*table.Table) error {
	return nil
}

func (s *levelHandler) subtraceSize(t *table.Table) {

}
func (s *levelHandler) addSize(t *table.Table) {

}

func (s *levelHandler) replaceTables(toDel, toAdd []*table.Table) error {
	return nil
}
func (s *levelHandler) addTable(t *table.Table) {

}
func (s *levelHandler) sortTables() {

}
func (s *levelHandler) tryAddLevel0Table(t *table.Table) bool {
	return false
}
func (s *levelHandler) numTables() int {
	return 0
}
func (s *levelHandler) close() error {
	return nil
}

func (s *levelHandler) getTableForKey(key []byte) ([]*table.Table, func() error) {
	return nil, func() error {
		return nil
	}
}
func (s *levelHandler) get(key []byte) (y.ValueStruct, error) {
	return y.ValueStruct{}, nil
}

func (s *levelHandler) overlappingTables(kr keyRange) (int, int) {
	return 0, 0
}

func decrRefs(tables []*table.Table) error {
	return nil
}
