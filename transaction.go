package transact

import (
	"context"
	"errors"

	"github.com/sklyar/go-transact/internal/txcontext"
	"github.com/sklyar/go-transact/txsql"
)

var (
	ErrNoTransaction     = errors.New("no transaction")
	ErrClosedTransaction = errors.New("transaction is closed")

	errCommittedTransaction = errors.New("operation failed: transaction has already been committed")
	errMarkedForRollback    = errors.New("operation failed: transaction has been marked for rollback and cannot be committed")
)

// Transaction is a transaction wrapper.
type Transaction struct {
	txsql.Tx

	id string

	commit   bool
	rollback bool
}

// newTransaction creates a new transaction.
func newTransaction(id string, tx txsql.Tx) *Transaction {
	return &Transaction{Tx: tx, id: id}
}

// Commit executes a transaction.
// If the transaction is a child, it does nothing and the original context is returned.
// If the transaction has already been committed or has been marked for deletion,
// it returns the original context along with the corresponding error (ErrNoTransaction or ErrClosedTransaction).
// After a successful commit, the transaction is marked as done within the context.
func (tx *Transaction) Commit(ctx context.Context) (context.Context, error) {
	if txcontext.IsChild(ctx) {
		return ctx, nil
	}

	v, exists := txcontext.From(ctx)
	if !exists {
		return ctx, ErrNoTransaction
	}
	if v.Done {
		return ctx, ErrClosedTransaction
	}

	if tx.commit {
		// unexpected commit after commit.
		return ctx, errCommittedTransaction
	}
	if tx.rollback {
		// unexpected commit after rollback.
		return ctx, errMarkedForRollback
	}

	tx.commit = true
	v.Done = true
	return txcontext.Wrap(ctx, v), tx.Tx.Commit(ctx)
}

// Rollback aborts a transaction.
// If the transaction is a child, it does nothing and the original context is returned.
// If the transaction doesn't exist in the context, it returns the original context along with ErrNoTransaction.
// If the transaction has already been rolled back or marked as done,
// it returns the original context along with ErrClosedTransaction.
// Upon a successful rollback, the transaction is marked as done within the context.
func (tx *Transaction) Rollback(ctx context.Context) (context.Context, error) {
	v, exists := txcontext.From(ctx)
	if !exists {
		return ctx, ErrNoTransaction
	}
	if v.Done {
		return ctx, ErrClosedTransaction
	}

	if tx.commit {
		// unexpected commit after commit.
		return ctx, errCommittedTransaction
	}

	fn := tx.Tx.Rollback
	if tx.rollback {
		fn = func(_ context.Context) error { return nil }
	} else {
		tx.rollback = true
	}

	v.Done = true
	return txcontext.Wrap(ctx, v), fn(ctx)
}

// ID returns a transaction ID.
func (tx *Transaction) ID() string {
	return tx.id
}
