package usecase

import (
	"birus/domain/entity/image"

	ozzo "github.com/go-ozzo/ozzo-validation/v4"
)

// OpticalCharacterRecognitionUsecase are usecases that define operations involving OCR operations
type OpticalCharacterRecognitionUsecase interface {
	ReadTextFromImage(request *ReadTextFromImageRequest) (string, error)
	ReadTextFromImages(request *ReadTextFromImagesRequest) ([]string, error)
}

type ReadTextFromImageRequest struct {
	Image   *image.Image
	Options []image.ProcessOptionFunc
}

func (r ReadTextFromImageRequest) Validate() error {
	return ozzo.ValidateStruct(&r,
		ozzo.Field(&r.Image, ozzo.Required),
		ozzo.Field(&r.Options),
	)
}

type ReadTextFromImagesRequest struct {
	Images  []*image.Image
	Options []image.ProcessOptionFunc
}

func (r ReadTextFromImagesRequest) Validate() error {
	return ozzo.ValidateStruct(&r,
		ozzo.Field(&r.Images, ozzo.Required),
		ozzo.Field(&r.Options),
	)
}
