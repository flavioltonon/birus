package controller

import (
	"mime/multipart"
	"net/http"

	"github.com/flavioltonon/birus/api/presenter"
	"github.com/flavioltonon/birus/internal/logger"

	"github.com/gin-gonic/gin"
	ozzo "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// classifyImage returns the model with the highest level of similarity with a given image
func (c *Controller) classifyImage(ctx *gin.Context) {
	req, err := newClassifyImageRequest(ctx)
	if err != nil {
		logger.Log().Error("failed to parse request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, ctx.Error(errors.WithMessage(err, "failed to parse request body")))
		return
	}

	classifierID, err := c.usecases.ImageClassification.ClassifyImage(ctx, req.Image)
	if err != nil {
		logger.Log().Error("failed to classify image", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, ctx.Error(errors.WithMessage(err, "failed to classify image")))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"model": presenter.Model{Name: classifierID}})
}

// ClassifyImageRequest is the request body for models creation requests
type ClassifyImageRequest struct {
	Image *multipart.FileHeader `json:"image"`
}

func newClassifyImageRequest(ctx *gin.Context) (*ClassifyImageRequest, error) {
	image, err := ctx.FormFile("image")
	if err != nil {
		return nil, errors.WithMessage(err, "failed to parse file from multipart form")
	}

	req := ClassifyImageRequest{
		Image: image,
	}

	if err := req.validate(); err != nil {
		return nil, errors.WithMessage(err, "failed to validate request body")
	}

	return &req, nil
}

func (r ClassifyImageRequest) validate() error {
	return ozzo.ValidateStruct(&r,
		ozzo.Field(&r.Image, ozzo.Required),
	)
}
