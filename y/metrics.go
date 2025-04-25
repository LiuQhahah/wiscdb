package y

import "expvar"

var (
	numIteratorsCreated *expvar.Int
	numBytesWrittenUser *expvar.Int
	numPuts             *expvar.Int
)

const (
	BADGER_METRIC_PREFIX = "badger_"
)

func init() {
	numIteratorsCreated = expvar.NewInt(BADGER_METRIC_PREFIX + "iterator_num_user")
}
func addInt(enabled bool, metric *expvar.Int, val int64) {
	if !enabled {
		return
	}
	metric.Add(val)
}
func NumIteratorsCreatedAdd(enabled bool, val int64) {
	addInt(enabled, numIteratorsCreated, val)
}

func NumBytesWrittenUserAdd(enabled bool, val int64) {
	addInt(enabled, numBytesWrittenUser, val)
}
func NumPutsAdd(enabled bool, val int64) {
	addInt(enabled, numPuts, val)
}
