package usecase

import (
	"mime/multipart"

	ozzo "github.com/go-ozzo/ozzo-validation/v4"
)

// OpticalCharacterRecognitionUsecase are usecases that define operations involving OCR operations
type OpticalCharacterRecognitionUsecase interface {
	ReadTextFromImage(request *ReadTextFromImageRequest) (string, error)
	ReadTextFromImages(request *ReadTextFromImagesRequest) ([]string, error)
}

type ReadTextFromImageRequest struct {
	File *multipart.FileHeader
}

func (r ReadTextFromImageRequest) Validate() error {
	return ozzo.ValidateStruct(&r,
		ozzo.Field(&r.File, ozzo.Required),
	)
}

type ReadTextFromImagesRequest struct {
	Files []*multipart.FileHeader
}

func (r ReadTextFromImagesRequest) Validate() error {
	return ozzo.ValidateStruct(&r,
		ozzo.Field(&r.Files, ozzo.Required),
	)
}
