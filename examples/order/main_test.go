package main

import (
	"context"
	"errors"
	"testing"

	"github.com/sklyar/go-transact"
	"github.com/sklyar/go-transact/examples/order/mock"
	"github.com/sklyar/go-transact/txsql"
	"github.com/sklyar/go-transact/txtest"
	"github.com/stretchr/testify/assert"
)

func TestApp_CreateOrder(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	type mocks struct {
		db *txtest.DB
		tx *txtest.Tx

		orderRepo     *mock.OrderRepository
		inventoryRepo *mock.InventoryRepository
	}

	type in struct {
		customerID int
		products   []int
	}

	tests := []struct {
		name string

		in in

		setup   func(in in, m *mocks)
		wantErr string
	}{
		{
			name: "create order",
			in:   in{customerID: 1, products: []int{1, 2}},
			setup: func(in in, m *mocks) {
				ctx := txtest.WithContext(ctx)

				m.db.On("Begin", ctx, (*txsql.TxOptions)(nil)).Return(m.tx, nil)
				m.orderRepo.On("Create", ctx, 1).Return(1, nil)

				for _, productID := range in.products {
					m.inventoryRepo.On("GetProductQuantity", ctx, productID).Return(3, nil)
					m.orderRepo.On("AddProduct", ctx, 1, productID).Return(nil)
					m.inventoryRepo.On("DecrementProductQuantity", ctx, productID).Return(nil)
				}

				m.tx.On("Commit", ctx).Return(nil)
			},
		},
		{
			name: "create order with rollback",
			in:   in{customerID: 1, products: []int{1}},
			setup: func(in in, m *mocks) {
				ctx := txtest.WithContext(ctx)

				m.db.On("Begin", ctx, (*txsql.TxOptions)(nil)).Return(m.tx, nil)
				m.orderRepo.On("Create", ctx, 1).Return(1, nil)

				productID := in.products[0]
				m.inventoryRepo.On("GetProductQuantity", ctx, productID).Return(3, nil)
				m.orderRepo.On("AddProduct", ctx, 1, productID).Return(nil)
				m.inventoryRepo.On("DecrementProductQuantity", ctx, productID).Return(errors.New("some error"))

				m.tx.On("Rollback", ctx).Return(nil)
			},
			wantErr: "some error",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := mocks{
				db:            txtest.NewDB(t),
				tx:            txtest.NewTx(t),
				orderRepo:     mock.NewOrderRepository(t),
				inventoryRepo: mock.NewInventoryRepository(t),
			}
			if tt.setup != nil {
				tt.setup(tt.in, &mocks)
			}

			adapter := func(_ transact.TransactionStore) (txsql.DB, error) { return mocks.db, nil }
			txManager, _, err := transact.NewManager(adapter)
			assert.NoError(t, err)

			orderService := NewOrderService(txManager, mocks.orderRepo, mocks.inventoryRepo)
			orderID, err := orderService.Create(context.Background(), tt.in.customerID, tt.in.products)
			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, 1, orderID)
		})
	}
}
