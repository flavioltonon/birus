package controller

import (
	"birus/application/usecase"

	"github.com/gin-gonic/gin"
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
