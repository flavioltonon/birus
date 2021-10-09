package controller

import (
	"net/http"

	"birus/api/presenter"
	"birus/infrastructure/logger"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// listClassifiers lists the existing classifiers
func (c *Controller) listClassifiers(ctx *gin.Context) {
	request, err := c.newListClassifiersRequest(ctx)
	if err != nil {
		logger.Log().Error("failed to parse request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, ctx.Error(errors.WithMessage(err, "failed to parse request body")))
		return
	}

	classifiers, err := c.usecases.TextClassification.ListClassifiers(ctx, request)
	if err != nil {
		logger.Log().Error("failed to list classifiers", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, ctx.Error(errors.WithMessage(err, "failed to list classifiers")))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"classifiers": presenter.NewClassifierList(classifiers)})
}
