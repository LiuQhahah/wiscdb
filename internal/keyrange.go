package internal

import "wiscdb/table"

type keyRange struct {
	left  []byte
	right []byte
	inf   bool
	size  int64
}

func (r keyRange) isEmpty() bool {
	return false
}

var infRange = keyRange{inf: true}

func (r keyRange) String() string {
	return ""
}
func (r keyRange) equals(dst keyRange) bool {
	return false
}

func (r *keyRange) extend(kr keyRange) {

}

func (r keyRange) overlapsWith(dst keyRange) bool {
	return false
}

func getKeyRange(tables ...*table.Table) keyRange {
	return keyRange{}
}
