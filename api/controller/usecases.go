package controller

import "birus/application/usecase"

type Usecases struct {
	ImageProcessing             usecase.ImageProcessingUsecase
	OpticalCharacterRecognition usecase.OpticalCharacterRecognitionUsecase
	TextClassification          usecase.TextClassificationUsecase
}
