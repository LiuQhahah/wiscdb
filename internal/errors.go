package internal

import "errors"

var (
	ErrDBClosed = errors.New("DB Closed")

	ErrDiscardedTxn     = errors.New("This transaction has been discarded. Create a new one")
	ErrReadOnlyTxn      = errors.New("No sets or deletes are allowed in a read-only transaction")
	ErrEmptyKey         = errors.New("Key cannot be empty")
	ErrInvalidKey       = errors.New("Key is using a reserved !wisc! prefix")
	ErrInvalidDataKeyID = errors.New("Invalid datakey id")
	ErrTxnTooBig        = errors.New("Txn is too big to fit into one request")
	ErrConflict         = errors.New("Transaction Conflict. Please retry")
	ErrBlockedWrites    = errors.New("Writes are blocked,possibly due to DropAll or Close")

	ErrBannedKey   = errors.New("Key is using the banned prefix")
	ErrKeyNotFound = errors.New("Key not found")
)
