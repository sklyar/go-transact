package txsql

import "context"

// IsolationLevel is the transaction isolation level used in TxOptions.
type IsolationLevel int

// Various isolation levels that drivers may support in BeginTx.
// If a driver does not support a given isolation level an error may be returned.
//
// See https://en.wikipedia.org/wiki/Isolation_(database_systems)#Isolation_levels.
const (
	LevelDefault IsolationLevel = iota
	LevelReadUncommitted
	LevelReadCommitted
	LevelWriteCommitted
	LevelRepeatableRead
	LevelSnapshot
	LevelSerializable
	LevelLinearizable
)

type TxOptions struct {
	// Isolation is the transaction isolation level.
	// If zero, the driver-specific default isolation level is used.
	Isolation IsolationLevel

	// ReadOnly is whether to set the transaction to read-only.
	ReadOnly bool
}

// TransactionBeginner provides functionality for starting a new transaction.
type TransactionBeginner interface {
	// Begin starts a new transaction and takes context and TxOptions as arguments.
	Begin(ctx context.Context, opts *TxOptions) (Tx, error)
}

// TransactionOption is a function that configures a TxOptions.
type TransactionOption func(options *TxOptions)

// WithIsolationLevel sets the transaction isolation level.
func WithIsolationLevel(level IsolationLevel) TransactionOption {
	return func(opts *TxOptions) {
		opts.Isolation = level
	}
}

// WithReadOnly sets the transaction to read-only.
func WithReadOnly() TransactionOption {
	return func(opts *TxOptions) {
		opts.ReadOnly = true
	}
}
