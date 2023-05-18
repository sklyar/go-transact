package txtest

import (
	"context"

	"github.com/sklyar/go-transact/internal/txcontext"
)

// WithContext returns a context with an embedded transaction context.
// The transaction context is created with default values.
func WithContext(ctx context.Context) context.Context {
	return WithContextValue(ctx, "1", false)
}

// WithChildContext returns a context with an embedded transaction context.
// The transaction context is created with default values and child flag set to true.
func WithChildContext(ctx context.Context) context.Context {
	return WithContextValue(ctx, "1", true)
}

// WithContextValue returns a context with an embedded transaction context.
// The transaction context is created with the provided id and child flag.
func WithContextValue(ctx context.Context, id string, child bool) context.Context {
	v := txcontext.Value{
		ID:    id,
		Child: child,
	}
	return txcontext.Wrap(ctx, v)
}
