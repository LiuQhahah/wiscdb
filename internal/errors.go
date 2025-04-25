package internal

import "errors"

var (
	ErrDBClosed = errors.New("DB Closed")

	ErrDiscardTxn = errors.New("This transaction has been discarded. Create a new one")
)
