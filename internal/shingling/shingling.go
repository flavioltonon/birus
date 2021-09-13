package shingling

import (
	"bytes"
	"encoding/gob"
)

// Shingling definition from Wikipedia: In natural language processing a w-shingling is a set of unique shingles
// (therefore n-grams) each of which is composed of contiguous subsequences of tokens within a document, which
// can then be used to ascertain the similarity between documents.
// Reference: https://en.wikipedia.org/wiki/W-shingling
type Shingling struct {
	shingles        []*Shingle
	shinglesCounter *ShinglesCounter
	multiplicity    int // defines the size of the n-grams contained in the Shingling
}

// New creates a new Shingling for a given set of tokens and size for its n-grams
func FromTokens(tokens []string, n int) *Shingling {
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
	shingles := make([]*Shingle, 0, len(tokens)-(n-1))

	for i := n; i <= len(tokens); i++ {
		shingle := NewShingle(tokens[i-n : i])
		shingles = append(shingles, shingle)
	}

	return FromShingles(shingles)
}

// FromShingles creates a new Shingling from a set of Shingles
func FromShingles(shingles []*Shingle) *Shingling {
	if len(shingles) == 0 {
		panic("at least one shingle should be given as input")
	}

	shingling := Shingling{
		shingles:        make([]*Shingle, 0, len(shingles)),
		shinglesCounter: NewShinglesCounter(),
		multiplicity:    len(shingles[0].tokens),
	}

	for i := range shingles {
		shingling.addShingle(shingles[i])
	}

	return &shingling
}

// GetShingles returns the unique Shingles that compose the Shingling
func (s *Shingling) GetShingles() []*Shingle {
	return s.shingles
}

// GetMultiplicity returns the multiplicity of the Shingling, which defines the size of all Shingles it should contain
func (s *Shingling) GetMultiplicity() int {
	return s.multiplicity
}

func (s *Shingling) addShingle(shingle *Shingle) {
	if shingle.GetMultiplicity() != s.multiplicity {
		panic("all shingles should contain the same multiplicity in a shingling")
	}

	h := shingle.GetHash()

	if _, exists := s.shinglesCounter.GetValue(h); !exists {
		s.shingles = append(s.shingles, shingle)
	}

	s.shinglesCounter.Increment(h)
}

// unionize merges two Shinglings into a single one
func unionize(s1, s2 *Shingling) *Shingling {
	allShingles := make([]*Shingle, 0, len(s1.shingles)+len(s2.shingles))
	allShingles = append(allShingles, s1.shingles...)
	allShingles = append(allShingles, s2.shingles...)
	return FromShingles(allShingles)
}

// intersect creates a new Shingling only with Shingles that are common between 2 given Shinglings
func intersect(s1, s2 *Shingling) *Shingling {
	var commonShingles []*Shingle

	for i := range s1.shingles {
		var (
			ss1 = s1.shingles[i]
			h1  = ss1.GetHash()
		)

		for j := range s2.shingles {
			var (
				ss2 = s2.shingles[j]
				h2  = ss2.GetHash()
			)

			if h1 == h2 {
				commonShingles = append(commonShingles, ss1)
				break
			}
		}
	}

	return &Shingling{shingles: commonShingles}
}

// JaccardSimilarity calculates the Jaccard similarity between 2 Shinglings
// https://www.cs.utah.edu/~jeffp/teaching/cs5955/L4-Jaccard+shingle.pdf
func JaccardSimilarity(s1, s2 *Shingling) float64 {
	var (
		intersection = intersect(s1, s2)
		union        = unionize(s1, s2)
	)

	return float64(len(intersection.shingles)) / float64(len(union.shingles))
}

type GobShingling struct {
	Shingles        []*Shingle
	ShinglesCounter *ShinglesCounter
	Multiplicity    int
}

func (s *Shingling) GobEncode() ([]byte, error) {
	var buffer bytes.Buffer

	if err := gob.NewEncoder(&buffer).Encode(&GobShingling{
		Shingles:        s.shingles,
		ShinglesCounter: s.shinglesCounter,
		Multiplicity:    s.multiplicity,
	}); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (s *Shingling) GobDecode(b []byte) error {
	var (
		buffer = bytes.NewBuffer(b)
		reader GobShingling
	)

	if err := gob.NewDecoder(buffer).Decode(&reader); err != nil {
		return err
	}

	s.shingles = reader.Shingles
	s.shinglesCounter = reader.ShinglesCounter
	s.multiplicity = reader.Multiplicity
	return nil
}
