package txcontext

import (
	"context"
)

// key is the key type for the context.
type key string

const (
	// transactKey is the key for the context.
	transactKey key = "transaction"
)

// Value is the value type for the context.
type Value struct {
	ID string

	Child bool
	Done  bool
}

// Wrap wraps the context with the transaction information.
func Wrap(ctx context.Context, value Value) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, transactKey, value)
}

// From returns the transaction information from the context.
// It returns false if the context doesn't have the transaction information.
func From(ctx context.Context) (Value, bool) {
	if ctx == nil {
		return Value{}, false
	}

	v, ok := ctx.Value(transactKey).(Value)
	return v, ok
}

// ID returns the transaction ID from the context.
// It returns false if the context doesn't have the transaction ID.
func ID(ctx context.Context) (string, bool) {
	v, ok := From(ctx)
	if !ok {
		return "", false
	}

	if v.Done {
		return "", false
	}

	return v.ID, true
}

// IsChild returns true if the context is a Child transaction.
func IsChild(ctx context.Context) bool {
	v, ok := From(ctx)
	return ok && v.Child
}

// WithTx adds a transaction information to the context.
// It returns a new context and a flag indicating whether the transaction is a Child.
func WithTx(ctx context.Context, nextID func() string) (context.Context, Value) {
	v, exists := From(ctx)
	if exists {
		v.Child = true
	} else {
		v.ID = nextID()
	}

	return Wrap(ctx, v), v
}
