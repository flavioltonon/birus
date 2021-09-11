package classifier

import (
	"bytes"
	"encoding/gob"
)

var _defaultClassifierOptions = classifierOptions{
	tfIdfCutOffThreshold:     0.1,
	scoreNormalizationFactor: 1,
}

type classifierOptions struct {
	tfIdfCutOffThreshold     float64
	scoreNormalizationFactor float64
}

type classifierOptionFunc func(*classifierOptions)

// SetTFIDFCutOffThreshold is a classifierOptionFunc that can be used to define TF-IDF cutoff threshold for a Classifier.
func SetTFIDFCutOffThreshold(threshold float64) classifierOptionFunc {
	return func(opts *classifierOptions) {
		opts.tfIdfCutOffThreshold = threshold
	}
}

type GobClassifierOptions struct {
	TfIdfCutOffThreshold     float64
	ScoreNormalizationFactor float64
}

func (o *classifierOptions) MarshalBinary() ([]byte, error) {
	var buffer bytes.Buffer

	if err := gob.NewEncoder(&buffer).Encode(&GobClassifierOptions{
		TfIdfCutOffThreshold:     o.tfIdfCutOffThreshold,
		ScoreNormalizationFactor: o.scoreNormalizationFactor,
	}); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (o *classifierOptions) UnmarshalBinary(b []byte) error {
	var (
		buffer = bytes.NewBuffer(b)
		reader GobClassifierOptions
	)

	if err := gob.NewDecoder(buffer).Decode(&reader); err != nil {
		return err
	}

	o.tfIdfCutOffThreshold = reader.TfIdfCutOffThreshold
	o.scoreNormalizationFactor = reader.ScoreNormalizationFactor
	return nil
}
