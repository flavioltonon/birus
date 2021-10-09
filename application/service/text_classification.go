package service

import (
	"context"

	"birus/application/usecase"
	"birus/domain/entity"
	"birus/domain/entity/shingling/classifier"

	"github.com/pkg/errors"
)

// TextClassificationService  interface
type TextClassificationService struct {
	opticalCharacterRecognition usecase.OpticalCharacterRecognitionUsecase
	textProcessing              usecase.TextProcessingUsecase
	classifierRepository        usecase.ClassifierRepository
}

// NewTextClassificationService creates new use case
func NewTextClassificationService(classifierRepository usecase.ClassifierRepository) usecase.TextClassificationUsecase {
	return &TextClassificationService{classifierRepository: classifierRepository}
}

// CreateClassifier creates a new typification model for a given name and a set of texts
func (s *TextClassificationService) CreateClassifier(ctx context.Context, request *usecase.CreateClassifierRequest) (*classifier.Classifier, error) {
	if err := request.Validate(); err != nil {
		return nil, errors.WithMessage(err, "failed to validate request body")
	}

	classifier := classifier.New(request.Name).Train(request.Texts...)

	if err := s.classifierRepository.CreateClassifier(ctx, classifier); err != nil {
		return nil, errors.WithMessage(err, "failed to persist classifier")
	}

	return classifier, nil
}

// ListClassifiers lists the existing classifiers
func (s *TextClassificationService) ListClassifiers(ctx context.Context, request *usecase.ListClassifiersRequest) ([]*classifier.Classifier, error) {
	if err := request.Validate(); err != nil {
		return nil, errors.WithMessage(err, "failed to validate request body")
	}

	return s.classifierRepository.ListClassifiers(ctx)
}

// DeleteClassifier deletes an existing classifier
func (s *TextClassificationService) DeleteClassifier(ctx context.Context, request *usecase.DeleteClassifierRequest) error {
	if err := request.Validate(); err != nil {
		return errors.WithMessage(err, "failed to validate request body")
	}

	if _, err := s.classifierRepository.GetClassifier(ctx, request.ID); errors.Is(err, entity.ErrNotFound) {
		return errors.WithMessage(err, "failed to get classifier")
	}

	return s.classifierRepository.DeleteClassifier(ctx, request.ID)
}

func (s *TextClassificationService) ClassifyText(ctx context.Context, request *usecase.ClassifyTextRequest) ([]*classifier.Score, error) {
	if err := request.Validate(); err != nil {
		return nil, errors.WithMessage(err, "failed to validate request body")
	}

	classifiers, err := s.classifierRepository.ListClassifiers(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to list classifiers")
	}

	set := classifier.NewSet()

	for _, classifier := range classifiers {
		set.AddClassifier(classifier)
	}

	return set.Classify(request.Text)
}
