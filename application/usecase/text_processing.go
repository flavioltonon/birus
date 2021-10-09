package usecase

// TextProcessingUsecase are usecases that define operations involving text normalization
type TextProcessingUsecase interface {
	ProcessText(text string) string
}
