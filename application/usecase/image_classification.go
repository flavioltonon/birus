package usecase

import (
	"context"
	"mime/multipart"

	"birus/domain/entity/shingling/classifier"

	ozzo "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// ImageClassificationUsecase are usecases that define operations involving image classification
type ImageClassificationUsecase interface {
	CreateClassifier(ctx context.Context, request *CreateClassifierRequest) (*classifier.Classifier, error)
	ListClassifiers(ctx context.Context, request *ListClassifiersRequest) ([]*classifier.Classifier, error)
	DeleteClassifier(ctx context.Context, request *DeleteClassifierRequest) error
	ClassifyImage(ctx context.Context, request *ClassifyImageRequest) (map[string]float64, error)
}

type CreateClassifierRequest struct {
	Name  string
	Files []*multipart.FileHeader
}

func (r CreateClassifierRequest) Validate() error {
	return ozzo.ValidateStruct(&r,
		ozzo.Field(&r.Name, ozzo.Required),
		ozzo.Field(&r.Files, ozzo.Required),
	)
}

type ListClassifiersRequest struct{}

func (r ListClassifiersRequest) Validate() error {
	return ozzo.ValidateStruct(&r)
}

type DeleteClassifierRequest struct {
	ID string
}

func (r DeleteClassifierRequest) Validate() error {
	return ozzo.ValidateStruct(&r,
		ozzo.Field(&r.ID, ozzo.Required, is.UUIDv4),
	)
}

type ClassifyImageRequest struct {
	File *multipart.FileHeader
}

func (r ClassifyImageRequest) Validate() error {
	return ozzo.ValidateStruct(&r,
		ozzo.Field(&r.File, ozzo.Required),
	)
}

type ClassifierRepository interface {
	GetClassifier(ctx context.Context, classifierID string) (*classifier.Classifier, error)
	CreateClassifier(ctx context.Context, classifier *classifier.Classifier) error
	ListClassifiers(ctx context.Context) ([]*classifier.Classifier, error)
	DeleteClassifier(ctx context.Context, classifierID string) error
}
