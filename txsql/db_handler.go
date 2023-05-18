package txsql

import "context"

// DBHandler combines methods for executing SQL queries and preparing statements
type DBHandler interface {
	// Exec executes a query without returning any rows.
	// The args are for any placeholder parameters in the query.
	Exec(ctx context.Context, query string, args ...any) (Result, error)

	// Query executes a query that returns rows, typically a SELECT.
	// The args are for any placeholder parameters in the query.
	Query(ctx context.Context, query string, args ...any) (Rows, error)

	// QueryRow executes a query that is expected to return at most one row.
	// QueryRow always returns a non-nil value. Errors are deferred until
	// Row's Scan method is called. If the query selects no rows, the *Row's Scan will
	// return ErrNoRows. Otherwise, the *Row's Scan scans the first selected
	// row and discards the rest.
	QueryRow(ctx context.Context, query string, args ...any) Row

	// Prepare creates a prepared statement for later queries or executions.
	// Multiple queries or executions may be run concurrently from the
	// returned statement.
	Prepare(ctx context.Context, query string) (Stmt, error)
}
