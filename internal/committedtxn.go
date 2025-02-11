package internal

type committedTxn struct {
	ts           uint64
	conflictKeys map[uint64]struct{}
}
