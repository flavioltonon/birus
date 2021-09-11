package controller

import (
	"mime/multipart"
	"net/http"

	"github.com/flavioltonon/birus/internal/logger"

	"github.com/gin-gonic/gin"
	ozzo "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// createModel creates a new image classification model with given name from a set of images
func (c *Controller) createModel(ctx *gin.Context) {
	req, err := newCreateModelRequest(ctx)
	if err != nil {
		logger.Log().Error("failed to parse request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, ctx.Error(errors.WithMessage(err, "failed to parse request body")))
		return
	}

	modelID, err := c.usecases.ImageClassification.CreateClassificationModel(ctx, req.Name, req.Images)
	if err != nil {
		logger.Log().Error("failed to create model", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, ctx.Error(errors.WithMessage(err, "failed to create model")))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"model_id": modelID})
}

// CreateModelRequest is the request body for models creation requests
type CreateModelRequest struct {
	Name   string                  `json:"name"`
	Images []*multipart.FileHeader `json:"images"`
}

func (r CreateModelRequest) validate() error {
	return ozzo.ValidateStruct(&r,
		ozzo.Field(&r.Name, ozzo.Required),
		ozzo.Field(&r.Images, ozzo.Required),
	)
}

func newCreateModelRequest(ctx *gin.Context) (*CreateModelRequest, error) {
	req := CreateModelRequest{
		Name:   ctx.Request.FormValue("name"),
		Images: ctx.Request.MultipartForm.File["images"],
	}

	if err := req.validate(); err != nil {
		return nil, err
	}

	return &req, nil
}
