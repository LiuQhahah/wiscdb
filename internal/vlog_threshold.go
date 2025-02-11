package internal

import (
	"github.com/dgraph-io/ristretto/v2/z"
	"sync/atomic"
)

type vlogThreshold struct {
	logger         Logger
	percentile     float64
	valueThreshold atomic.Int64
	valueCh        chan []int64
	clearCh        chan bool
	closer         *z.Closer
	vlMetrics      *z.HistogramData
}

func initVlogThreshold(opt *Options) *vlogThreshold {
	return &vlogThreshold{}
}

func (v *vlogThreshold) Clear(opt Options) {

}
func (v *vlogThreshold) update(sizes []int64) {

}
func (v *vlogThreshold) close() {

}

func (v *vlogThreshold) listenForValueThresholdUpdate() {

}
