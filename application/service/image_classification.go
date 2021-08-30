package service

import (
	"context"
	"mime/multipart"

	"github.com/flavioltonon/birus/application/usecase"
	"github.com/flavioltonon/birus/domain/entity"
	"github.com/flavioltonon/birus/infrastructure/engine"
	"github.com/flavioltonon/birus/internal/image"

	ozzo "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
)

// Ensure that ImageClassificationService implements usecase.TypificationUsecase
var _ usecase.ImageClassificationUsecase = (*ImageClassificationService)(nil)

// ImageClassificationService  interface
type ImageClassificationService struct {
	models usecase.ModelsUsecase
	engine engine.Engine
}

// NewImageClassificationService creates new use case
func NewImageClassificationService(m usecase.ModelsUsecase, e engine.Engine) *ImageClassificationService {
	return &ImageClassificationService{
		models: m,
		engine: e,
	}
}

// CreateModel creates a new typification model for a given name and a set of images
func (s *ImageClassificationService) CreateModel(ctx context.Context, name string, files []*multipart.FileHeader) (string, error) {
	if err := ozzo.Required.Validate(files); err != nil {
		return "", errors.WithMessage(err, "failed to validate files")
	}

	images, err := image.FromBulkMultipartFileHeaders(files)
	if err != nil {
		return "", errors.WithMessage(err, "failed to read images from multipart file headers")
	}

	texts := make([]string, 0, len(images))

	for _, image := range images {
		text, err := s.engine.ExtractTextFromImage(image)
		if err != nil {
			return "", errors.WithMessage(err, "failed to extract text from image")
		}

		texts = append(texts, text)
	}

	return s.models.CreateModel(ctx, name, texts)
}

func (s *ImageClassificationService) ClassifyImage(ctx context.Context, file *multipart.FileHeader) (*entity.Model, error) {
	image, err := image.FromMultipartFileHeader(file)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read image from multipart file header")
	}

	text, err := s.engine.ExtractTextFromImage(image)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to extract text from image")
	}

	models, err := s.models.ListModels(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to list models")
	}

	var (
		highestSimilarityModel *entity.Model
		highestSimilarity      = 0.0
	)

	for _, model := range models {
		similarity, err := model.Compare(text)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to compare text to model")
		}

		if similarity > highestSimilarity {
			highestSimilarityModel = &model
		}
	}

	if highestSimilarityModel == nil {
		return nil, ErrNoTypificationMatches
	}

	return highestSimilarityModel, nil
}
