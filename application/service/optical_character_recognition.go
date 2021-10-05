package service

import (
	"birus/application/usecase"
	"birus/domain/entity/image"
	"birus/infrastructure/engine"

	"github.com/pkg/errors"
)

// OpticalCharacterRecognitionService is a text extraction service
type OpticalCharacterRecognitionService struct {
	e engine.Engine
}

// NewOpticalCharacterRecognitionService creates a new OpticalCharacterRecognitionService
func NewOpticalCharacterRecognitionService(e engine.Engine) usecase.OpticalCharacterRecognitionUsecase {
	return &OpticalCharacterRecognitionService{e: e}
}

// ReadTextFromImage uses an OCR engine to extract text from a given multipart.FileHeader
func (s *OpticalCharacterRecognitionService) ReadTextFromImage(request *usecase.ReadTextFromImageRequest) (string, error) {
	image, err := image.FromMultipartFileHeader(request.File)
	if err != nil {
		return "", errors.WithMessage(err, "failed to read image from file")
	}

	text, err := s.e.ExtractTextFromImage(image.Bytes())
	if err != nil {
		return "", errors.WithMessage(err, "failed to extract text from image")
	}

	return text, nil
}

// ReadTextFromImages uses an OCR engine to extract texts from a given set of multipart.FileHeaders
func (s *OpticalCharacterRecognitionService) ReadTextFromImages(request *usecase.ReadTextFromImagesRequest) ([]string, error) {
	texts := make([]string, 0, len(request.Files))

	for _, file := range request.Files {
		text, err := s.ReadTextFromImage(&usecase.ReadTextFromImageRequest{
			File: file,
		})
		if err != nil {
			return nil, errors.WithMessage(err, "failed to extract text from file")
		}

		texts = append(texts, text)
	}

	return texts, nil
}
