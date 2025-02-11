package y

type mark struct {
	index   uint64
	waiter  chan struct{}
	indices []uint64
	done    bool
}
