package y

type Filter []byte

func (f Filter) MayContainKey(key []byte) bool {
	return false
}

func (f Filter) MayContain(h uint32) bool {
	return false
}

func BloomBitsPerKey(numEntries int, fp float64) int {
	return 0
}

func NewFilter(keys []uint32, bitsPerKey int) Filter {
	return Filter{}
}

func appendFilter(buf []byte, keys []uint32, bitsPerKey int) []byte {

	return nil
}

func extend(b []byte, n int) (overall, trailer []byte) {
	return nil, nil
}

func Hash(key []byte) uint32 {
	return 0
}
