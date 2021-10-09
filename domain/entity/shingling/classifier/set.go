package classifier

import (
	"math"
	"sort"

	"github.com/pkg/errors"
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

type Score struct {
	Name       string
	Confidence float64
}

// NewScore creates a new Score
func NewScore(name string, confidence float64) *Score {
	return &Score{
		Name:       name,
		Confidence: confidence,
	}
}

// Classify returns the similarity scores with all the Set classifiers for a given Shingling, sorted by confidence (descending)
// The first score in the results will always represent the best match.
func (s *Set) Classify(text string) ([]*Score, error) {
	if len(s.classifiers) == 0 {
		return nil, errors.New("no classifiers available in current set")
	}

	var (
		scoresByName = make(map[string]float64, len(s.classifiers))
		totalScore   float64
	)

	for i := range s.classifiers {
		classifier := s.classifiers[i]
		score := classifier.Classify(text)

		// use the square of the scores to accentuate their differences
		scoreSquare := math.Pow(score, 2)
		scoresByName[classifier.Name()] = scoreSquare
		totalScore += scoreSquare
	}

	scores := make([]*Score, 0, len(scoresByName))

	// normalize the scores so their sum will be equal to 1 (100%) and add them to a slice to allow sorting
	for name, score := range scoresByName {
		scores = append(scores, &Score{
			Name:       name,
			Confidence: score / totalScore,
		})
	}

	sort.Sort(byConfidenceDesc(scores))

	return scores, nil
}

type byConfidenceDesc []*Score

func (a byConfidenceDesc) Len() int           { return len(a) }
func (a byConfidenceDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byConfidenceDesc) Less(i, j int) bool { return a[i].Confidence > a[j].Confidence }
