package transact

import (
	"context"
	"errors"
	"sync"

	"github.com/sklyar/go-transact/internal/txcontext"
)

// store is the store for transactions.
type store struct {
	txs map[string]*Transaction
	mu  sync.RWMutex
}

func newStore() *store {
	return &store{
		txs: make(map[string]*Transaction),
	}
}

// Transaction returns the transaction for the given context.
// If there is no transaction in the context, it returns false.
func (s *store) Transaction(ctx context.Context) (*Transaction, bool) {
	tid, ok := txcontext.ID(ctx)
	if !ok {
		return nil, false
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.txs[tid], true
}

// Add adds the transaction to the store.
func (s *store) Add(tx *Transaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tid := tx.ID()
	if _, ok := s.txs[tid]; ok {
		return errors.New("transaction already exists")
	}

	s.txs[tid] = tx

	return nil
}

// Delete deletes the transaction from the store.
func (s *store) Delete(ctx context.Context, tx *Transaction) error {
	if txcontext.IsChild(ctx) {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	tid := tx.ID()
	if _, ok := s.txs[tid]; !ok {
		return errors.New("transaction not found")
	}

	delete(s.txs, tid)

	return nil
}

// Len returns the number of transactions in the store.
func (s *store) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.txs)
}
