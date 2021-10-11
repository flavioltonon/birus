package controller

import (
	"birus/infrastructure/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// readTextFromImages returns the text contained in an image, extracted via OCR
func (c *Controller) readTextFromImages(ctx *gin.Context) {
	request, err := c.newReadTextFromImagesRequest(ctx)
	if err != nil {
		logger.Log().Error("failed to parse request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, ctx.Error(errors.WithMessage(err, "failed to parse request body")))
		return
	}

	texts, err := c.usecases.OpticalCharacterRecognition.ReadTextFromImages(request)
	if err != nil {
		logger.Log().Error("failed to read text from images", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, ctx.Error(errors.WithMessage(err, "failed to read text from images")))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"texts": texts})
}
