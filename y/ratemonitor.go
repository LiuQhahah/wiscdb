package y

import "time"

type RateMonitor struct {
	start       time.Time
	lastSent    uint64
	lastCapture time.Time
	rates       []float64
	idx         int
}

func NewRateMonitor(numSamples int) *RateMonitor {
	return &RateMonitor{}
}

const minRate = 0.0001

func (rm *RateMonitor) Capture(sent uint64) {

}

func (rm *RateMonitor) Rate() uint64 {
	return 0
}
