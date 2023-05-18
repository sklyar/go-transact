package txsql

import "context"

// Tx represents a database transaction.
type Tx interface {
	DBHandler

	// Commit commits the transaction.
	Commit(ctx context.Context) error

	// Rollback aborts the transaction.
	Rollback(ctx context.Context) error

	// Stmt returns a transaction-specific prepared statement
	// from an existing statement.
	Stmt(stmt Stmt) Stmt
}
