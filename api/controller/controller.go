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

	// ImageClassification
	api.POST("/classifiers", c.createClassifier)
	api.GET("/classifiers", c.listClassifiers)
	api.DELETE("/classifiers/:classifier_id", c.deleteClassifier)
	api.POST("/classifiers/classify", c.classifyImage)

	// OpticalCharacterRecognition
	api.POST("/ocr/read", c.readTextFromImage)

	return router
}
