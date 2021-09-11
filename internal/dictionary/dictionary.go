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

// similarityFunc is a function that should return true when two input words are considered as similar
type similarityFunc func(w1, w2 string) bool

func (fn similarityFunc) areSimilar(w1, w2 string) bool {
	return fn(w1, w2)
}

// Equal returns true if both input words are equal
func Equal(w1, w2 string) bool {
	return w1 == w2
}

// LevenshteinDistance returns true if the Levenshtein distance between two given words is smaller than/equal
// to a given threshold
func LevenshteinDistance(threshold int) similarityFunc {
	return func(w1, w2 string) bool {
		return levenshtein.ComputeDistance(w1, w2) <= threshold
	}
}

// FindWordBySimilarity looks for a word in the Dictionary, checking its similarity according to a given similarityFunc
func (d *Dictionary) FindWordBySimilarity(word string, fn similarityFunc) (string, bool) {
	if _, exists := d.words[word]; !exists {
		for w := range d.words {
			if fn.areSimilar(word, w) {
				return w, true
			}
		}
	}

	return word, false
}

// ReplaceWordsBySimilarity replaces words in a given set for similar words in the Dictionary, checking their similarities according
// to a given similarityFunc
func (d *Dictionary) ReplaceWordsBySimilarity(words []string, fn similarityFunc) []string {
	for i := range words {
		if replacement, found := d.FindWordBySimilarity(words[i], fn); found {
			words[i] = replacement
		}
	}

	return words
}
