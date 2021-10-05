package controller

import (
	"birus/application/usecase"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (c *Controller) newCreateClassifierRequest(ctx *gin.Context) (*usecase.CreateClassifierRequest, error) {
	request := usecase.CreateClassifierRequest{
		Name:  ctx.Request.FormValue("name"),
		Files: ctx.Request.MultipartForm.File["images"],
	}

	if err := request.Validate(); err != nil {
		return nil, err
	}

	return &request, nil
}

func (c *Controller) newListClassifiersRequest(ctx *gin.Context) (*usecase.ListClassifiersRequest, error) {
	request := usecase.ListClassifiersRequest{}

	if err := request.Validate(); err != nil {
		return nil, err
	}

	return &request, nil
}

func (c *Controller) newDeleteClassifierRequest(ctx *gin.Context) (*usecase.DeleteClassifierRequest, error) {
	request := usecase.DeleteClassifierRequest{
		ID: ctx.Param("classifier_id"),
	}

	if err := request.Validate(); err != nil {
		return nil, err
	}

	return &request, nil
}

func (c *Controller) newClassifyImageRequest(ctx *gin.Context) (*usecase.ClassifyImageRequest, error) {
	file, err := ctx.FormFile("image")
	if err != nil {
		return nil, errors.WithMessage(err, "failed to parse file from multipart form")
	}

	request := &usecase.ClassifyImageRequest{
		File: file,
	}

	if err := request.Validate(); err != nil {
		return nil, err
	}

	return request, nil
}

func (c *Controller) newReadTextFromImageRequest(ctx *gin.Context) (*usecase.ReadTextFromImageRequest, error) {
	file, err := ctx.FormFile("image")
	if err != nil {
		return nil, errors.WithMessage(err, "failed to parse file from multipart form")
	}

	request := usecase.ReadTextFromImageRequest{
		File: file,
	}

	if err := request.Validate(); err != nil {
		return nil, err
	}

	return &request, nil
}
