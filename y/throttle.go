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

// 抢占资源
func (t *Throttle) Do() error {
	for {
		select {
		case t.ch <- struct{}{}:
			t.wg.Add(1)
			return nil
		case err := <-t.errCh:
			if err != nil {
				return err
			}
		}
	}
}

func (t *Throttle) Done(err error) {
	if err != nil {
		t.errCh <- err
	}
	select {
	case <-t.ch: // channel中输出一个位置
	default:
		panic("Throttle Do Done mismatch")
	}
	// 释放一个位置
	t.wg.Done()

}

func (t *Throttle) Finish() error {
	t.once.Do(func() {
		//等待wg 都处于done的状态
		t.wg.Wait()
		close(t.ch)
		close(t.errCh)
		for err := range t.errCh {
			if err != nil {
				t.finishErr = err
				return
			}
		}
	})
	return t.finishErr
}
