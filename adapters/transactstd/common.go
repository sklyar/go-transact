package transactstd

import (
	stdsql "database/sql"
	"errors"

	"github.com/sklyar/go-transact/txsql"
)

// Row implements txsql.Row interface.
type Row struct {
	row *stdsql.Row
	err error
}

// newRow creates new Row.
func newRow(row *stdsql.Row, err error) txsql.Row {
	return &Row{row: row, err: err}
}

// Scan implements txsql.Row interface.
func (r *Row) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}
	return r.row.Scan(dest...)
}

// Close implements txsql.Row interface.
func (r *Row) Close() error {
	if r.err != nil {
		return r.err
	}
	return r.row.Err()
}

// Err implements txsql.Row interface.
func (r *Row) Err() error {
	return errors.Join(r.err, r.row.Err())
}

// Rows implements txsql.Rows interface.
type Rows struct {
	*stdsql.Rows
}

// newRows creates new Rows.
func newRows(rows *stdsql.Rows) *Rows {
	return &Rows{Rows: rows}
}

// Result implements txsql.Result interface.
type Result struct {
	stdsql.Result
}

// newResult creates new Result.
func newResult(result stdsql.Result) *Result {
	return &Result{Result: result}
}

// Stmt implements txsql.Stmt interface.
type Stmt struct {
	*stdsql.Stmt
}

// newStmt creates new Stmt.
func newStmt(stmt *stdsql.Stmt) *Stmt {
	return &Stmt{Stmt: stmt}
}

// Exec implements txsql.Stmt interface.
func (s *Stmt) Exec(args ...any) (txsql.Result, error) {
	res, err := s.Stmt.Exec(args...)
	if err != nil {
		return nil, err
	}

	return newResult(res), nil
}

// Query implements txsql.Stmt interface.
func (s *Stmt) Query(args ...any) (txsql.Rows, error) {
	rows, err := s.Stmt.Query(args...)
	if err != nil {
		return nil, err
	}

	return newRows(rows), nil
}

// QueryRow implements txsql.Stmt interface.
func (s *Stmt) QueryRow(args ...any) txsql.Row {
	row := s.Stmt.QueryRow(args...)
	return newRow(row, nil)
}
