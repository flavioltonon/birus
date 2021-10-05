package controller

import (
	"net/http"

	"birus/infrastructure/logger"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// classifyImage returns the model with the highest level of similarity with a given image
func (c *Controller) classifyImage(ctx *gin.Context) {
	request, err := c.newClassifyImageRequest(ctx)
	if err != nil {
		logger.Log().Error("failed to parse request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, ctx.Error(errors.WithMessage(err, "failed to parse request body")))
		return
	}

	scores, err := c.usecases.ImageClassification.ClassifyImage(ctx, request)
	if err != nil {
		logger.Log().Error("failed to classify image", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, ctx.Error(errors.WithMessage(err, "failed to classify image")))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"scores": scores})
}
