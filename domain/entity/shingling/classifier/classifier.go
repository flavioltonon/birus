package classifier

import (
	"bytes"
	"encoding/gob"
	"math"

	"birus/domain/entity/shingling"

	"github.com/google/uuid"
)

// Classifier is a Shingling classifier
type Classifier struct {
	id              string
	name            string
	model           *shingling.Shingling
	shinglings      []*shingling.Shingling
	shinglesMapper  *shingling.ShinglesMapper
	shinglesCounter *shingling.ShinglesCounter
	shinglesTotal   uint16
	shinglingsTotal uint16
	options         classifierOptions
}

// New creates a new Shingling Classifier
func New(name string) *Classifier {
	return &Classifier{
		id:              uuid.NewString(),
		name:            name,
		shinglesMapper:  shingling.NewShinglesMapper(),
		shinglesCounter: shingling.NewShinglesCounter(),
		options:         _defaultClassifierOptions,
	}
}

// ID returns Classifier's unique ID
func (c *Classifier) ID() string {
	return c.id
}

// Name returns the Classifier's name
func (c *Classifier) Name() string {
	return c.name
}

// SetTFIDFCutOffThreshold sets a custom TF-IDF cutoff threshold for the Classifier
func (c *Classifier) SetTFIDFCutOffThreshold(threshold float64) {
	c.options.tfIdfCutOffThreshold = threshold
}

// SetShinglingMultiplicity sets a multiplicity for shinglings in the Classifier
func (c *Classifier) SetShinglingMultiplicity(shinglingMultiplicity int) {
	c.options.shinglingMultiplicity = shinglingMultiplicity
}

// Train trains a Classifier with a set of texts
func (c *Classifier) Train(texts ...string) *Classifier {
	for _, text := range texts {
		c.addShingling(shingling.FromText(text, c.options.shinglingMultiplicity))
	}

	c.model = shingling.FromShingles(c.cutOffShingles())
	c.options.scoreNormalizationFactor = c.calculateScoreNormalizationFactor(texts)
	return c
}

// Classify returns a similarity score by comparing a given Shingling with the Classifier model
func (c *Classifier) Classify(text string) float64 {
	s := shingling.FromText(text, c.options.shinglingMultiplicity)
	return shingling.JaccardSimilarity(c.model, s) * c.options.scoreNormalizationFactor
}

func (c *Classifier) addShingling(s *shingling.Shingling) {
	c.addShingles(s.GetShingles())
	c.shinglings = append(c.shinglings, s)
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

func (c *Classifier) cutOffShingles() []*shingling.Shingle {
	var shingles []*shingling.Shingle

	for hash, tfIdf := range c.calculateTFIDFs() {
		if tfIdf > c.options.tfIdfCutOffThreshold {
			continue
		}

		if shingle, exists := c.shinglesMapper.GetValue(hash); exists {
			shingles = append(shingles, shingle)
		}
	}

	return shingles
}

func (c *Classifier) calculateScoreNormalizationFactor(texts []string) float64 {
	var highestScore float64

	for _, text := range texts {
		if score := c.Classify(text); score > highestScore {
			highestScore = score
		}
	}

	return 1 / highestScore
}

type GobClassifier struct {
	ID              string
	Name            string
	Model           *shingling.Shingling
	ShinglesMapper  *shingling.ShinglesMapper
	ShinglesCounter *shingling.ShinglesCounter
	ShinglesTotal   uint16
	ShinglingsTotal uint16
	Options         classifierOptions
}

func (c *Classifier) GobEncode() ([]byte, error) {
	var buffer bytes.Buffer

	if err := gob.NewEncoder(&buffer).Encode(&GobClassifier{
		ID:              c.id,
		Name:            c.name,
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

	c.id = reader.ID
	c.name = reader.Name
	c.model = reader.Model
	c.shinglesMapper = reader.ShinglesMapper
	c.shinglesCounter = reader.ShinglesCounter
	c.shinglesTotal = reader.ShinglesTotal
	c.shinglingsTotal = reader.ShinglingsTotal
	c.options = reader.Options

	return nil
}
