package presenter

import (
	"encoding/base64"

	"birus/domain/entity/image"
)

// Image is a image.Image presenter
type Image struct {
	Base64 string `json:"base64"`
}

// NewImage creates a new Image presenter
func NewImage(image *image.Image) *Image {
	return &Image{
		Base64: base64.StdEncoding.EncodeToString(image.Bytes()),
	}
}
