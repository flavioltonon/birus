package classifier

import (
	"math"
)

// Set is a set of Classifiers
type Set struct {
	classifiers []*Classifier
}

// NewSet creates a new Set
func NewSet() *Set {
	return new(Set)
}

// AddClassifier adds a classifier to the set
func (s *Set) AddClassifier(classifier *Classifier) {
	s.classifiers = append(s.classifiers, classifier)
}

// Classify returns the similarity score with all the Set classifiers for a given Shingling
func (s *Set) Classify(text string) map[string]float64 {
	var (
		scores     = make(map[string]float64, len(s.classifiers))
		totalScore float64
	)

	for i := range s.classifiers {
		classifier := s.classifiers[i]
		score := classifier.Classify(text)

		// use the square of the scores to accentuate their differences
		scoreSquare := math.Pow(score, 2)
		scores[classifier.Name()] = scoreSquare
		totalScore += scoreSquare
	}

	// normalize the scores so their sum will be equal to 1 (100%)
	for id, score := range scores {
		scores[id] = score / totalScore
	}

	return scores
}
