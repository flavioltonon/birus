package controller

import (
	"birus/api/middleware"
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

	// TextClassification
	textClassification := api.Group("/text-classification")
	textClassification.POST("/classifiers", c.createClassifier)
	textClassification.GET("/classifiers", c.listClassifiers)
	textClassification.DELETE("/classifiers/:classifier_id", c.deleteClassifier)
	textClassification.POST("/classify", c.classifyText)

	// OpticalCharacterRecognition
	ocr := api.Group("/ocr")
	ocr.POST("/read", c.readTextFromImage)
	ocr.POST("/read/batch", c.readTextFromImages)

	router.Use(middleware.SetRequestID)

	return router
}
