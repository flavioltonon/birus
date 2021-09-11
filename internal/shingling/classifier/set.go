package classifier

import "github.com/flavioltonon/birus/internal/shingling"

type Set struct {
	classifiers []*Classifier
}

func NewSet() *Set {
	return new(Set)
}

func (s *Set) AddClassifier(classifier *Classifier) {
	s.classifiers = append(s.classifiers, classifier)
}

func (s *Set) Classify(shingling *shingling.Shingling) map[string]float64 {
	var (
		scores     = make(map[string]float64, len(s.classifiers))
		totalScore float64
	)

	for i := range s.classifiers {
		classifier := s.classifiers[i]
		score := classifier.Classify(shingling)
		scores[classifier.ID] = score
		totalScore += score
	}

	for id, score := range scores {
		scores[id] = score / totalScore
	}

	return scores
}
