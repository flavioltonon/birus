package usecase

import (
	"context"
	"mime/multipart"
)

// ImageClassificationUsecase interface
type ImageClassificationUsecase interface {
	CreateClassificationModel(ctx context.Context, name string, files []*multipart.FileHeader) (modelID string, err error)
	ClassifyImage(ctx context.Context, file *multipart.FileHeader) (classifierID string, err error)
}
