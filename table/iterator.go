package table

var (
	REVERSED int = 2
	NOCACHE  int = 4
)

type Iterator struct {
	t    *Table
	bpos int
	bi   blockIterator
	err  error
	opt  int
}

type blockIterator struct {
	data         []byte
	idx          int
	err          error
	baseKey      []byte
	key          []byte
	val          []byte
	entryOffsets []uint32
	block        *Block
	tableID      uint64
	blockID      int
	prevOverlap  uint16
}

func (itr *Iterator) Rewind() {
	if itr.opt&REVERSED == 0 {
		itr.seekToFirst()
	} else {
		itr.seekToFirst()
	}
}

func (itr *Iterator) seekToFirst() {

}
func (itr *Iterator) seekToLast() {

}
func (itr *blockIterator) seekToFirst() {

}
func (itr *blockIterator) seekToLast() {

}
func (itr *Iterator) Close() error {
	return nil
}
func (itr *Iterator) Valid() bool {
	return false
}

func (itr *Iterator) Key() []byte {
	return nil
}
