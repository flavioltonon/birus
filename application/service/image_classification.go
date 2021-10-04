package service

import (
	"context"

	"birus/application/usecase"
	"birus/domain/entity/shingling/classifier"

	"github.com/pkg/errors"
)

// ImageClassificationService  interface
type ImageClassificationService struct {
	textExtraction       usecase.TextExtractionUsecase
	textProcessing       usecase.TextProcessingUsecase
	classifierRepository usecase.ClassifierRepository
}

// NewImageClassificationService creates new use case
func NewImageClassificationService(
	textExtraction usecase.TextExtractionUsecase,
	textProcessing usecase.TextProcessingUsecase,
	classifierRepository usecase.ClassifierRepository,
) usecase.ImageClassificationUsecase {
	return &ImageClassificationService{
		textExtraction:       textExtraction,
		textProcessing:       textProcessing,
		classifierRepository: classifierRepository,
	}
}

// CreateClassifier creates a new typification model for a given name and a set of images
func (s *ImageClassificationService) CreateClassifier(ctx context.Context, request *usecase.CreateClassifierRequest) (*classifier.Classifier, error) {
	if err := request.Validate(); err != nil {
		return nil, errors.WithMessage(err, "failed to validate request body")
	}

	texts, err := s.textExtraction.ExtractTextFromFiles(request.Files)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to extract texts from files")
	}

	classifier := classifier.New(request.Name).Train(texts...)

	if err := s.classifierRepository.CreateClassifier(ctx, classifier); err != nil {
		return nil, errors.WithMessage(err, "failed to persist classifier")
	}

	return classifier, nil
}

func (s *ImageClassificationService) ClassifyImage(ctx context.Context, request *usecase.ClassifyImageRequest) (map[string]float64, error) {
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

	text, err := s.textExtraction.ExtractTextFromFile(request.File)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read image from file")
	}

	return set.Classify(text), nil
}
