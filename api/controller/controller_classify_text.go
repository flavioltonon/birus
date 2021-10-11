package controller

import (
	"net/http"

	"birus/api/presenter"
	"birus/infrastructure/logger"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// classifyText returns the model with the highest level of similarity with a given text
func (c *Controller) classifyText(ctx *gin.Context) {
	request, err := c.newClassifyTextRequest(ctx)
	if err != nil {
		logger.Log().Error("failed to parse request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, ctx.Error(errors.WithMessage(err, "failed to parse request body")))
		return
	}

	scores, err := c.usecases.TextClassification.ClassifyText(ctx, request)
	if err != nil {
		logger.Log().Error("failed to classify text", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, ctx.Error(errors.WithMessage(err, "failed to classify text")))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"result": presenter.NewScore(scores[0])})
}
