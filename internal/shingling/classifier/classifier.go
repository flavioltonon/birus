package classifier

import (
	"bytes"
	"encoding/gob"
	"math"

	"github.com/flavioltonon/birus/internal/shingling"
)

// Classifier is a Shingling classifier
type Classifier struct {
	ID              string
	model           *shingling.Shingling
	shinglesMapper  *shingling.ShinglesMapper
	shinglesCounter *shingling.ShinglesCounter
	shinglesTotal   uint16
	shinglingsTotal uint16
	options         *classifierOptions
}

// New creates a new Shingling classifier with a given TF-IDF cut-off function.
func New(id string, funcs ...classifierOptionFunc) *Classifier {
	options := _defaultClassifierOptions

	for _, fn := range funcs {
		fn(&options)
	}

	return &Classifier{
		ID:              id,
		shinglesMapper:  shingling.NewShinglesMapper(),
		shinglesCounter: shingling.NewShinglesCounter(),
		options:         &options,
	}
}

// Train trains a Classifier with a set of Shinglings
func (c *Classifier) Train(shinglings ...*shingling.Shingling) *Classifier {
	for i := range shinglings {
		c.addShingling(shinglings[i])
	}

	var shingles []*shingling.Shingle

	for hash, tfIdf := range c.calculateTFIDFs() {
		if tfIdf > c.options.tfIdfCutOffThreshold {
			continue
		}

		if shingle, exists := c.shinglesMapper.GetValue(hash); exists {
			shingles = append(shingles, shingle)
		}
	}

	c.model = shingling.FromShingles(shingles)

	var highestScore float64

	for i := range shinglings {
		score := c.Classify(shinglings[i])

		if score > highestScore {
			highestScore = score
		}
	}

	c.options.scoreNormalizationFactor = 1 / highestScore

	return c
}

// Classify returns a similarity score by comparing a given Shingling with the Classifier model
func (c *Classifier) Classify(s *shingling.Shingling) float64 {
	return shingling.JaccardSimilarity(c.model, s) * c.options.scoreNormalizationFactor
}

func (c *Classifier) addShingling(s *shingling.Shingling) {
	c.addShingles(s.GetShingles())
	c.shinglingsTotal++
}

func (c *Classifier) addShingles(shingles []*shingling.Shingle) {
	for i := range shingles {
		c.addShingle(shingles[i])
	}
}

func (c *Classifier) addShingle(s *shingling.Shingle) {
	hash := s.GetHash()
	c.shinglesMapper.AddValue(hash, s)
	c.shinglesCounter.Increment(hash)
	c.shinglesTotal++
}

func (c *Classifier) calculateTFIDFs() map[string]float64 {
	var (
		tfs  = c.calculateTermsFrequencies()
		idfs = c.calculateInverseDocumentFrequencies()
	)

	tfIdfs := make(map[string]float64)

	c.shinglesMapper.Each(func(key string, value *shingling.Shingle) {
		tfIdfs[key] = tfs[key] * idfs[key]
	})

	return tfIdfs
}

func (c *Classifier) calculateTermsFrequencies() map[string]float64 {
	return c.shinglesCounter.CalculateTermsFrequencies()
}

func (c *Classifier) calculateInverseDocumentFrequencies() map[string]float64 {
	idfs := make(map[string]float64, c.shinglesCounter.Length())

	c.shinglesCounter.Each(func(key string, value uint16) {
		idfs[key] = math.Log(float64(c.shinglingsTotal) / float64(value))
	})

	return idfs
}

type GobClassifier struct {
	ID              string
	Model           *shingling.Shingling
	ShinglesMapper  *shingling.ShinglesMapper
	ShinglesCounter *shingling.ShinglesCounter
	ShinglesTotal   uint16
	ShinglingsTotal uint16
	Options         *classifierOptions
}

func (c *Classifier) GobEncode() ([]byte, error) {
	var buffer bytes.Buffer

	if err := gob.NewEncoder(&buffer).Encode(&GobClassifier{
		ID:              c.ID,
		Model:           c.model,
		ShinglesMapper:  c.shinglesMapper,
		ShinglesCounter: c.shinglesCounter,
		ShinglesTotal:   c.shinglesTotal,
		ShinglingsTotal: c.shinglingsTotal,
		Options:         c.options,
	}); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (c *Classifier) GobDecode(b []byte) error {
	var (
		buffer = bytes.NewBuffer(b)
		reader GobClassifier
	)

	if err := gob.NewDecoder(buffer).Decode(&reader); err != nil {
		return err
	}

	c.ID = reader.ID
	c.model = reader.Model
	c.shinglesMapper = reader.ShinglesMapper
	c.shinglesCounter = reader.ShinglesCounter
	c.shinglesTotal = reader.ShinglesTotal
	c.shinglingsTotal = reader.ShinglingsTotal
	c.options = reader.Options

	return nil
}
