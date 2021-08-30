package usecase

import (
	"context"
	"mime/multipart"

	"github.com/flavioltonon/birus/domain/entity"
)

// ImageClassificationUsecase interface
type ImageClassificationUsecase interface {
	CreateModel(ctx context.Context, name string, files []*multipart.FileHeader) (string, error)
	ClassifyImage(ctx context.Context, file *multipart.FileHeader) (*entity.Model, error)
}
