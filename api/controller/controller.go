package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Controller is a controller for API handlers
type Controller struct {
	usecases *Usecases
}

// New creates a new Controller
func New(usecases *Usecases) *Controller {
	return &Controller{
		usecases: usecases,
	}
}

func (c *Controller) NewRouter() http.Handler {
	router := gin.Default()

	api := router.Group("/api")
	api.POST("/tax-receipts/classifiers", c.createClassifier)
	api.GET("/tax-receipts/classifiers", c.listClassifiers)
	api.DELETE("/tax-receipts/classifiers/:classifier_id", c.deleteClassifier)
	api.POST("/tax-receipts/classify", c.classifyImage)

	return router
}
