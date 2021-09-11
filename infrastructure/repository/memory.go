package repository

import (
	"context"
	"sync"
)

type MemoryRepository struct {
	Models *ModelsMemoryRepository
}

// NewMemoryRepository creates a new MemoryRepository
func NewMemoryRepository() (*MemoryRepository, error) {
	return &MemoryRepository{
		Models: &ModelsMemoryRepository{
			models: make(map[string]*Model),
			mu:     new(sync.RWMutex),
		},
	}, nil
}

// Connect initializes the repository
func (r *MemoryRepository) Connect(ctx context.Context) error {
	return nil
}

// Disconnect stops the repository
func (r *MemoryRepository) Disconnect(ctx context.Context) error {
	return nil
}
