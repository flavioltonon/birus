package controller

import (
	"net/http"

	"birus/api/presenter"
	"birus/infrastructure/logger"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// processImage processes an image using a given set of options
func (c *Controller) processImage(ctx *gin.Context) {
	request, err := c.newProcessImageRequest(ctx)
	if err != nil {
		logger.Log().Error("failed to parse request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, ctx.Error(errors.WithMessage(err, "failed to parse request body")))
		return
	}

	image, err := c.usecases.ImageProcessing.ProcessImage(request)
	if err != nil {
		logger.Log().Error("failed to process image", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, ctx.Error(errors.WithMessage(err, "failed to process image")))
		return
	}

	if gin.Mode() == gin.DebugMode {
		// Save image for debugging
		image.Save("/output/image.jpg")
	}

	ctx.JSON(http.StatusOK, gin.H{"image": presenter.NewImage(image)})
}
