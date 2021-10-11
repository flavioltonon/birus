package image

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/gabriel-vasile/mimetype"
	"github.com/pkg/errors"
)

// Image is a set of bytes
type Image struct {
	data     []byte
	mimetype *mimetype.MIME
}

// New creates a new Image
func New(data []byte) *Image {
	return &Image{
		data:     data,
		mimetype: mimetype.Detect(data),
	}
}

// FromMultipartFileHeader creates an Image from a given *multipart.FileHeader
func FromMultipartFileHeader(fileHeader *multipart.FileHeader) (*Image, error) {
	f, err := fileHeader.Open()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to open file")
	}

	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read data from file")
	}

	return New(data), nil
}

// FromMultipartFileHeaders returns a set of Images for a given set of *multipart.FileHeader
func FromMultipartFileHeaders(fileHeaders []*multipart.FileHeader) ([]*Image, error) {
	images := make([]*Image, 0, len(fileHeaders))

	for _, fileHeader := range fileHeaders {
		image, err := FromMultipartFileHeader(fileHeader)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to read image from multipart file header")
		}

		images = append(images, image)
	}

	return images, nil
}

// Bytes returns Image's content as a set of bytes
func (i *Image) Bytes() []byte {
	return i.data
}

type ProcessOptionFunc func(img image.Image) *image.NRGBA

func (fn ProcessOptionFunc) apply(img image.Image) *image.NRGBA { return fn(img) }

// Process processes an Image with a given set of ProcessOptionFunc. If no options are provided, the
// original image will be returned.
//
// Example:
//
//  image := New([]byte{/* image bytes */})
//
//  processedImage, err := image.Process(Resize(1080, 720), Grayscale())
func (i *Image) Process(opts ...ProcessOptionFunc) (*Image, error) {
	if len(opts) == 0 {
		return i, nil
	}

	extension, err := imaging.FormatFromExtension(i.mimetype.Extension())
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(i.Bytes())

	image, err := imaging.Decode(reader)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		image = opt.apply(image)
	}

	var buffer bytes.Buffer

	if err := imaging.Encode(&buffer, image, extension); err != nil {
		return nil, err
	}

	return New(buffer.Bytes()), nil
}

// Save saves an image to a given file path
func (i *Image) Save(path string) error {
	reader := bytes.NewReader(i.Bytes())

	image, err := imaging.Decode(reader)
	if err != nil {
		return err
	}

	return imaging.Save(image, path)
}

// ParseProcessOptions reads a string of options and parses it into a list of ProcessOptionFuncs.
//
// Example:
//
//  options := ParseProcessOptions("grayscale;sharpen:3.5")
func ParseProcessOptions(optionsStr string) ([]ProcessOptionFunc, error) {
	optionsStrs := strings.Split(optionsStr, ";")

	fns := make([]ProcessOptionFunc, 0, len(optionsStrs))

	for _, option := range optionsStrs {
		fn, err := parseProcessOption(option)
		if err != nil {
			return nil, err
		}

		fns = append(fns, fn)
	}

	return fns, nil
}

// parseProcessOption reads a string and parses it into a single ProcessOptionFunc.
//
// Example:
//
//  options := ParseProcessOption("sharpen:3.5")
func parseProcessOption(optionStr string) (ProcessOptionFunc, error) {
	split := strings.Split(optionStr, ":")

	switch split[0] {
	case "blur":
		sigma, err := strconv.ParseFloat(split[1], 64)
		if err != nil {
			return nil, errors.New("sigma should be a number")
		}

		return GaussianBlur(sigma), nil
	case "resize":
		if len(split) != 2 {
			return nil, errors.New("resize requires a weight and a height")
		}

		dimensions := strings.Split(split[1], ",")

		if len(dimensions) != 2 {
			return nil, errors.New("invalid number of dimensions")
		}

		w, err := strconv.Atoi(dimensions[0])
		if err != nil {
			return nil, errors.New("dimensions should be integers")
		}

		h, err := strconv.Atoi(dimensions[1])
		if err != nil {
			return nil, errors.New("dimensions should be integers")
		}

		return Resize(w, h), nil
	case "grayscale":
		return Grayscale(), nil
	case "adjust-contrast":
		percentage, err := strconv.ParseFloat(split[1], 64)
		if err != nil {
			return nil, errors.New("percentage should be a number")
		}

		return AdjustContrast(percentage), nil
	case "adjust-brightness":
		percentage, err := strconv.ParseFloat(split[1], 64)
		if err != nil {
			return nil, errors.New("percentage should be a number")
		}

		return AdjustBrightness(percentage), nil
	case "sharpen":
		sigma, err := strconv.ParseFloat(split[1], 64)
		if err != nil {
			return nil, errors.New("sigma should be a number")
		}

		return Sharpen(sigma), nil
	case "":
		return nil, errors.New("option string cannot be empty")
	default:
		return nil, fmt.Errorf("invalid option string '%s'", split[0])
	}
}

// Resize resizes an image to a given width and height in pixels
func Resize(width int, height int) ProcessOptionFunc {
	return func(img image.Image) *image.NRGBA {
		if img.Bounds().Dx() > img.Bounds().Dy() {
			height = 0
		} else {
			width = 0
		}

		return imaging.Resize(img, width, height, imaging.Lanczos)
	}
}

// Grayscale transforms the image colors to shades of grey
func Grayscale() ProcessOptionFunc {
	return func(img image.Image) *image.NRGBA {
		return imaging.Grayscale(img)
	}
}

// AdjustContrast sets the contrast in the image to a given percentage
func AdjustContrast(percentage float64) ProcessOptionFunc {
	return func(img image.Image) *image.NRGBA {
		return imaging.AdjustContrast(img, percentage)
	}
}

// AdjustBrightness sets the brightness in the image to a given percentage
func AdjustBrightness(percentage float64) ProcessOptionFunc {
	return func(img image.Image) *image.NRGBA {
		return imaging.AdjustBrightness(img, percentage)
	}
}

// Sharpen produces a sharpened version of the image. Sigma should be a positive number.
func Sharpen(sigma float64) ProcessOptionFunc {
	return func(img image.Image) *image.NRGBA {
		return imaging.Sharpen(img, sigma)
	}
}

// GaussianBlur produces a blurred version of the image. Sigma should be a positive number.
func GaussianBlur(sigma float64) ProcessOptionFunc {
	return func(img image.Image) *image.NRGBA {
		return imaging.Blur(img, sigma)
	}
}
