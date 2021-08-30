package engine

import (
	ozzo "github.com/go-ozzo/ozzo-validation"
	"github.com/otiai10/gosseract/v2"
	"github.com/pkg/errors"
)

// Gosseract wraps gosseract.Client to make it implement Engine interface
type Gosseract struct {
	source *gosseract.Client
}

// GosseractOptions are options for Gosseract
type GosseractOptions struct {
	// TessdataPrefix is the path to Tesseract models directory.
	TessdataPrefix string

	// Language is the language that Tesseract should use
	Language string
}

func (opts GosseractOptions) validate() error {
	return ozzo.ValidateStruct(&opts,
		ozzo.Field(&opts.TessdataPrefix, ozzo.Required),
		ozzo.Field(&opts.Language, ozzo.Required),
	)
}

func (opts GosseractOptions) apply(engine *Gosseract) error {
	if err := opts.validate(); err != nil {
		return errors.WithMessage(err, "failed to validate options")
	}

	if err := engine.source.SetTessdataPrefix(opts.TessdataPrefix); err != nil {
		return errors.WithMessage(err, "failed to set Tessdata prefix")
	}

	if err := engine.source.SetLanguage(opts.Language); err != nil {
		return errors.WithMessage(err, "failed to set language")
	}

	return nil
}

// NewGosseract creates a new Gosseract
func NewGosseract(opts GosseractOptions) (*Gosseract, error) {
	engine := &Gosseract{
		source: gosseract.NewClient(),
	}

	if err := opts.apply(engine); err != nil {
		return nil, errors.WithMessage(err, "failed to apply options")
	}

	return engine, nil
}

// Close closes Gosseract engine connection with Tesseract's API
func (e *Gosseract) Stop() error {
	return e.source.Close()
}

// ExtractTextFromImage makes the receiver implement Engine interface
func (e *Gosseract) ExtractTextFromImage(image []byte) (text string, err error) {
	if err := e.source.SetImageFromBytes(image); err != nil {
		return "", errors.WithMessage(err, "failed to set image from bytes")
	}

	return e.source.Text()
}

// ExtractTextFromImage makes the receiver implement LineTextExtractionEngine interface
func (e *Gosseract) ExtractTextLinesFromImage(image []byte) ([]Block, error) {
	if err := e.source.SetImageFromBytes(image); err != nil {
		return nil, errors.WithMessage(err, "failed to set image from bytes")
	}

	boxes, err := e.source.GetBoundingBoxes(gosseract.RIL_TEXTLINE)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get bounding boxes from image")
	}

	return newBlocksFromGosseractBoundingBoxes(boxes), nil
}
