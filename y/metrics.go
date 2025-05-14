package y

import "expvar"

var (
	numIteratorsCreated *expvar.Int
	numBytesWrittenUser *expvar.Int
	numPuts             *expvar.Int
	numGets             *expvar.Int
	numMemtableGets     *expvar.Int
	numGetsWithResults  *expvar.Int
	pendingWrites       *expvar.Map
	numBytesWrittenToL0 *expvar.Int
	numWritesVlog       *expvar.Int
	numBytesVlogWritten *expvar.Int
	numBytesReadLSM     *expvar.Int
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

func NumGetsAdd(enabled bool, val int64) {
	addInt(enabled, numGets, val)
}

func NumMemtableGetsAdd(enabled bool, val int64) {
	addInt(enabled, numMemtableGets, val)
}

func NumBytesReadsLSMAdd(enabled bool, val int64) {
	addInt(enabled, numBytesReadLSM, val)
}
func NumGetsWithResultAdd(enabled bool, val int64) {
	addInt(enabled, numGetsWithResults, val)
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
