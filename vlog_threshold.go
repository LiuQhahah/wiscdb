package main

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

// 将当前valuelog的文件长度写到value channel中
func (v *vlogThreshold) update(sizes []int64) {
	v.valueCh <- sizes
}
func (v *vlogThreshold) close() {

}

func (v *vlogThreshold) listenForValueThresholdUpdate() {

}
