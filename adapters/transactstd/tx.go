package transactstd

import (
	"context"
	"database/sql"

	"github.com/sklyar/go-transact/txsql"
)

type tx struct {
	*sql.Tx
}

func (t *tx) Exec(ctx context.Context, query string, args ...any) (txsql.Result, error) {
	return t.Tx.ExecContext(ctx, query, args...)
}

func (t *tx) Query(ctx context.Context, query string, args ...any) (txsql.Rows, error) {
	rows, err := t.Tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return newRows(rows), nil
}

func (t *tx) QueryRow(ctx context.Context, query string, args ...any) txsql.Row {
	row := t.Tx.QueryRowContext(ctx, query, args...)
	return newRow(row, nil)
}

func (t *tx) Prepare(ctx context.Context, query string) (txsql.Stmt, error) {
	stmt, err := t.Tx.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	return newStmt(stmt), nil
}

func (t *tx) Commit(_ context.Context) error {
	return t.Tx.Commit()
}

func (t *tx) Rollback(_ context.Context) error {
	return t.Tx.Rollback()
}

func (t *tx) Stmt(stmt txsql.Stmt) txsql.Stmt {
	return newStmt(t.Tx.Stmt(stmt.(*Stmt).Stmt))
}
