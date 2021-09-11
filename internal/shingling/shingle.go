package shingling

import (
	"bytes"
	"encoding/gob"
	"strings"
)

// Shingle is a sequence of tokens
type Shingle struct {
	tokens       []string
	hash         string
	multiplicity int
}

// NewShingle creates a new Shingle from a given set of tokens
func NewShingle(tokens []string) *Shingle {
	return &Shingle{
		tokens:       tokens,
		hash:         hash(strings.Join(tokens, " ")),
		multiplicity: len(tokens),
	}
}

// GetHash returns the Shingle unique hash
func (s *Shingle) GetHash() string {
	return s.hash
}

// GetMultiplicity returns the multiplicity of the Shingle
func (s *Shingle) GetMultiplicity() int {
	return s.multiplicity
}

type GobShingle struct {
	Tokens       []string
	Hash         string
	Multiplicity int
}

func (s *Shingle) MarshalBinary() ([]byte, error) {
	var buffer bytes.Buffer

	if err := gob.NewEncoder(&buffer).Encode(&GobShingle{
		Tokens:       s.tokens,
		Hash:         s.hash,
		Multiplicity: s.multiplicity,
	}); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (s *Shingle) UnmarshalBinary(b []byte) error {
	var (
		buffer = bytes.NewBuffer(b)
		reader GobShingle
	)

	if err := gob.NewDecoder(buffer).Decode(&reader); err != nil {
		return err
	}

	s.tokens = reader.Tokens
	s.hash = reader.Hash
	s.multiplicity = reader.Multiplicity
	return nil
}
