package service

import (
	"birus/application/usecase"
	"birus/domain/entity/image"
	"birus/infrastructure/engine"
	"mime/multipart"

	"github.com/pkg/errors"
)

// TextExtrationService is a text extraction service
type TextExtrationService struct {
	e engine.Engine
}

// NewTextExtrationService creates a new TextExtrationService
func NewTextExtrationService(e engine.Engine) usecase.TextExtractionUsecase {
	return &TextExtrationService{e: e}
}

// ExtractTextFromFile uses an OCR engine to extract text from a given multipart.FileHeader
func (s *TextExtrationService) ExtractTextFromFile(file *multipart.FileHeader) (string, error) {
	image, err := image.FromMultipartFileHeader(file)
	if err != nil {
		return "", errors.WithMessage(err, "failed to read image from file")
	}

	text, err := s.e.ExtractTextFromImage(image.Bytes())
	if err != nil {
		return "", errors.WithMessage(err, "failed to extract text from image")
	}

	return text, nil
}

// ExtractTextFromFiles uses an OCR engine to extract texts from a given set of multipart.FileHeaders
func (s *TextExtrationService) ExtractTextFromFiles(files []*multipart.FileHeader) ([]string, error) {
	texts := make([]string, 0, len(files))

	for _, file := range files {
		text, err := s.ExtractTextFromFile(file)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to extract text from file")
		}

		texts = append(texts, text)
	}

	return texts, nil
}
