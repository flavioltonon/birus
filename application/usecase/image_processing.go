package usecase

import (
	"birus/domain/entity/image"

	ozzo "github.com/go-ozzo/ozzo-validation/v4"
)

// ImageProcessingUsecase are usecases that define operations involving image classification
type ImageProcessingUsecase interface {
	ProcessImage(request *ProcessImageRequest) (*image.Image, error)
	ProcessImages(request *ProcessImagesRequest) ([]*image.Image, error)
}

type ProcessImageRequest struct {
	Image   *image.Image
	Options []image.ProcessOptionFunc
}

func (r ProcessImageRequest) Validate() error {
	return ozzo.ValidateStruct(&r,
		ozzo.Field(&r.Image, ozzo.Required),
		ozzo.Field(&r.Options),
	)
}

type ProcessImagesRequest struct {
	Images  []*image.Image
	Options []image.ProcessOptionFunc
}

func (r ProcessImagesRequest) Validate() error {
	return ozzo.ValidateStruct(&r,
		ozzo.Field(&r.Images, ozzo.Required),
		ozzo.Field(&r.Options),
	)
}
