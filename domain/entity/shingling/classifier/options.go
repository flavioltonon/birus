package classifier

import (
	"bytes"
	"encoding/gob"
)

var _defaultClassifierOptions = classifierOptions{
	tfIdfCutOffThreshold:     0.1,
	scoreNormalizationFactor: 1,
	shinglingMultiplicity:    1,
}

type classifierOptions struct {
	tfIdfCutOffThreshold     float64
	scoreNormalizationFactor float64
	shinglingMultiplicity    int
}

type GobClassifierOptions struct {
	TfIdfCutOffThreshold     float64
	ScoreNormalizationFactor float64
	ShinglingMultiplicity    int
}

func (o *classifierOptions) GobEncode() ([]byte, error) {
	var buffer bytes.Buffer

	if err := gob.NewEncoder(&buffer).Encode(&GobClassifierOptions{
		TfIdfCutOffThreshold:     o.tfIdfCutOffThreshold,
		ScoreNormalizationFactor: o.scoreNormalizationFactor,
		ShinglingMultiplicity:    o.shinglingMultiplicity,
	}); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (o *classifierOptions) GobDecode(b []byte) error {
	var (
		buffer = bytes.NewBuffer(b)
		reader GobClassifierOptions
	)

	if err := gob.NewDecoder(buffer).Decode(&reader); err != nil {
		return err
	}

	o.tfIdfCutOffThreshold = reader.TfIdfCutOffThreshold
	o.scoreNormalizationFactor = reader.ScoreNormalizationFactor
	o.shinglingMultiplicity = reader.ShinglingMultiplicity
	return nil
}
