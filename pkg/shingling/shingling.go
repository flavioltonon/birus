package shingling

import (
	"crypto/sha256"
	"strings"
)

// Shingling definition from Wikipedia: In natural language processing a w-shingling is a set of unique shingles
// (therefore n-grams) each of which is composed of contiguous subsequences of tokens within a document, which
// can then be used to ascertain the similarity between documents.
// Reference: https://en.wikipedia.org/wiki/W-shingling
type Shingling struct {
	shingles []shingle
}

// shingle is a sequence of tokens
type shingle struct {
	tokens []string
}

func (s shingle) hash() string {
	return string(sha256.New().Sum([]byte(s.ngram())))
}

func (s shingle) ngram() string {
	return strings.Join(s.tokens, " ")
}

// New creates a new Shingling for a given set of tokens and size for its n-grams
func FromTokens(tokens []string, n int) Shingling {
	if len(tokens) < n {
		panic("n should be greater than/equal to the number of tokens in the document")
	}

	// Calculating the capacity of the slice of shingles based on the number of tokens and the size
	// of the shingles we want:
	// - In a document with 4 tokens, we have 1 4-gram;
	// - In a document with 5 tokens, we have 2 4-grams;
	// - In a document with 6 tokens, we have 3 4-grams;
	//
	// - In a document with 4 tokens, we have 2 3-grams;
	// - In a document with 5 tokens, we have 3 3-grams;
	// - In a document with 6 tokens, we have 4 5-grams;
	//
	// Generalizing:
	// - In a document with t tokens, we have t-(n-1) n-grams
	shingles := make([]shingle, 0, len(tokens)-(n-1))

	for i := n; i <= len(tokens); i++ {
		shingles = append(shingles, shingle{tokens: tokens[i-n : i]})
	}

	return fromShingles(shingles)
}

// fromShingles creates a new Shingling from a set of Shingles
func fromShingles(shingles []shingle) Shingling {
	return Shingling{shingles: removeDuplicates(shingles)}
}

// removeDuplicates returns a set of unique Shingles contained in a given set of Shingles
func removeDuplicates(shingles []shingle) []shingle {
	m := make(map[string]shingle)

	for _, shingle := range shingles {
		hash := shingle.hash()

		if _, exists := m[hash]; !exists {
			m[hash] = shingle
		}
	}

	uniqueShingles := make([]shingle, 0, len(m))

	for _, shingle := range m {
		uniqueShingles = append(uniqueShingles, shingle)
	}

	return uniqueShingles
}

// unionize merges two Shinglings into a single one
func unionize(s1, s2 Shingling) Shingling {
	allShingles := make([]shingle, 0, len(s1.shingles)+len(s2.shingles))
	allShingles = append(allShingles, s1.shingles...)
	allShingles = append(allShingles, s2.shingles...)
	return fromShingles(allShingles)
}

// intersect creates a new Shingling only with Shingles that are common between 2 given Shinglings
func intersect(s1, s2 Shingling) Shingling {
	var commonShingles []shingle

	for _, ss1 := range s1.shingles {
		for _, ss2 := range s2.shingles {
			if ss1.hash() == ss2.hash() {
				commonShingles = append(commonShingles, ss1)
				break
			}
		}
	}

	return Shingling{shingles: commonShingles}
}

// JaccardSimilarity calculates the Jaccard similarity between 2 Shinglings
// https://www.cs.utah.edu/~jeffp/teaching/cs5955/L4-Jaccard+shingle.pdf
func JaccardSimilarity(s1, s2 Shingling) float64 {
	var (
		intersection = intersect(s1, s2)
		union        = unionize(s1, s2)
	)

	return float64(len(intersection.shingles)) / float64(len(union.shingles))
}
