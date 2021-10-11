package controller

import (
	"encoding/base64"

	"birus/application/usecase"
	"birus/domain/entity/image"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (c *Controller) newCreateClassifierRequest(ctx *gin.Context) (*usecase.CreateClassifierRequest, error) {
	var request usecase.CreateClassifierRequest

	if err := ctx.BindJSON(&request); err != nil {
		return nil, errors.WithMessage(err, "failed to decode request body")
	}

	if err := request.Validate(); err != nil {
		return nil, errors.WithMessage(err, "failed to validate request body")
	}

	return &request, nil
}

func (c *Controller) newListClassifiersRequest(ctx *gin.Context) (*usecase.ListClassifiersRequest, error) {
	var request usecase.ListClassifiersRequest

	if err := request.Validate(); err != nil {
		return nil, errors.WithMessage(err, "failed to validate request body")
	}

	return &request, nil
}

func (c *Controller) newDeleteClassifierRequest(ctx *gin.Context) (*usecase.DeleteClassifierRequest, error) {
	request := usecase.DeleteClassifierRequest{
		ID: ctx.Param("classifier_id"),
	}

	if err := request.Validate(); err != nil {
		return nil, errors.WithMessage(err, "failed to validate request body")
	}

	return &request, nil
}

func (c *Controller) newClassifyTextRequest(ctx *gin.Context) (*usecase.ClassifyTextRequest, error) {
	var request usecase.ClassifyTextRequest

	if err := ctx.BindJSON(&request); err != nil {
		return nil, errors.WithMessage(err, "failed to decode request body")
	}

	if err := request.Validate(); err != nil {
		return nil, errors.WithMessage(err, "failed to validate request body")
	}

	return &request, nil
}

func (c *Controller) newReadTextFromImageRequest(ctx *gin.Context) (*usecase.ReadTextFromImageRequest, error) {
	var request usecase.ReadTextFromImageRequest

	switch ctx.ContentType() {
	case "application/json":
		wrapper := new(struct {
			Base64  string `json:"base64"`
			Options string `json:"options"`
		})

		if err := ctx.BindJSON(wrapper); err != nil {
			return nil, errors.WithMessage(err, "failed to decode JSON body")
		}

		raw, err := base64.StdEncoding.DecodeString(wrapper.Base64)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to decode base64 image data")
		}

		request.Image = image.FromBytes(raw)

		if wrapper.Options != "" {
			request.Options, err = image.ParseProcessOptions(wrapper.Options)
			if err != nil {
				return nil, errors.WithMessage(err, "failed to parse process options")
			}
		}
	case "multipart/form-data":
		file, err := ctx.FormFile("file")
		if err != nil {
			return nil, errors.WithMessage(err, "failed to parse file from multipart form")
		}

		request.Image, err = image.FromMultipartFileHeader(file)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to read image from file")
		}

		request.Options, err = image.ParseProcessOptions(ctx.Request.FormValue("options"))
		if err != nil {
			return nil, errors.WithMessage(err, "failed to parse process options")
		}
	}

	if err := request.Validate(); err != nil {
		return nil, errors.WithMessage(err, "failed to validate request body")
	}

	return &request, nil
}

func (c *Controller) newReadTextFromImagesRequest(ctx *gin.Context) (*usecase.ReadTextFromImagesRequest, error) {
	var (
		request usecase.ReadTextFromImagesRequest
		err     error
	)

	switch ctx.ContentType() {
	case "application/json":
		wrapper := new(struct {
			Base64List []string `json:"base64_list"`
			Options    string   `json:"options"`
		})

		if err := ctx.BindJSON(wrapper); err != nil {
			return nil, errors.WithMessage(err, "failed to decode JSON body")
		}

		request.Images, err = image.FromBase64List(wrapper.Base64List)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to read images from base64 list")
		}

		request.Options, err = image.ParseProcessOptions(wrapper.Options)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to parse process options")
		}
	case "multipart/form-data":
		form, err := ctx.MultipartForm()
		if err != nil {
			return nil, errors.WithMessage(err, "failed to parse file from multipart form")
		}

		request.Images, err = image.FromMultipartFileHeaders(form.File["files"])
		if err != nil {
			return nil, errors.WithMessage(err, "failed to read image from file")
		}

		request.Options, err = image.ParseProcessOptions(ctx.Request.FormValue("options"))
		if err != nil {
			return nil, errors.WithMessage(err, "failed to parse process options")
		}
	}

	if err := request.Validate(); err != nil {
		return nil, errors.WithMessage(err, "failed to validate request body")
	}

	return &request, nil
}

func (c *Controller) newProcessImageRequest(ctx *gin.Context) (*usecase.ProcessImageRequest, error) {
	var request usecase.ProcessImageRequest

	switch ctx.ContentType() {
	case "application/json":
		wrapper := new(struct {
			Base64  string `json:"base64"`
			Options string `json:"options"`
		})

		if err := ctx.BindJSON(wrapper); err != nil {
			return nil, errors.WithMessage(err, "failed to decode JSON body")
		}

		raw, err := base64.StdEncoding.DecodeString(wrapper.Base64)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to decode base64 image data")
		}

		request.Image = image.FromBytes(raw)

		request.Options, err = image.ParseProcessOptions(wrapper.Options)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to parse process options")
		}
	case "multipart/form-data":
		file, err := ctx.FormFile("file")
		if err != nil {
			return nil, errors.WithMessage(err, "failed to parse file from multipart form")
		}

		request.Image, err = image.FromMultipartFileHeader(file)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to read image from file")
		}

		request.Options, err = image.ParseProcessOptions(ctx.Request.FormValue("options"))
		if err != nil {
			return nil, errors.WithMessage(err, "failed to parse process options")
		}
	}

	if err := request.Validate(); err != nil {
		return nil, errors.WithMessage(err, "failed to validate request body")
	}

	return &request, nil
}
