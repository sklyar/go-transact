package txsql

import (
	"context"
	"database/sql/driver"
)

// DB is a database handle representing a pool of zero or more
// underlying connections. It's safe for concurrent use by multiple
// goroutines.
type DB interface {
	DBHandler
	TransactionBeginner
	ConnManager

	// Ping verifies a connection to the database is still alive,
	// establishing a connection if necessary.
	Ping(ctx context.Context) error

	// Driver returns the database's underlying driver.
	Driver() driver.Driver
}
