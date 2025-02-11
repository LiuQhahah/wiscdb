package options

type ChecksumVerificationMode int

const (
	NoVerification ChecksumVerificationMode = iota
	OnTableRead
	OnBlockRead
	OnTableAndBlockRead
)

type CompressionType uint32

const (
	None CompressionType = iota
	Snappy
	ZSTD
)
