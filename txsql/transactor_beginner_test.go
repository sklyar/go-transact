package txsql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithIsolationLevel(t *testing.T) {
	opts := new(TxOptions)
	WithIsolationLevel(LevelDefault)(opts)
	assert.Equal(t, LevelDefault, opts.Isolation)
}

func TestWithReadOnly(t *testing.T) {
	opts := new(TxOptions)
	WithReadOnly()(opts)
	assert.True(t, opts.ReadOnly)
}
