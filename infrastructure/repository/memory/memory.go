package memory

import (
	"context"
	"sync"

	"birus/domain/entity/shingling/classifier"
)

type Repository struct {
	ClassifierRepository *ClassifierRepository
}

// NewRepository creates a new Repository
func NewRepository() (*Repository, error) {
	return &Repository{
		ClassifierRepository: &ClassifierRepository{
			classifiers: make(map[string]*classifier.Classifier),
			mu:          new(sync.RWMutex),
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
