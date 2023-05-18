package transact

import (
	"context"
	"errors"
	"testing"

	"github.com/sklyar/go-transact/internal/txcontext"
	"github.com/sklyar/go-transact/txtest"
	"github.com/stretchr/testify/assert"
)

func TestTransaction_Commit(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	txContext := txtest.WithContext(ctx)
	childTxContext := txtest.WithChildContext(txContext)

	errTest := errors.New("test error")

	tests := []struct {
		name string
		ctx  context.Context
		tx   *Transaction

		setup func(sqlTx *txtest.Tx)

		// returned params
		wantCtx context.Context
		wantErr error

		// modified params
		wantCommit bool
		wantDone   bool
	}{
		{
			name: "commit succeeds",
			ctx:  txContext,
			tx:   &Transaction{id: "id"},
			setup: func(sqlTx *txtest.Tx) {
				sqlTx.On("Commit", txContext).Return(nil)
			},
			wantCtx:    setContextAsDone(t, txContext),
			wantCommit: true,
		},
		{
			name:    "no transaction",
			ctx:     ctx,
			tx:      &Transaction{id: "id"},
			wantCtx: ctx,
			wantErr: ErrNoTransaction,
		},
		{
			name:    "transaction is done",
			ctx:     setContextAsDone(t, txContext),
			tx:      &Transaction{id: "id"},
			wantCtx: setContextAsDone(t, txContext),
			wantErr: ErrClosedTransaction,
		},
		{
			name:    "transaction is a child",
			ctx:     childTxContext,
			tx:      &Transaction{id: "id"},
			wantCtx: childTxContext,
			wantErr: nil,
		},
		{
			name:    "transaction is committed",
			ctx:     txContext,
			tx:      &Transaction{id: "id", commit: true},
			wantCtx: txContext,
			wantErr: errCommittedTransaction,
		},
		{
			name:    "transaction is marked for rollback",
			ctx:     txContext,
			tx:      &Transaction{id: "id", rollback: true},
			wantCtx: txContext,
			wantErr: errMarkedForRollback,
		},
		{
			name: "commit fails",
			ctx:  txContext,
			tx:   &Transaction{id: "id"},
			setup: func(sqlTx *txtest.Tx) {
				sqlTx.On("Commit", txContext).Return(errTest)
			},
			// if commit fails, we still mark transaction as done.
			wantCtx: setContextAsDone(t, txContext),
			wantErr: errTest,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sqlTx := txtest.NewTx(t)
			if tt.setup != nil {
				tt.setup(sqlTx)
			}

			tx := tt.tx
			tx.Tx = sqlTx

			ctx, err := tx.Commit(tt.ctx)
			if err != nil {
				assert.Equal(t, tt.wantErr, err)
				assert.Equal(t, tt.wantCtx, ctx)
				return
			}

			txContextValue, ok := txcontext.From(ctx)
			assert.True(t, ok)

			expTxContextValue, ok := txcontext.From(tt.wantCtx)
			assert.True(t, ok)

			assert.Equal(t, expTxContextValue, txContextValue)
			assert.Equal(t, tt.wantCommit, tx.commit)
		})
	}
}

func TestTransaction_Rollback(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	txContext := txtest.WithContext(ctx)
	childTxContext := txtest.WithChildContext(txContext)

	errTest := errors.New("test error")

	tests := []struct {
		name string
		ctx  context.Context
		tx   *Transaction

		setup func(sqlTx *txtest.Tx)

		// returned params
		wantCtx context.Context
		wantErr error

		// modified params
		wantRollback bool
		wantDone     bool
	}{
		{
			name: "rollback succeeds",
			ctx:  txContext,
			tx:   &Transaction{id: "id"},
			setup: func(sqlTx *txtest.Tx) {
				sqlTx.On("Rollback", txContext).Return(nil)
			},
			wantCtx:      setContextAsDone(t, txContext),
			wantDone:     true,
			wantRollback: true,
		},
		{
			name: "rollback succeeds for child transaction",
			ctx:  childTxContext,
			tx:   &Transaction{id: "id"},
			setup: func(sqlTx *txtest.Tx) {
				sqlTx.On("Rollback", childTxContext).Return(nil)
			},
			wantCtx:      setContextAsDone(t, childTxContext),
			wantDone:     true,
			wantRollback: true,
		},
		{
			name:    "no transaction",
			ctx:     ctx,
			tx:      &Transaction{id: "id"},
			wantCtx: ctx,
			wantErr: ErrNoTransaction,
		},
		{
			name:    "transaction is done",
			ctx:     setContextAsDone(t, txContext),
			tx:      &Transaction{id: "id"},
			wantCtx: setContextAsDone(t, txContext),
			wantErr: ErrClosedTransaction,
		},
		{
			name: "rollback fails",
			ctx:  txContext,
			tx:   &Transaction{id: "id"},
			setup: func(sqlTx *txtest.Tx) {
				sqlTx.On("Rollback", txContext).Return(errTest)
			},
			wantCtx:      setContextAsDone(t, txContext),
			wantErr:      errTest,
			wantDone:     true,
			wantRollback: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sqlTx := txtest.NewTx(t)
			if tt.setup != nil {
				tt.setup(sqlTx)
			}

			tx := tt.tx
			tx.Tx = sqlTx

			ctx, err := tx.Rollback(tt.ctx)
			if err != nil {
				assert.Equal(t, tt.wantErr, err)
				assert.Equal(t, tt.wantCtx, ctx)
				return
			}

			txContextValue, ok := txcontext.From(ctx)
			assert.True(t, ok)

			expTxContextValue, ok := txcontext.From(tt.wantCtx)
			assert.True(t, ok)

			assert.Equal(t, expTxContextValue, txContextValue)
			assert.Equal(t, tt.wantRollback, tx.rollback)
		})
	}
}

func setContextAsDone(t *testing.T, ctx context.Context) context.Context {
	t.Helper()

	v, exists := txcontext.From(ctx)
	assert.True(t, exists)

	v.Done = true
	return txcontext.Wrap(ctx, v)
}
