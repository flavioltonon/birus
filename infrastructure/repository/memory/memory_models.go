package memory

import (
	"context"
	"sync"

	"github.com/flavioltonon/birus/domain/entity"
	"github.com/flavioltonon/birus/internal/shingling/classifier"

	"github.com/google/uuid"
)

// ModelsRepository is a repository for Models
type ModelsRepository struct {
	models map[string]*Model
	mu     *sync.RWMutex
}

type Model struct {
	ID         string
	Classifier *classifier.Classifier
}

// Get returns a that matches a given ID. If no Models are found, an entity.ErrNotFound will be returned.
func (r *ModelsRepository) Get(ctx context.Context, modelID string) (*entity.Model, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	model, exists := r.models[modelID]
	if !exists {
		return nil, entity.ErrNotFound
	}

	return entity.NewModel(model.Classifier)
}

// Create creates a Model
func (r *ModelsRepository) Create(ctx context.Context, e *entity.Model) (string, error) {
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
func (r *ModelsRepository) List(ctx context.Context) ([]*entity.Model, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	models := make([]*entity.Model, 0, len(r.models))

	for _, model := range r.models {
		e, err := entity.NewModel(model.Classifier)
		if err != nil {
			return nil, err
		}

		models = append(models, e)
	}

	return models, nil
}
