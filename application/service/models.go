package service

import (
	"context"

	"github.com/flavioltonon/birus/application/usecase"
	"github.com/flavioltonon/birus/domain/entity"
	"github.com/flavioltonon/birus/infrastructure/engine"
	"github.com/flavioltonon/birus/internal/shingling"
	"github.com/flavioltonon/birus/internal/shingling/classifier"
)

// Ensure that ModelsService implements usecase.ModelsUsecase
var _ usecase.ModelsUsecase = (*ModelsService)(nil)

// ModelsService  interface
type ModelsService struct {
	repository usecase.ModelsRepository
}

// NewModelsService creates new use case
func NewModelsService(r usecase.ModelsRepository, e engine.Engine) *ModelsService {
	return &ModelsService{
		repository: r,
	}
}

// GetModel fetches a Model with a given modelID from the repository
func (s *ModelsService) GetModel(ctx context.Context, modelID string) (*entity.Model, error) {
	return s.repository.Get(ctx, modelID)
}

// ListModels fetches a set of Models from the repository
func (s *ModelsService) ListModels(ctx context.Context) ([]*entity.Model, error) {
	return s.repository.List(ctx)
}

// CreateModel creates and persists a new model into the repository
func (s *ModelsService) CreateModel(ctx context.Context, name string, shinglings []*shingling.Shingling) (string, error) {
	classifier := classifier.New(name).Train(shinglings...)

	model, err := entity.NewModel(classifier)
	if err != nil {
		return "", err
	}

	return s.repository.Create(ctx, model)
}
