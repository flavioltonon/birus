package controller

import (
	"net/http"

	"birus/domain/entity"
	"birus/infrastructure/logger"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// deleteClassifier deletes a classifier
func (c *Controller) deleteClassifier(ctx *gin.Context) {
	request, err := c.newDeleteClassifierRequest(ctx)
	if err != nil {
		logger.Log().Error("failed to parse request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, ctx.Error(errors.WithMessage(err, "failed to parse request body")))
		return
	}

	if err := c.usecases.TextClassification.DeleteClassifier(ctx, request); err != nil {
		logger.Log().Error("failed to delete classifier", zap.Error(err))

		status := http.StatusNotFound

		if !errors.Is(err, entity.ErrNotFound) {
			status = http.StatusInternalServerError
		}

		ctx.JSON(status, ctx.Error(errors.WithMessage(err, "failed to delete classifier")))
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
