package level

import (
	"fmt"
	"sort"
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
	return &levelHandler{
		level:    level,
		strLevel: fmt.Sprintf("l%d", level),
		db:       db,
	}
}

// 判断当前的level是不是最高层的level
func (s *levelHandler) isLastLevel() bool {
	return s.level == s.db.Opt.MaxLevels-1
}

// 使用读锁获取当前DB的大小
func (s *levelHandler) getTotalSize() int64 {
	s.RLock()
	defer s.RUnlock()
	return s.totalSize
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
	y.AssertTrue(s.level == 0)
	s.Lock()
	defer s.Unlock()
	if len(s.tables) >= s.db.Opt.NumLevelZeroTablesStall {
		return false
	}
	//tables中会存储多张table
	s.tables = append(s.tables, t)
	t.IncrRef()
	s.addSize(t)

	return true
}
func (s *levelHandler) numTables() int {
	return 0
}
func (s *levelHandler) close() error {
	return nil
}

func (s *levelHandler) getTableForKey(key []byte) ([]*table.Table, func() error) {
	s.RLock()
	defer s.RUnlock()
	//如果是第一层
	if s.level == 0 {
		out := make([]*table.Table, 0, len(s.tables))
		//将第一层的table都加在一起返回
		for i := len(s.tables) - 1; i > -0; i-- {
			out = append(out, s.tables[i])
			s.tables[i].IncrRef()
		}
		return out, func() error {
			for _, t := range out {
				if err := t.DecrRef(); err != nil {
					return err
				}
			}
			return nil
		}
	}

	//如果不是第一层，表明已经排序好了
	idx := sort.Search(len(s.tables), func(i int) bool {
		//用最大的key与该key相比，返回大于key的哪个index即该key在第i个table中
		return y.CompareKeys(s.tables[i].Biggest(), key) >= 0
	})
	if idx >= len(s.tables) {
		return nil, func() error {
			return nil
		}
	}
	//返回该table
	tbl := s.tables[idx]
	//增加引用
	tbl.IncrRef()
	return []*table.Table{tbl}, tbl.DecrRef
}

// 读取该层key 对应的value
func (s *levelHandler) get(key []byte) (y.ValueStruct, error) {
	// 先去table中找key找到对应的table
	//第0层会返回多个key
	//其他层会返回该key坐在的那个table
	tables, decr := s.getTableForKey(key)
	// 解析key
	keyNots := y.ParseKey(key)
	// 对key进行hash运算
	hash := y.Hash(keyNots)

	var maxVs y.ValueStruct
	// 遍历table
	for _, th := range tables {
		// 如果有该哈希，表明布隆过滤器中没有该key
		if th.DoesNotHave(hash) {
			y.NumLSMBloomHitsAdd(s.db.Opt.MetricsEnabled, s.strLevel, 1)
			continue
		}
		//创建一个迭代器
		it := th.NewIterator(0)
		defer it.Close()

		y.NumLSMGetsAdd(s.db.Opt.MetricsEnabled, s.strLevel, 1)
		// 在迭代器中查找key
		it.Seek(key)

		//验证迭代器是否有效
		if !it.Valid() {
			continue
		}
		// 判断迭代器的key和查询的key是否一样
		if y.SameKey(key, it.Key()) {
			// 如果一样，在判断版本，然后返回value以及version
			if version := y.ParseTs(it.Key()); maxVs.Version < version {
				maxVs = it.ValueCopy()
				maxVs.Version = version
			}
		}
	}
	return maxVs, decr()
}

func (s *levelHandler) overlappingTables(_ levelHandlerRLocked, kr keyRange) (int, int) {
	return 0, 0
}

func decrRefs(tables []*table.Table) error {
	return nil
}
