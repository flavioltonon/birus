package presenter

import "birus/domain/entity/shingling/classifier"

// Score is a entity.Score presenter
type Score struct {
	Name       string  `json:"name"`
	Confidence float64 `json:"confidence"`
}

// NewScore creates a new Score presenter
func NewScore(score *classifier.Score) *Score {
	return &Score{
		Name:       score.Name,
		Confidence: score.Confidence,
	}
}

// NewScoreList creates a list of Score presenters
func NewScoreList(scores []*classifier.Score) []*Score {
	result := make([]*Score, 0, len(scores))

	for _, score := range scores {
		result = append(result, NewScore(score))
	}

	return result
}
