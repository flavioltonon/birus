package shingling

import (
	"bytes"
	"encoding/gob"

	"github.com/flavioltonon/birus/internal/mapper"
)

// ShinglesCounter is a counter for shingles occurrencies
type ShinglesCounter struct {
	m mapper.Mapper
}

// NewShinglesCounter creates a new ShinglesCounter
func NewShinglesCounter() *ShinglesCounter {
	return &ShinglesCounter{
		m: make(mapper.Mapper),
	}
}

// Increment adds 1 to the value of a given key in the ShinglesCounter. New keys don't need to be initialized
// with a zero value beforehand.
func (c *ShinglesCounter) Increment(key string) {
	if value, exists := c.GetValue(key); exists {
		c.addValue(key, value+1)
	} else {
		c.addValue(key, 1)
	}
}

func (c *ShinglesCounter) addValue(key string, value uint16) {
	c.m.AddValue(key, value)
}

// GetValue returns the count of a shingle with a given key
func (c *ShinglesCounter) GetValue(key string) (uint16, bool) {
	if value, exists := c.m.GetValue(key); exists {
		return value.(uint16), true
	}

	return 0, false
}

// Each executes a given function iterating over all key/value pairs in the ShinglesCounter
func (c *ShinglesCounter) Each(fn func(key string, value uint16)) {
	c.m.Each(func(k string, v interface{}) { fn(k, v.(uint16)) })
}

// Length returns the length of the ShinglesCounter
func (c *ShinglesCounter) Length() int {
	return len(c.m)
}

// CalculateTermsFrequencies calculates shingles frequencies based on their counts
func (c *ShinglesCounter) CalculateTermsFrequencies() map[string]float64 {
	mostOccurrentTermCount := c.getHighestCount()

	normalizedTFs := make(map[string]float64, len(c.m))

	c.Each(func(key string, value uint16) {
		// calculate the frequency of a given term, applying a augmentation factor to avoid possible bias towards
		// longer documents
		normalizedTFs[key] = 0.5 + 0.5*(float64(value)/float64(mostOccurrentTermCount))
	})

	return normalizedTFs
}

func (c *ShinglesCounter) getHighestCount() uint16 {
	var highestValue uint16

	c.Each(func(key string, value uint16) {
		if value > highestValue {
			highestValue = value
		}
	})

	return highestValue
}

type GobShinglesCounter struct {
	M mapper.Mapper
}

func (m *ShinglesCounter) GobEncode() ([]byte, error) {
	var buffer bytes.Buffer

	if err := gob.NewEncoder(&buffer).Encode(&GobShinglesCounter{M: m.m}); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (m *ShinglesCounter) GobDecode(b []byte) error {
	var (
		buffer = bytes.NewBuffer(b)
		reader GobShinglesCounter
	)

	if err := gob.NewDecoder(buffer).Decode(&reader); err != nil {
		return err
	}

	m.m = reader.M
	return nil
}
