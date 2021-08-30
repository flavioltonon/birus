package engine

import (
	"io"
	"io/ioutil"
	"os"
)

// Engine is an OCR engine
type Engine interface {
	// ExtractTextFromImage reads text from an image set of bytes
	ExtractTextFromImage(image []byte) (text string, err error)
}

// StoppableEngine is an engine that can be stopped
type StoppableEngine interface {
	Engine

	// Stop stops the engine
	Stop() error
}

// LineTextExtractionEngine is an OCR engine capable of extracting text from an image line by line
type LineTextExtractionEngine interface {
	Engine

	// ExtractTextLinesFromImage reads text from an image set of bytes and returns a set of Blocks
	ExtractTextLinesFromImage(image []byte) ([]Block, error)
}

// Stop stops an Engine
func Stop(e Engine) error {
	if stopper, implements := e.(StoppableEngine); implements {
		return stopper.Stop()
	}

	return nil
}

// ExtractTextFromFile uses an Engine to extract text from a file in a given path
func ExtractTextFromFile(e Engine, path string) (text string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}

	return ExtractTextFromReader(e, f)
}

// ExtractTextFromReader uses an Engine to extract text from a given io.Reader
func ExtractTextFromReader(e Engine, reader io.Reader) (text string, err error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return e.ExtractTextFromImage(b)
}
