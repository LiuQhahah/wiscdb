package y

import "sync"

type Throttle struct {
	once      sync.Once
	wg        sync.WaitGroup
	ch        chan struct{}
	errCh     chan error
	finishErr error
}

// 控制channel的大小
func NewThrottle(max int) *Throttle {
	return &Throttle{
		ch:    make(chan struct{}, max),
		errCh: make(chan error, max),
	}
}

func (t *Throttle) Do() error {
	return nil
}

func (t *Throttle) Done(err error) {

}

func (t *Throttle) Finish() error {
	return nil
}
