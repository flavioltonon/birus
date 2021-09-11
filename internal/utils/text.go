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
		text, err := ReadTextFromFile(path)
		if err != nil {
			return nil, err
		}

		texts = append(texts, text)
	}

	return texts, nil
}

// ReadTextFromFile reads a file from a given path and returns its contents
func ReadTextFromFile(path string) (string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
