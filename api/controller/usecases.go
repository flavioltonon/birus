package controller

import "birus/application/usecase"

type Usecases struct {
	ImageClassification         usecase.ImageClassificationUsecase
	OpticalCharacterRecognition usecase.OpticalCharacterRecognitionUsecase
}
