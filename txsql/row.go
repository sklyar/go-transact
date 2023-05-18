package txsql

// Row is the result of calling QueryRow to select a single row.
type Row interface {
	// Scan copies the columns from the matched row into the values
	// pointed at by dest. See the documentation on Rows.Scan for details.
	// If more than one row matches the query,
	// Scan uses the first row and discards the rest. If no row matches
	// the query, Scan returns ErrNoRows.
	Scan(dest ...any) error

	// Err provides a way for wrapping packages to check for
	// query errors without calling Scan.
	// Err returns the error, if any, that was encountered while running the query.
	// If this error is not nil, this error will also be returned from Scan.
	Err() error
}
