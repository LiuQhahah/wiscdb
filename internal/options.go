package internal

import (
	"time"
	"wiscdb/options"
)

const (
	maxValueThreshold = (1 << 20)
)

type Options struct {
	testOnlyExtensions
	Dir               string
	ValueDir          string
	SyncWrites        bool
	NumVersionsToKeep int
	ReadOnly          bool
	Logger            Logger
	Compression       options.CompressionType
	InMemory          bool
	MetricsEnabled    bool
	NumGoroutine      int

	MemTableSize        int64
	BaseTableSize       int64
	BaseLevelSize       int64
	LevelSizeMultiplier int
	TableSizeMultiplier int
	MaxLevels           int

	VLogPercentile float64
	ValueThreshold int64
	NumMemTables   int

	BlockSize          int
	BloomFalsePositive float64
	BlockCacheSize     int64
	IndexCacheSize     int64

	NumLevelZeroTables            int
	NumLevelZeroTablesStall       int
	ValueLogFileSize              int64
	ValueLogMaxEntries            uint32
	NumCompactors                 int
	CompactL0OnClose              bool
	LMaxCompaction                bool
	ZSTDCompressionLevel          int
	VerifyValueChecksum           bool
	EncryptionKey                 []byte
	EncryptionKeyRotationDuration time.Duration
	BypassLockGuard               bool
	ChecksumVerificationMode      options.ChecksumVerificationMode
	DetectConflicts               bool
	NamespaceOffset               int
	ExternalMagicVersion          uint16
	managedTxns                   bool
	maxBatchCount                 int64
	maxBatchSize                  int64
	maxValueThreshold             float64
}

func DefaultOptions(path string) Options {
	return Options{
		Dir:                           path,
		ValueDir:                      path,
		MemTableSize:                  64 << 20, //default 64MB
		BaseTableSize:                 2 << 20,  // default 2 MB
		BaseLevelSize:                 10 << 20, //default 10MB
		TableSizeMultiplier:           2,
		LevelSizeMultiplier:           10,
		MaxLevels:                     7,
		NumGoroutine:                  8,
		MetricsEnabled:                true,
		NumCompactors:                 4,
		NumLevelZeroTables:            5,
		NumLevelZeroTablesStall:       15,
		NumMemTables:                  5,
		BloomFalsePositive:            0.01,
		BlockSize:                     4 * 1024,
		SyncWrites:                    false,
		NumVersionsToKeep:             1,
		CompactL0OnClose:              false,
		VerifyValueChecksum:           false,
		Compression:                   options.Snappy,
		BlockCacheSize:                256 << 20, // default 256MB
		IndexCacheSize:                0,
		ZSTDCompressionLevel:          1,
		ValueLogFileSize:              1<<30 - 1, // default 1GB - 1
		ValueLogMaxEntries:            1000000,
		VLogPercentile:                0.0,
		ValueThreshold:                maxValueThreshold,
		Logger:                        defaultLogger(INFO),
		EncryptionKey:                 []byte{},
		EncryptionKeyRotationDuration: 10 * 24 * time.Hour,
		DetectConflicts:               true,
		NamespaceOffset:               -1,
	}
}
