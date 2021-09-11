package repository

import (
	"context"
	"sync"

	"github.com/flavioltonon/birus/domain/entity"
	"github.com/flavioltonon/birus/internal/shingling/classifier"

	"github.com/google/uuid"
)

// ModelsMemoryRepository is a repository for Models
type ModelsMemoryRepository struct {
	models map[string]*Model
	mu     *sync.RWMutex
}

type Model struct {
	ID         string
	Classifier *classifier.Classifier
}

// Get returns a that matches a given ID. If no Models are found, an entity.ErrNotFound will be returned.
func (r *ModelsMemoryRepository) Get(ctx context.Context, modelID string) (*entity.Model, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	model, exists := r.models[modelID]
	if !exists {
		return nil, entity.ErrNotFound
	}

	return &entity.Model{Classifier: model.Classifier}, nil
}

// Create creates a Model
func (r *ModelsMemoryRepository) Create(ctx context.Context, e *entity.Model) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := uuid.NewString()

	r.models[id] = &Model{
		ID:         id,
		Classifier: e.Classifier,
	}

	return id, nil
}

// List returns a set of Models
func (r *ModelsMemoryRepository) List(ctx context.Context) ([]*entity.Model, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	models := make([]*entity.Model, 0, len(r.models))

	for _, model := range r.models {
		models = append(models, &entity.Model{Classifier: model.Classifier})
	}

	return models, nil
}
