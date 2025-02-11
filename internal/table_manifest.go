package internal

import "wiscdb/options"

type TableManifest struct {
	Level       uint8
	KeyID       uint64
	Compression options.CompressionType
}
