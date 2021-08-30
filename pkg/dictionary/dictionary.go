package dictionary

import "github.com/agnivade/levenshtein"

// Dictionary is a dictionary of words
type Dictionary struct {
	words []string
}

// New creates a new Dictionary
func New(words ...string) *Dictionary {
	return &Dictionary{words: words}
}

// similarityFunc is a function that should return true when two input words are considered as similar
type similarityFunc func(w1, w2 string) bool

// Equal returns true if both input words are equal
func Equal(w1, w2 string) bool {
	return w1 == w2
}

// LevenshteinDistanceThreshold returns a similarityFunc that returns true when the distance between two words is
// less than/equal to a defined threshold
func LevenshteinDistanceThreshold(threshold int) similarityFunc {
	return func(w1, w2 string) bool {
		return levenshtein.ComputeDistance(w1, w2) <= threshold
	}
}

// FindWord looks for an exact match of a word in the Dictionary
func (d *Dictionary) FindWord(word string) (string, bool) {
	return d.FindWordBySimilarity(word, Equal)
}

// FindWordBySimilarity looks for a word in the Dictionary, checking its similarity according to a given similarityFunc
func (d *Dictionary) FindWordBySimilarity(word string, fn similarityFunc) (string, bool) {
	for _, w := range d.words {
		if fn(word, w) {
			return w, true
		}
	}

	return word, false
}
