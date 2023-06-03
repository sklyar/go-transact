package transact

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync/atomic"

	"github.com/sklyar/go-transact/internal/txcontext"
	"github.com/sklyar/go-transact/txsql"
)

type TransactionStore interface {
	// Transaction returns the transaction for the given context.
	// If there is no transaction in the context, it returns false.
	Transaction(ctx context.Context) (*Transaction, bool)
}

type AdapterFactoryFunc func(transactionStore TransactionStore) (txsql.DB, error)

type TransactionFunc func(tx context.Context) error

type Manager struct {
	db    txsql.DB
	store *store

	// lastID is the last transaction id.
	// It is used to generate a new transaction id.
	lastID uint64
}

func NewManager(adapterFactory AdapterFactoryFunc) (*Manager, txsql.DB, error) {
	store := newStore()
	db, err := adapterFactory(store)
	if err != nil {
		return nil, nil, err
	}

	return &Manager{db: db, store: store}, db, nil
}

func (m *Manager) BeginFunc(ctx context.Context, fn TransactionFunc, opts ...txsql.TransactionOption) (err error) {
	ctx, tx, err := m.transaction(ctx, opts)
	if err != nil {
		return err
	}

	defer func() {
		if derr := m.store.Delete(ctx, tx); derr != nil {
			err = errors.Join(err, fmt.Errorf("failed to delete transaction from context: %w", derr))
			return
		}
	}()

	if err = fn(ctx); err != nil {
		err = fmt.Errorf("failed to execute transaction function: %w", err)
		if _, rerr := tx.Rollback(ctx); rerr != nil {
			err = errors.Join(err, fmt.Errorf("failed to rollback transaction: %w", rerr))
		}
		return err
	}

	if _, err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (m *Manager) Begin(ctx context.Context, opts ...txsql.TransactionOption) (context.Context, *Transaction, error) {
	return m.transaction(ctx, opts)
}

func (m *Manager) transaction(ctx context.Context, opts []txsql.TransactionOption) (context.Context, *Transaction, error) {
	ctx, ctxVal := txcontext.WithTx(ctx, m.nextID)
	if ctxVal.Done {
		return nil, nil, errors.New("transaction already done")
	}
	if ctxVal.Child {
		tx, transacted := m.store.Transaction(ctx)
		if !transacted {
			return nil, nil, errors.New("failed to find parent transaction")
		}

		return ctx, tx, nil
	}

	var txOptions *txsql.TxOptions
	if len(opts) > 0 {
		txOptions = new(txsql.TxOptions)
	}
	for _, opt := range opts {
		opt(txOptions)
	}

	sqlTx, err := m.db.Begin(ctx, txOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	tid, _ := txcontext.ID(ctx)
	tx := newTransaction(tid, sqlTx)

	if err := m.store.Add(tx); err != nil {
		addErr := fmt.Errorf("failed to add transaction: %w", err)
		ctx, err := tx.Rollback(ctx)
		if err != nil {
			return ctx, nil, errors.Join(addErr, fmt.Errorf("failed to rollback transaction: %w", err))
		}
		return ctx, nil, err
	}

	return ctx, tx, nil
}

// nextID returns the next transaction id.
func (m *Manager) nextID() string {
	id := atomic.AddUint64(&m.lastID, 1)
	return strconv.FormatUint(id, 10)
}
