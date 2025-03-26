package inmemory

import (
	"context"
	"sync"
)

type InMemoryRepository struct {
	mu    sync.RWMutex
	users map[int64]struct{}
}

func New() *InMemoryRepository {
	return &InMemoryRepository{
		users: make(map[int64]struct{}),
	}
}

func (r *InMemoryRepository) CreateUser(ctx context.Context, userID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[userID] = struct{}{}
	return nil
}

func (r *InMemoryRepository) UserExists(ctx context.Context, userID int64) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.users[userID]
	return exists, nil
}
