package internal

type testOnlyExtensions struct {
	syncChan               chan string
	onClosedDiscardCapture map[uint64]uint64
}
