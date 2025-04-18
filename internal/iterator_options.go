package internal

import "wiscdb/table"

type IteratorOptions struct {
	PrefetchSize   int
	PrefetchValues bool
	Reverse        bool
	AllVersions    bool
	InternalAccess bool
	prefixIsKey    bool
	Prefix         []byte
	SinceTs        uint64
}

// 设置默认的迭代选项、预捞取value、捞取100个、不反转、不捞取所有版本
var DefaultIteratorOptions = IteratorOptions{
	PrefetchValues: true,
	PrefetchSize:   100,
	Reverse:        false,
	AllVersions:    false,
}

func (opt *IteratorOptions) compareToPrefix(key []byte) int {
	return 0
}

func (opt *IteratorOptions) pickTable(t table.TableInterface) bool {
	return false
}

func (opt *IteratorOptions) pickTables(all []*table.Table) []*table.Table {
	return nil
}
