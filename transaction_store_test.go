package transact

import (
	"context"
	"testing"

	"github.com/sklyar/go-transact/txtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_transactionStore_Transaction(t *testing.T) {
	t.Parallel()

	s := newTransactionStore()
	require.NoError(t, s.Add(&Transaction{id: "1"}), "Expected no error when adding transaction")
	require.NoError(t, s.Add(&Transaction{id: "2"}), "Expected no error when adding transaction")

	t.Run("nil context", func(t *testing.T) {
		tx, transacted := s.Transaction(nil) //nolint:staticcheck
		assert.False(t, transacted, "Expected no transaction")
		assert.Nil(t, tx, "Expected no transaction")
	})

	t.Run("no transaction", func(t *testing.T) {
		tx, transacted := s.Transaction(context.Background())
		assert.False(t, transacted, "Expected no transaction")
		assert.Nil(t, tx, "Expected no transaction")
	})

	t.Run("has transaction", func(t *testing.T) {
		ctx := txtest.WithContext(context.Background())
		tx, transacted := s.Transaction(ctx)
		assert.True(t, transacted, "Expected transaction")
		assert.Equal(t, "1", tx.ID(), "Expected transaction ID to be 1")
	})

	t.Run("nested transaction", func(t *testing.T) {
		ctx := txtest.WithContextValue(context.Background(), "1", true)
		tx, transacted := s.Transaction(ctx)
		assert.True(t, transacted, "Expected transaction")
		// in this case, the transaction store should return the parent transaction.
		assert.Equal(t, "1", tx.ID(), "Expected transaction ID to be 1")
	})
}

func Test_transactionStore_Add(t *testing.T) {
	t.Parallel()

	s := newTransactionStore()
	assert.NoError(t, s.Add(&Transaction{id: "1"}), "Expected no error when adding transaction")
	assert.NoError(t, s.Add(&Transaction{id: "2"}), "Expected no error when adding transaction")

	assert.Equal(t, 2, len(s.txs), "Expected 2 transactions")

	assert.Error(t, s.Add(&Transaction{id: "1"}), "Expected error when adding transaction with same id")
}

func Test_transactionStore_Delete(t *testing.T) {
	t.Parallel()

	t.Run("delete a transaction", func(t *testing.T) {
		s := newTransactionStore()
		tx1 := &Transaction{id: "1"}
		tx2 := &Transaction{id: "2"}

		require.NoError(t, s.Add(tx1), "Expected no error when adding transaction")
		require.NoError(t, s.Add(tx2), "Expected no error when adding transaction")
		require.Equal(t, 2, len(s.txs), "Expected 2 transactions")

		baseContext := context.Background()

		ctx1 := txtest.WithContextValue(baseContext, "1", false)
		assert.NoError(t, s.Delete(ctx1, tx1), "Expected no error when deleting transaction")

		ctx2 := txtest.WithContextValue(baseContext, "2", false)
		assert.NoError(t, s.Delete(ctx2, tx2), "Expected no error when deleting transaction")

		assert.Equal(t, 0, len(s.txs), "Expected 0 transactions")
	})

	t.Run("delete a nested transaction", func(t *testing.T) {
		s := newTransactionStore()
		tx1 := &Transaction{id: "1"}
		tx2 := &Transaction{id: "2"}

		require.NoError(t, s.Add(tx1), "Expected no error when adding transaction")
		require.NoError(t, s.Add(tx2), "Expected no error when adding transaction")
		require.Equal(t, 2, len(s.txs), "Expected 2 transactions")

		baseContext := context.Background()
		childTxContext := txtest.WithContextValue(baseContext, "2", true)

		// In this case, the transaction store shouldn`t delete the transaction.
		// Because the transaction is a child transaction.
		assert.NoError(t, s.Delete(childTxContext, tx2), "Expected no error when deleting transaction")
		assert.Len(t, s.txs, 2, "Expected 2 transactions")
	})
}

func Test_transactionStore_Len(t *testing.T) {
	store := newTransactionStore()
	assert.Equal(t, 0, store.Len(), "Expected 0 transactions")

	_ = store.Add(&Transaction{id: "1"})
	assert.Equal(t, 1, store.Len(), "Expected 1 transaction")

	_ = store.Add(&Transaction{id: "2"})
	assert.Equal(t, 2, store.Len(), "Expected 2 transactions")

	_ = store.Delete(context.Background(), &Transaction{id: "1"})
	assert.Equal(t, 1, store.Len(), "Expected 1 transaction")
}
