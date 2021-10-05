package presenter

import "birus/domain/entity/shingling/classifier"

// Classifier is a entity.Classifier presenter
type Classifier struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// NewClassifier creates a new Classifier presenter
func NewClassifier(classifier *classifier.Classifier) *Classifier {
	return &Classifier{
		ID:   classifier.ID(),
		Name: classifier.Name(),
	}
}

// NewClassifierList creates a list of Classifier presenters
func NewClassifierList(classifiers []*classifier.Classifier) []*Classifier {
	result := make([]*Classifier, 0, len(classifiers))

	for _, classifier := range classifiers {
		result = append(result, NewClassifier(classifier))
	}

	return result
}
