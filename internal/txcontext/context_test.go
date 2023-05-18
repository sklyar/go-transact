package txcontext

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapContext(t *testing.T) {
	ctx := context.Background()
	value := Value{
		ID:    "test",
		Child: true,
	}

	t.Run("with non-nil context", func(t *testing.T) {
		newCtx := Wrap(ctx, value)

		v, ok := newCtx.Value(transactKey).(Value)
		assert.True(t, ok, "Expected value to be present in context")
		assert.Equal(t, value, v, "Expected value to be equal to inserted value")
	})

	t.Run("with nil context", func(t *testing.T) {
		newCtx := Wrap(nil, value) //nolint:staticcheck

		v, ok := newCtx.Value(transactKey).(Value)
		assert.True(t, ok, "Expected value to be present in context")
		assert.Equal(t, value, v, "Expected value to be equal to inserted value")
	})
}

func TestFromContext(t *testing.T) {
	ctx := context.Background()
	value := Value{
		ID:    "test",
		Child: true,
	}

	t.Run("value in context", func(t *testing.T) {
		newCtx := Wrap(ctx, value)

		v, ok := From(newCtx)
		assert.True(t, ok, "Expected value to be present in context")
		assert.Equal(t, value, v, "Expected value to be equal to inserted value")
	})

	t.Run("no value in context", func(t *testing.T) {
		v, ok := From(ctx)
		assert.False(t, ok, "Expected no value to be present in context")
		assert.Equal(t, Value{}, v, "Expected value to be default value")
	})
}

func TestIDFromContext(t *testing.T) {
	ctx := context.Background()
	value := Value{
		ID:    "test",
		Child: true,
	}

	t.Run("ID present in context", func(t *testing.T) {
		newCtx := Wrap(ctx, value)

		id, ok := ID(newCtx)
		assert.True(t, ok, "Expected ID to be present in context")
		assert.Equal(t, value.ID, id, "Expected ID to be equal to inserted ID")
	})

	t.Run("ID not present in context", func(t *testing.T) {
		id, ok := ID(ctx)
		assert.False(t, ok, "Expected no ID to be present in context")
		assert.Equal(t, "", id, "Expected ID to be empty string")
	})
}

func TestIsChildFromContext(t *testing.T) {
	t.Run("Child present in context", func(t *testing.T) {
		ctx := context.Background()
		value := Value{
			ID:    "test",
			Child: true,
		}
		newCtx := Wrap(ctx, value)

		isChild := IsChild(newCtx)
		assert.True(t, isChild, "Expected Child to be true when inserted value is true")
	})

	t.Run("transaction not present in context", func(t *testing.T) {
		isChild := IsChild(context.Background())
		assert.False(t, isChild, "Expected Child to be false when no value is present")
	})

	t.Run("Child not present in context", func(t *testing.T) {
		ctx := context.Background()
		value := Value{
			ID:    "test",
			Child: false,
		}
		newCtx := Wrap(ctx, value)

		isChild := IsChild(newCtx)
		assert.False(t, isChild, "Expected Child to be false when inserted value is false")
	})
}

func TestAddTxToContext(t *testing.T) {
	t.Parallel()

	emptyContext := context.Background()

	tests := []struct {
		name   string
		ctx    context.Context
		nextID string

		want         context.Context
		wantCtxValue Value
	}{
		{
			name:         "add transaction info to empty context",
			ctx:          emptyContext,
			nextID:       "10",
			want:         Wrap(emptyContext, Value{ID: "10", Child: false}),
			wantCtxValue: Value{ID: "10", Child: false},
		},
		{
			name:         "add Child transaction info to transaction context",
			ctx:          Wrap(emptyContext, Value{ID: "10", Child: false}),
			want:         Wrap(emptyContext, Value{ID: "10", Child: true}),
			wantCtxValue: Value{ID: "10", Child: true},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, txContextValue := WithTx(tt.ctx, func() string { return tt.nextID })
			assertContext(t, ctx, tt.want)
			assert.Equal(t, tt.wantCtxValue, txContextValue)
		})
	}
}

func assertContext(t *testing.T, ctx, expCtx context.Context) {
	t.Helper()

	v1, ok1 := From(ctx)
	v2, ok2 := From(expCtx)

	assert.Equal(t, ok1, ok2, "Expected both contexts to have same value")
	assert.Equal(t, v1, v2, "Expected both contexts to have same value")
}
