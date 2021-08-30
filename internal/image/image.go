package image

import (
	"io/ioutil"
	"mime/multipart"

	"github.com/pkg/errors"
)

// Image is a set of bytes
type Image []byte

// Bytes returns Image's content as a set of bytes
func (i Image) Bytes() []byte {
	return []byte(i)
}

// FromBulkMultipartFileHeaders returns a set of Images for a given set of *multipart.FileHeader
func FromBulkMultipartFileHeaders(fhs []*multipart.FileHeader) ([]Image, error) {
	images := make([]Image, 0, len(fhs))

	for _, fh := range fhs {
		image, err := FromMultipartFileHeader(fh)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to read image from multipart file header")
		}

		images = append(images, image)
	}

	return images, nil
}

// FromMultipartFileHeader returns an Image for a given *multipart.FileHeader
func FromMultipartFileHeader(fh *multipart.FileHeader) (Image, error) {
	f, err := fh.Open()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to open file")
	}

	defer f.Close()

	image, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read data from file")
	}

	return image, nil
}
