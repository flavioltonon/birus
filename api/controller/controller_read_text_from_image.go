package controller

import (
	"net/http"

	"birus/infrastructure/logger"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// readTextFromImage returns the text contained in an image, extracted via OCR
func (c *Controller) readTextFromImage(ctx *gin.Context) {
	request, err := c.newReadTextFromImageRequest(ctx)
	if err != nil {
		logger.Log().Error("failed to parse request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, ctx.Error(errors.WithMessage(err, "failed to parse request body")))
		return
	}

	text, err := c.usecases.OpticalCharacterRecognition.ReadTextFromImage(request)
	if err != nil {
		logger.Log().Error("failed to read text from image", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, ctx.Error(errors.WithMessage(err, "failed to read text from image")))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"text": text})
}
