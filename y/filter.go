package y

import "math"

type Filter []byte

func (f Filter) MayContainKey(key []byte) bool {
	return false
}

func (f Filter) MayContain(h uint32) bool {
	return false
}

// BloomBitsPerKey 返回 bloomfilter 根据误报率计算的每个密钥所需的比特数。
func BloomBitsPerKey(numEntries int, fp float64) int {
	size := -1 * float64(numEntries) * math.Log(fp) / math.Pow(float64(0.69314718056), 2)
	locs := math.Ceil(float64(0.69314718056) * size / float64(numEntries))
	return int(locs)
}

func NewFilter(keys []uint32, bitsPerKey int) Filter {
	return Filter(appendFilter(nil, keys, bitsPerKey))
}

func appendFilter(buf []byte, keys []uint32, bitsPerKey int) []byte {
	if bitsPerKey < 0 {
		bitsPerKey = 0
	}

	k := uint32(float64(bitsPerKey) * 0.69)
	if k < 1 {
		k = 1
	}
	if k > 30 {
		k = 30
	}
	nBits := len(keys) * bitsPerKey

	if nBits < 64 {
		nBits = 64
	}
	nBytes := (nBits + 7) / 8
	nBits = nBytes * 8
	buf, filter := extend(buf, nBytes+1)
	for _, h := range keys {
		delta := h>>17 | h<<15
		for j := uint32(0); j < k; j++ {
			bitPos := h % uint32(nBits)
			filter[bitPos/8] |= 1 << (bitPos % 8)
			h += delta
		}
	}
	filter[nBytes] = uint8(k)
	return buf
}

func extend(b []byte, n int) (overall, trailer []byte) {
	want := n + len(b)
	if want <= cap(b) {
		overall = b[:want]
		trailer = overall[len(b):]
		for i := range trailer {
			trailer[i] = 0
		}
	} else {
		c := 1024
		for c < want {
			c += c / 4
		}
		overall = make([]byte, want, c)
		trailer = overall[len(b):]
		copy(overall, b)
	}
	return overall, trailer
}

// 实现 Murmur 散列函数
func Hash(key []byte) uint32 {
	const (
		seed = 0xbc9f1d34
		m    = 0xc6a4a793
	)
	h := uint32(seed) ^ uint32(len(key))*m
	for ; len(key) >= 4; key = key[4:] {
		h += uint32(key[0]) | uint32(key[1])<<8 | uint32(key[2])<<16 | uint32(key[3])<<24
		h *= m
		h ^= h >> 16
	}
	switch len(key) {
	case 3:
		h += uint32(key[2]) << 16
		fallthrough
	case 2:
		h += uint32(key[1]) << 8
		fallthrough
	case 1:
		h += uint32(key[0])
		h *= m
		h ^= h >> 24
	}
	return h
}
