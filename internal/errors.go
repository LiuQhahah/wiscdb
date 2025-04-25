package internal

import "errors"

var (
	ErrDBClosed = errors.New("DB Closed")

	ErrDiscardedTxn = errors.New("This transaction has been discarded. Create a new one")
	ErrReadOnlyTxn  = errors.New("No sets or deletes are allowed in a read-only transaction")
	ErrEmptyKey     = errors.New("Key cannot be empty")
	ErrInvalidKey   = errors.New("Key is using a reserved !wisc! prefix")
	ErrTxnTooBig    = errors.New("Txn is too big to fit into one request")
)
