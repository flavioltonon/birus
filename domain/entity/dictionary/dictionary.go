package dictionary

import (
	"github.com/agnivade/levenshtein"
)

// Dictionary is a dictionary of words
type Dictionary struct {
	words map[string]struct{}
}

// New creates a new Dictionary
func New(words ...string) *Dictionary {
	w := make(map[string]struct{})

	for _, word := range words {
		w[word] = struct{}{}
	}

	return &Dictionary{words: w}
}

// SimilarityFunc is a function that should return true when two input words are considered as similar
type SimilarityFunc func(w1, w2 string) bool

func (fn SimilarityFunc) areSimilar(w1, w2 string) bool { return fn(w1, w2) }

// Equal returns true if both input words are equal
func Equal(w1, w2 string) bool {
	return w1 == w2
}

// LevenshteinDistance returns true if the Levenshtein distance between two given words is smaller than/equal
// to a given threshold
func LevenshteinDistance(threshold int) SimilarityFunc {
	return func(w1, w2 string) bool {
		return levenshtein.ComputeDistance(w1, w2) <= threshold
	}
}

// FindWordBySimilarity looks for a word in the Dictionary, checking its similarity according to a given SimilarityFunc
func (d *Dictionary) FindWordBySimilarity(word string, fn SimilarityFunc) (string, bool) {
	if _, exists := d.words[word]; exists {
		return word, true
	}

	for w := range d.words {
		if fn.areSimilar(word, w) {
			return w, true
		}
	}

	return word, false
}
