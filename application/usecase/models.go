package usecase

import (
	"context"

	"github.com/flavioltonon/birus/domain/entity"
	"github.com/flavioltonon/birus/internal/shingling"
)

// ModelsReader is a reader interface for Models
type ModelsReader interface {
	Get(ctx context.Context, modelID string) (*entity.Model, error)
	List(ctx context.Context) ([]*entity.Model, error)
}

// ModelsWriter is a writer interface for Models
type ModelsWriter interface {
	Create(ctx context.Context, e *entity.Model) (string, error)
}

// ModelsRepository interface
type ModelsRepository interface {
	ModelsReader
	ModelsWriter
}

// ModelsUsecase interface
type ModelsUsecase interface {
	GetModel(ctx context.Context, modelID string) (*entity.Model, error)
	ListModels(ctx context.Context) ([]*entity.Model, error)
	CreateModel(ctx context.Context, name string, shinglings []*shingling.Shingling) (modelID string, err error)
}
