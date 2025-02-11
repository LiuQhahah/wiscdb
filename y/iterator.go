package y

type Iterator interface {
	Next()
	ReWind()
	Seek(key []byte)
	Value() ValueStruct
	Valid() bool
	Close() error
}
