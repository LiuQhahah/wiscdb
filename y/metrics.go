package y

import "expvar"

var (
	numIteratorsCreated *expvar.Int
	numBytesWrittenUser *expvar.Int
	numPuts             *expvar.Int
	pendingWrites       *expvar.Map
	numBytesWrittenToL0 *expvar.Int
	numWritesVlog       *expvar.Int
	numBytesVlogWritten *expvar.Int
	numLSMBloomHits     *expvar.Map
	numLSMGets          *expvar.Map
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

func NumBytesWrittenToL0Add(enabled bool, val int64) {
	addInt(enabled, numBytesWrittenToL0, val)
}

func NumWritesVlogAdd(enabled bool, val int64) {
	addInt(enabled, numWritesVlog, val)
}

func NumBytesWrittenVlogAdd(enabled bool, val int64) {
	addInt(enabled, numBytesVlogWritten, val)
}

func NumLSMBloomHitsAdd(enabled bool, key string, val int64) {
	addToMap(enabled, numLSMBloomHits, key, val)
}

func NumLSMGetsAdd(enabled bool, key string, val int64) {
	addToMap(enabled, numLSMGets, key, val)
}
func addToMap(enabled bool, metric *expvar.Map, key string, val int64) {
	if !enabled {
		return
	}
	metric.Add(key, val)
}
