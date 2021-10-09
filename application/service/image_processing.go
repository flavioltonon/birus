package service

import (
	"birus/application/usecase"
	"birus/domain/entity/image"

	"github.com/pkg/errors"
)

// ImageProcessingService is a service for image processing
type ImageProcessingService struct{}

// NewImageProcessingService creates new use case
func NewImageProcessingService() usecase.ImageProcessingUsecase {
	return new(ImageProcessingService)
}

// ProcessImage processes an image with a given set of options
func (h *ImageProcessingService) ProcessImage(request *usecase.ProcessImageRequest) (*image.Image, error) {
	if err := request.Validate(); err != nil {
		return nil, errors.WithMessage(err, "failed to validate request body")
	}

	return request.Image.Process(request.Options...)
}

// ProcessImages processes a list of images with a given set of options
func (h *ImageProcessingService) ProcessImages(request *usecase.ProcessImagesRequest) ([]*image.Image, error) {
	if err := request.Validate(); err != nil {
		return nil, errors.WithMessage(err, "failed to validate request body")
	}

	images := make([]*image.Image, 0, len(request.Images))

	for _, image := range request.Images {
		image, err := h.ProcessImage(&usecase.ProcessImageRequest{
			Image:   image,
			Options: request.Options,
		})
		if err != nil {
			return nil, errors.WithMessage(err, "failed to process image")
		}

		images = append(images, image)
	}

	return images, nil
}
