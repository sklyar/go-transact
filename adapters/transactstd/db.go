package transactstd

import (
	"context"
	stdsql "database/sql"

	"github.com/sklyar/go-transact"
	"github.com/sklyar/go-transact/txsql"
)

// Database is a wrapper around stdsql.DB.
type Database struct {
	*stdsql.DB
	txs transact.TransactionStore
}

// Wrap creates new wrapper for stdsql.DB.
func Wrap(db *stdsql.DB) transact.AdapterFactoryFunc {
	return func(transactionStore transact.TransactionStore) (txsql.DB, error) {
		return &Database{
			DB:  db,
			txs: transactionStore,
		}, nil
	}
}

func (db *Database) Exec(ctx context.Context, query string, args ...any) (txsql.Result, error) {
	if tx, transacted := db.txs.Transaction(ctx); transacted {
		res, err := tx.Exec(ctx, query, args...)
		if err != nil {
			return nil, err
		}

		return newResult(res), nil
	}

	res, err := db.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return newResult(res), nil
}

func (db *Database) Query(ctx context.Context, query string, args ...any) (txsql.Rows, error) {
	if tx, transacted := db.txs.Transaction(ctx); transacted {
		rows, err := tx.Query(ctx, query, args...)
		if err != nil {
			return nil, err
		}
		return rows, nil
	}

	rows, err := db.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return newRows(rows), nil
}

func (db *Database) QueryRow(ctx context.Context, query string, args ...any) txsql.Row {
	if tx, transacted := db.txs.Transaction(ctx); transacted {
		return tx.QueryRow(ctx, query, args...)
	}

	row := db.DB.QueryRowContext(ctx, query, args...)
	return newRow(row, nil)
}

func (db *Database) Prepare(ctx context.Context, query string) (txsql.Stmt, error) {
	if tx, transacted := db.txs.Transaction(ctx); transacted {
		return tx.Prepare(ctx, query)
	}

	stmt, err := db.DB.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	return newStmt(stmt), nil
}

func (db *Database) Begin(ctx context.Context, opts *txsql.TxOptions) (txsql.Tx, error) {
	var stdOpts *stdsql.TxOptions
	if opts != nil {
		stdOpts = &stdsql.TxOptions{
			Isolation: stdsql.IsolationLevel(opts.Isolation),
			ReadOnly:  opts.ReadOnly,
		}
	}

	sqlTx, err := db.DB.BeginTx(ctx, stdOpts)
	if err != nil {
		return nil, err
	}

	return &tx{Tx: sqlTx}, nil
}

func (db *Database) Ping(ctx context.Context) error {
	return db.DB.PingContext(ctx)
}
