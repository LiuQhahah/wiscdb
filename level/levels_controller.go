package level

import (
	"github.com/dgraph-io/ristretto/v2/z"
	"sync/atomic"
	"wiscdb/internal"
	"wiscdb/pb"
	"wiscdb/table"
	"wiscdb/y"
)

type levelsController struct {
	nextFileID atomic.Uint64
	l0stallMs  atomic.Int64
	levels     []*levelHandler
	kv         *internal.DB
	cstatus    compactStatus
}

func newLevelController(db *internal.DB, mf *internal.Manifest) (*levelsController, error) {
	return &levelsController{}, nil
}

func (s *levelsController) close() error {
	return nil
}

func (s *levelsController) cleanupLevels() error {
	return nil
}

func (s *levelsController) addLevel0Table(t *table.Table) {

}

func (s *levelsController) startCompact(lc *z.Closer) {

}

func (s *levelsController) runCompactor(id int, lc *z.Closer) {

}

func (s *levelsController) pickCompactLevels(priosBuffer []compactionPriority) (prios []compactionPriority) {
	return nil
}

func (s *levelsController) levelTargets() targets {
	return targets{}
}

func (s *levelsController) lastLevel() *levelHandler {
	return nil
}

func (s *levelsController) doCompact(id int, p compactionPriority) error {
	return nil
}

func (s *levelsController) runCompactDef(id, l int, cd compactDef) (err error) {
	return nil
}

func tableToString(tables []*table.Table) []string {
	return nil
}

func buildChangeSet(cd *compactDef, newTables []*table.Table) pb.ManifestChangeSet {
	return pb.ManifestChangeSet{}
}

func (s *levelsController) compactBuildTables(lev int, cd compactDef) ([]*table.Table, func() error, error) {
	return nil, func() error {
		return nil
	}, nil
}

func (s *levelsController) subCompact(it y.Iterator, kr keyRange, cd compactDef, inflightBuilders *y.Throttle, res chan<- *table.Table) {

}

func hasAnyPrefixes(s []byte, listOfPrefixes [][]byte) bool {
	return false
}

func (s *levelsController) checkOverlap(tables []*table.Table, lev int) bool {
	return false
}
func (s *levelsController) addSplits(cd *compactDef) {

}
func (s *levelsController) fillTables(cd *compactDef) bool {
	return false
}

func (s *levelsController) sortByHeuristic(tables []*table.Table, cd *compactDef) {

}

func (s *levelsController) fillMaxLevelTables(tables []*table.Table, cd *compactDef) bool {

	return false
}

func (s *levelsController) sortByStaleDataSize(tables []*table.Table, cd *compactDef) {

}

func (s *levelsController) fillTablesL0(cd *compactDef) bool {
	return false
}
func (s *levelsController) fillTablesL0ToL0(cd *compactDef) bool {
	return false
}

func (s *levelsController) fillTablesL0ToLBase(cd *compactDef) bool {
	return false
}

func (s *levelsController) appendIterator(iters []y.Iterator, opt *internal.IteratorOptions) []y.Iterator {
	return nil
}

func (s *levelsController) keySplits(numPerTable int, prefix []byte) []string {
	return nil
}

func (s *levelsController) getLevelInfo() []LevelInfo {
	return nil
}

func (s *levelsController) verifyChecksum() error {
	return nil
}
