package text

import (
	"io/ioutil"
	"os"
)

// WriteToFile writes a given text on a new file, created on a given path
func WriteToFile(text, filePath string) error {
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

// ReadFromFiles reads files from a given set of paths and returns their contents
func ReadFromFiles(paths []string) ([]string, error) {
	texts := make([]string, 0, len(paths))

	for _, path := range paths {
		text, err := ReadFromFile(path)
		if err != nil {
			return nil, err
		}

		texts = append(texts, text)
	}

	return texts, nil
}

// ReadFromFile reads a file from a given path and returns its contents
func ReadFromFile(path string) (string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
