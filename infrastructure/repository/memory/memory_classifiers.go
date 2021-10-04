package memory

import (
	"context"
	"sync"

	"birus/domain/entity/shingling/classifier"
)

// ClassifierRepository is a repository for Classifiers
type ClassifierRepository struct {
	classifiers map[string]*classifier.Classifier
	mu          *sync.RWMutex
}

// Create creates a Classifier
func (r *ClassifierRepository) CreateClassifier(ctx context.Context, classifier *classifier.Classifier) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.classifiers[classifier.ID()] = classifier

	return nil
}

// ListClassifiers returns a set of Classifiers
func (r *ClassifierRepository) ListClassifiers(ctx context.Context) ([]*classifier.Classifier, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	classifiers := make([]*classifier.Classifier, 0, len(r.classifiers))

	for _, classifier := range r.classifiers {
		classifiers = append(classifiers, classifier)
	}

	return classifiers, nil
}
