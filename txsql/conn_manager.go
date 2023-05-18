package txsql

import (
	"context"
	"database/sql"
	"time"
)

// ConnManager provides functionality for managing database connections.
type ConnManager interface {
	// Close closes the database and prevents new queries from starting.
	Close() error

	// Conn returns a single connection by either opening a new connection
	// or returning an existing connection from the connection pool.
	Conn(ctx context.Context) (*sql.Conn, error)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	SetMaxOpenConns(n int)

	// SetMaxIdleConns sets the maximum number of connections in the idle
	// connection pool.
	SetMaxIdleConns(n int)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	SetConnMaxLifetime(d time.Duration)

	// SetConnMaxIdleTime sets the maximum amount of time a connection may be idle
	// before it is closed.
	SetConnMaxIdleTime(d time.Duration)
}
