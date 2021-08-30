package utils

import (
	"io/ioutil"
	"os"
)

// WriteTextToFile writes a given text on a new file, created on a given path
func WriteTextToFile(text, filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err := f.WriteString(text); err != nil {
		return err
	}

	return nil
}

// ReadTextsFromFiles reads files from a given set of paths and returns their contents
func ReadTextsFromFiles(paths []string) ([]string, error) {
	texts := make([]string, 0, len(paths))

	for _, path := range paths {
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}

		texts = append(texts, string(b))
	}

	return texts, nil
}
