package memory

import (
	"context"
	"sync"
)

type Repository struct {
	Models *ModelsRepository
}

// NewRepository creates a new Repository
func NewRepository() (*Repository, error) {
	return &Repository{
		Models: &ModelsRepository{
			models: make(map[string]*Model),
			mu:     new(sync.RWMutex),
		},
	}, nil
}

// Connect initializes the repository
func (r *Repository) Connect(ctx context.Context) error {
	return nil
}

// Disconnect stops the repository
func (r *Repository) Disconnect(ctx context.Context) error {
	return nil
}
