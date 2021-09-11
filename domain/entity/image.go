package entity

import (
	"io/ioutil"
	"mime/multipart"

	"github.com/pkg/errors"
)

// Image is a set of bytes
type Image struct {
	data []byte
}

// NewImages returns a set of Images for a given set of *multipart.FileHeader
func NewImages(fhs []*multipart.FileHeader) ([]*Image, error) {
	images := make([]*Image, 0, len(fhs))

	for _, fh := range fhs {
		image, err := NewImage(fh)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to read image from multipart file header")
		}

		images = append(images, image)
	}

	return images, nil
}

// NewImage returns an Image for a given *multipart.FileHeader
func NewImage(fh *multipart.FileHeader) (*Image, error) {
	f, err := fh.Open()
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
