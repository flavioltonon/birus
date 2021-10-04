package usecase

import (
	"mime/multipart"
)

// TextExtractionUsecase are usecases that define operations involving text extraction using an OCR engine
type TextExtractionUsecase interface {
	ExtractTextFromFile(file *multipart.FileHeader) (string, error)
	ExtractTextFromFiles(files []*multipart.FileHeader) ([]string, error)
}
