package service

import (
	"birus/application/usecase"
	"birus/infrastructure/engine"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// OpticalCharacterRecognitionService is a text extraction service
type OpticalCharacterRecognitionService struct {
	imageProcessing usecase.ImageProcessingUsecase
	textProcessing  usecase.TextProcessingUsecase
	options         OpticalCharacterRecognitionServiceOptions
}

// OpticalCharacterRecognitionServiceOptions are options for a OpticalCharacterRecognitionService
type OpticalCharacterRecognitionServiceOptions struct {
	TessdataPrefix string
	Language       string
}

// NewOpticalCharacterRecognitionService creates a new OpticalCharacterRecognitionService
func NewOpticalCharacterRecognitionService(
	imageProcessing usecase.ImageProcessingUsecase,
	textProcessing usecase.TextProcessingUsecase,
	options OpticalCharacterRecognitionServiceOptions,
) usecase.OpticalCharacterRecognitionUsecase {
	return &OpticalCharacterRecognitionService{
		imageProcessing: imageProcessing,
		textProcessing:  textProcessing,
		options:         options,
	}
}

// ReadTextFromImage uses an OCR engine to extract text from a given multipart.FileHeader
func (s *OpticalCharacterRecognitionService) ReadTextFromImage(request *usecase.ReadTextFromImageRequest) (string, error) {
	image, err := s.imageProcessing.ProcessImage(&usecase.ProcessImageRequest{
		Image:   request.Image,
		Options: request.Options,
	})
	if err != nil {
		return "", errors.WithMessage(err, "failed to process image")
	}

	// Gosseract clients are not thread-safe, so we need to start a new client for every call to Tesseract
	e, err := engine.NewGosseract(engine.GosseractOptions{
		TessdataPrefix: s.options.TessdataPrefix,
		Language:       s.options.Language,
	})
	if err != nil {
		return "", errors.WithMessage(err, "failed to initialize OCR engine")
	}

	defer e.Stop()

	text, err := e.ExtractTextFromImage(image.Bytes())
	if err != nil {
		return "", errors.WithMessage(err, "failed to extract text from image")
	}

	return s.textProcessing.ProcessText(text), nil
}

// ReadTextFromImages uses an OCR engine to extract texts from a given set of image.Images
func (s *OpticalCharacterRecognitionService) ReadTextFromImages(request *usecase.ReadTextFromImagesRequest) ([]string, error) {
	var (
		texts = make([]string, 0, len(request.Images))
		g     errgroup.Group
		t     = make(chan string, len(request.Images))
	)

	defer close(t)

	for i := range request.Images {
		image := request.Images[i]

		g.Go(func() error {
			text, err := s.ReadTextFromImage(&usecase.ReadTextFromImageRequest{
				Image: image,
			})
			if err != nil {
				return errors.WithMessage(err, "failed to extract text from image")
			}

			t <- text
			return nil
		})
	}

	go func() {
		for text := range t {
			texts = append(texts, text)
		}
	}()

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return texts, nil
}
