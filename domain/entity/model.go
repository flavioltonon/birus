package entity

import (
	"github.com/google/uuid"
)

// Model is a NLP model
type Model struct {
	ID    string   `json:"_id"`
	Name  string   `json:"name"`
	Texts []string `json:"texts"`
}

// NewModel creates a new Model with a given name and a given set of texts
func NewModel(name string, texts []string) (*Model, error) {
	return &Model{
		ID:    uuid.NewString(),
		Name:  name,
		Texts: texts,
	}, nil
}

// Compare returns a similarity score (0-1) between an input text to any element in the text corpus in the model
func (m *Model) Compare(text string) (float64, error) {
	return 0, nil
}
