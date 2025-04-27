package y

import "expvar"

var (
	numIteratorsCreated *expvar.Int
	numBytesWrittenUser *expvar.Int
	numPuts             *expvar.Int
	pendingWrites       *expvar.Map
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

func PendingWritesSet(enabled bool, key string, val expvar.Var) {
	storeToMap(enabled, pendingWrites, key, val)
}

func storeToMap(enabled bool, metric *expvar.Map, key string, value expvar.Var) {
	if !enabled {
		return
	}
	metric.Set(key, value)
}
