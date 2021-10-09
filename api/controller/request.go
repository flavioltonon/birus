package controller

import (
	"birus/application/usecase"
	"birus/domain/entity/image"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (c *Controller) newCreateClassifierRequest(ctx *gin.Context) (*usecase.CreateClassifierRequest, error) {
	var request usecase.CreateClassifierRequest

	if err := ctx.BindJSON(&request); err != nil {
		return nil, errors.WithMessage(err, "failed to decode request body")
	}

	if err := request.Validate(); err != nil {
		return nil, errors.WithMessage(err, "failed to validate request body")
	}

	return &request, nil
}

func (c *Controller) newListClassifiersRequest(ctx *gin.Context) (*usecase.ListClassifiersRequest, error) {
	var request usecase.ListClassifiersRequest

	if err := request.Validate(); err != nil {
		return nil, errors.WithMessage(err, "failed to validate request body")
	}

	return &request, nil
}

func (c *Controller) newDeleteClassifierRequest(ctx *gin.Context) (*usecase.DeleteClassifierRequest, error) {
	request := usecase.DeleteClassifierRequest{
		ID: ctx.Param("classifier_id"),
	}

	if err := request.Validate(); err != nil {
		return nil, errors.WithMessage(err, "failed to validate request body")
	}

	return &request, nil
}

func (c *Controller) newClassifyTextRequest(ctx *gin.Context) (*usecase.ClassifyTextRequest, error) {
	var request usecase.ClassifyTextRequest

	if err := ctx.BindJSON(&request); err != nil {
		return nil, errors.WithMessage(err, "failed to decode request body")
	}

	if err := request.Validate(); err != nil {
		return nil, errors.WithMessage(err, "failed to validate request body")
	}

	return &request, nil
}

func (c *Controller) newReadTextFromImageRequest(ctx *gin.Context) (*usecase.ReadTextFromImageRequest, error) {
	var request usecase.ReadTextFromImageRequest

	file, err := ctx.FormFile("file")
	if err != nil {
		return nil, errors.WithMessage(err, "failed to parse file from multipart form")
	}

	request.Image, err = image.FromMultipartFileHeader(file)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read image from file")
	}

	if optionsStr := ctx.Request.FormValue("options"); optionsStr != "" {
		request.Options, err = image.ParseProcessOptions(optionsStr)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to parse process options")
		}
	}

	if err := request.Validate(); err != nil {
		return nil, errors.WithMessage(err, "failed to validate request body")
	}

	return &request, nil
}

func (c *Controller) newReadTextFromImagesRequest(ctx *gin.Context) (*usecase.ReadTextFromImagesRequest, error) {
	var request usecase.ReadTextFromImagesRequest

	form, err := ctx.MultipartForm()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to parse file from multipart form")
	}

	request.Images, err = image.FromMultipartFileHeaders(form.File["files"])
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read image from file")
	}

	if optionsStr := ctx.Request.FormValue("options"); optionsStr != "" {
		request.Options, err = image.ParseProcessOptions(optionsStr)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to parse process options")
		}
	}

	if err := request.Validate(); err != nil {
		return nil, errors.WithMessage(err, "failed to validate request body")
	}

	return &request, nil
}
