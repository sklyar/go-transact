package transact

import (
	"context"
	"errors"
	"testing"

	"github.com/sklyar/go-transact/txsql"
	"github.com/sklyar/go-transact/txtest"
	"github.com/stretchr/testify/assert"
)

var nilTxOptions = (*txsql.TxOptions)(nil)

func TestNewManager(t *testing.T) {
	t.Run("successful adapter factory invocation", func(t *testing.T) {
		mockDB := txtest.NewDB(t)
		adapterFactory := func(transactionStore TransactionStore) (txsql.DB, error) {
			return mockDB, nil
		}

		manager, db, err := NewManager(adapterFactory)
		assert.NoError(t, err)
		assert.NotNil(t, manager)
		assert.Equal(t, mockDB, db)
	})

	t.Run("adapter factory returns error", func(t *testing.T) {
		expectedError := errors.New("mock error")
		adapterFactory := func(transactionStore TransactionStore) (txsql.DB, error) {
			return nil, expectedError
		}

		manager, db, err := NewManager(adapterFactory)
		assert.Error(t, err)
		assert.Nil(t, manager)
		assert.Nil(t, db)
		assert.Equal(t, expectedError, err)
	})
}

func TestBeginFunc(t *testing.T) {
	db := txtest.NewDB(t)
	manager := &Manager{
		db:    db,
		store: newTransactionStore(),
	}

	baseContext := context.Background()
	txContext := txtest.WithContext(baseContext)

	tx := txtest.NewTx(t)
	tx.On("Commit", txContext).Return(nil)
	db.On("Begin", txContext, nilTxOptions).Return(tx, nil)

	err := manager.BeginFunc(baseContext, func(_ context.Context) error { return nil })
	assert.NoError(t, err)
	assert.Equal(t, 0, manager.store.Len())
}

func TestBeginFuncTransactionFunctionReturnsError(t *testing.T) {
	db := txtest.NewDB(t)
	manager := &Manager{
		db:    db,
		store: newTransactionStore(),
	}

	baseContext := context.Background()
	txContext := txtest.WithContext(baseContext)

	tx := txtest.NewTx(t)
	tx.On("Rollback", txContext).Return(nil)
	db.On("Begin", txContext, nilTxOptions).Return(tx, nil)

	someErr := errors.New("some error")
	err := manager.BeginFunc(baseContext, func(tx context.Context) error {
		return someErr
	})
	assert.ErrorContains(t, err, someErr.Error())
	assert.Equal(t, 0, manager.store.Len())
}

func TestBeginFuncSuccessfulTransactionWithChildTransaction(t *testing.T) {
	db := txtest.NewDB(t)
	manager := &Manager{
		db:    db,
		store: newTransactionStore(),
	}

	ctx := context.Background()
	txContext := txtest.WithContext(ctx)
	childTxContext := txtest.WithChildContext(txContext)

	tx := txtest.NewTx(t)
	db.On("Begin", txContext, nilTxOptions).Return(tx, nil)
	db.On("Query", txContext, "SELECT 1").Return(nil, nil)
	db.On("Query", childTxContext, "SELECT 2").Return(nil, nil)
	tx.On("Commit", txContext).Return(nil)

	err := manager.BeginFunc(ctx, func(tx context.Context) error {
		_, err := db.Query(tx, "SELECT 1")
		assert.NoError(t, err)

		return manager.BeginFunc(tx, func(tx context.Context) error {
			_, err := db.Query(tx, "SELECT 2")
			return err
		})
	})
	assert.NoError(t, err)
	assert.Equal(t, 0, manager.store.Len())
}

func TestBeginFuncErrorOnChildTransactionRetrieval(t *testing.T) {
	db := txtest.NewDB(t)
	manager := &Manager{
		db:    db,
		store: newTransactionStore(),
	}

	ctx := context.Background()
	txContext := txtest.WithContext(ctx)
	childTxContext := txtest.WithChildContext(txContext)

	someErr := errors.New("some error")

	tx := txtest.NewTx(t)
	db.On("Begin", txContext, nilTxOptions).Return(tx, nil)
	db.On("Query", txContext, "SELECT 1").Return(nil, nil)
	db.On("Query", childTxContext, "SELECT 2").Return(nil, someErr)
	tx.On("Rollback", childTxContext).Return(nil)

	err := manager.BeginFunc(ctx, func(tx context.Context) error {
		_, err := db.Query(tx, "SELECT 1")
		assert.NoError(t, err)

		return manager.BeginFunc(tx, func(tx context.Context) error {
			_, err := db.Query(tx, "SELECT 2")
			return err
		})
	})
	assert.ErrorContains(t, err, someErr.Error())
	assert.Equal(t, 0, manager.store.Len())
}
