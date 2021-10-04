package controller

import (
	"net/http"

	"birus/api/presenter"
	"birus/application/usecase"
	"birus/infrastructure/logger"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// createClassifier creates a new image classification model with given name from a set of images
func (c *Controller) createClassifier(ctx *gin.Context) {
	request, err := c.newCreateClassifierRequest(ctx)
	if err != nil {
		logger.Log().Error("failed to parse request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, ctx.Error(errors.WithMessage(err, "failed to parse request body")))
		return
	}

	classifier, err := c.usecases.ImageClassification.CreateClassifier(ctx, request)
	if err != nil {
		logger.Log().Error("failed to create classifier", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, ctx.Error(errors.WithMessage(err, "failed to create classifier")))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"classifier": presenter.NewClassifier(classifier)})
}

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
