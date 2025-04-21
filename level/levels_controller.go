package level

import (
	"github.com/dgraph-io/ristretto/v2/z"
	"sync/atomic"
	"wiscdb/internal"
	"wiscdb/pb"
	"wiscdb/table"
	"wiscdb/y"
)

type LevelsController struct {
	nextFileID atomic.Uint64
	l0stallMs  atomic.Int64
	levels     []*levelHandler
	kv         *internal.DB
	cstatus    compactStatus
}

func newLevelController(db *internal.DB, mf *internal.Manifest) (*LevelsController, error) {
	return &LevelsController{}, nil
}

func (s *LevelsController) close() error {
	return nil
}

func (s *LevelsController) cleanupLevels() error {
	return nil
}

func (s *LevelsController) addLevel0Table(t *table.Table) {

}

func (s *LevelsController) startCompact(lc *z.Closer) {

}

func (s *LevelsController) runCompactor(id int, lc *z.Closer) {

}

func (s *LevelsController) pickCompactLevels(priosBuffer []compactionPriority) (prios []compactionPriority) {
	return nil
}

func (s *LevelsController) levelTargets() targets {
	return targets{}
}

func (s *LevelsController) lastLevel() *levelHandler {
	return nil
}

func (s *LevelsController) doCompact(id int, p compactionPriority) error {
	return nil
}

func (s *LevelsController) runCompactDef(id, l int, cd compactDef) (err error) {
	return nil
}

func tableToString(tables []*table.Table) []string {
	return nil
}

func buildChangeSet(cd *compactDef, newTables []*table.Table) pb.ManifestChangeSet {
	return pb.ManifestChangeSet{}
}

func (s *LevelsController) compactBuildTables(lev int, cd compactDef) ([]*table.Table, func() error, error) {
	return nil, func() error {
		return nil
	}, nil
}

func (s *LevelsController) subCompact(it y.Iterator, kr keyRange, cd compactDef, inflightBuilders *y.Throttle, res chan<- *table.Table) {

}

func hasAnyPrefixes(s []byte, listOfPrefixes [][]byte) bool {
	return false
}

func (s *LevelsController) checkOverlap(tables []*table.Table, lev int) bool {
	return false
}
func (s *LevelsController) addSplits(cd *compactDef) {

}
func (s *LevelsController) fillTables(cd *compactDef) bool {
	return false
}

func (s *LevelsController) sortByHeuristic(tables []*table.Table, cd *compactDef) {

}

func (s *LevelsController) fillMaxLevelTables(tables []*table.Table, cd *compactDef) bool {

	return false
}

func (s *LevelsController) sortByStaleDataSize(tables []*table.Table, cd *compactDef) {

}

func (s *LevelsController) fillTablesL0(cd *compactDef) bool {
	return false
}
func (s *LevelsController) fillTablesL0ToL0(cd *compactDef) bool {
	return false
}

func (s *LevelsController) fillTablesL0ToLBase(cd *compactDef) bool {
	return false
}

func (s *LevelsController) appendIterator(iters []y.Iterator, opt *internal.IteratorOptions) []y.Iterator {
	return nil
}

func (s *LevelsController) keySplits(numPerTable int, prefix []byte) []string {
	return nil
}

func (s *LevelsController) getLevelInfo() []LevelInfo {
	return nil
}

func (s *LevelsController) verifyChecksum() error {
	return nil
}

func (s *LevelsController) AppendIterators(iters []y.Iterator, opt *internal.IteratorOptions) []y.Iterator {
	return nil
}
