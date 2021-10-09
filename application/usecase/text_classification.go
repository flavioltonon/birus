package usecase

import (
	"context"

	"birus/domain/entity/shingling/classifier"

	ozzo "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// TextClassificationUsecase are usecases that define operations involving text classification
type TextClassificationUsecase interface {
	CreateClassifier(ctx context.Context, request *CreateClassifierRequest) (*classifier.Classifier, error)
	ListClassifiers(ctx context.Context, request *ListClassifiersRequest) ([]*classifier.Classifier, error)
	DeleteClassifier(ctx context.Context, request *DeleteClassifierRequest) error
	ClassifyText(ctx context.Context, request *ClassifyTextRequest) ([]*classifier.Score, error)
}

type CreateClassifierRequest struct {
	Name  string
	Texts []string
}

func (r CreateClassifierRequest) Validate() error {
	return ozzo.ValidateStruct(&r,
		ozzo.Field(&r.Name, ozzo.Required),
		ozzo.Field(&r.Texts, ozzo.Required),
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

type ClassifyTextRequest struct {
	Text string
}

func (r ClassifyTextRequest) Validate() error {
	return ozzo.ValidateStruct(&r,
		ozzo.Field(&r.Text, ozzo.Required),
	)
}

type ClassifierRepository interface {
	GetClassifier(ctx context.Context, classifierID string) (*classifier.Classifier, error)
	CreateClassifier(ctx context.Context, classifier *classifier.Classifier) error
	ListClassifiers(ctx context.Context) ([]*classifier.Classifier, error)
	DeleteClassifier(ctx context.Context, classifierID string) error
}
