package image

import (
	"io/ioutil"
	"mime/multipart"

	"github.com/pkg/errors"
)

// Image is a set of bytes
type Image struct {
	data []byte
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

	return &Image{data: data}, nil
}

// Bytes returns Image's content as a set of bytes
func (i *Image) Bytes() []byte {
	return i.data
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
