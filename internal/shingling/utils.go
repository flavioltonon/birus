package shingling

import (
	"crypto/sha256"
	"encoding/base64"
	"math"
)

// TrainTestSplit splits a set of Shinglings into two separate sets for training and testing shingling classifiers
func TrainTestSplit(shinglings []*Shingling, trainPercentage float64) (train []*Shingling, test []*Shingling) {
	if trainPercentage > 1 {
		panic("train percentage should range between 0-1")
	}

	trainSize := int(math.Ceil(float64(len(shinglings)) * trainPercentage))

	return shinglings[:trainSize], shinglings[trainSize:]
}

// hash encodes a given string to base64 after hashing it using a SHA256 algorithm
func hash(s string) string {
	hash := sha256.New()
	hash.Write([]byte(s))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}
