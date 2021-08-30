package controller

import "github.com/flavioltonon/birus/application/usecase"

type Usecases struct {
	Models              usecase.ModelsUsecase
	ImageClassification usecase.ImageClassificationUsecase
	// DataExtraction      usecase.DataExtraction
}
