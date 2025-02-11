package y

import "hash/crc32"

var (
	CastTagNoLiCrcTable = crc32.MakeTable(crc32.Castagnoli)
)
