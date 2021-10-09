package controller

import "birus/application/usecase"

type Usecases struct {
	OpticalCharacterRecognition usecase.OpticalCharacterRecognitionUsecase
	TextClassification          usecase.TextClassificationUsecase
}
