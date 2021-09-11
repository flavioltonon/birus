package entity

import (
	"github.com/flavioltonon/birus/internal/shingling/classifier"

	ozzo "github.com/go-ozzo/ozzo-validation"
)

// Model is a machine learning model with a classifier based on the w-shingling NLP technique
type Model struct {
	Classifier *classifier.Classifier `json:"classifier"`
}

// NewModel creates a new Model from a given classifier
func NewModel(classifier *classifier.Classifier) (*Model, error) {
	if err := ozzo.Required.Validate(classifier); err != nil {
		return nil, err
	}

	return &Model{Classifier: classifier}, nil
}
