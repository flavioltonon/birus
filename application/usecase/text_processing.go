package usecase

// TextProcessingUsecase are usecases that define operations involving text normalization
type TextProcessingUsecase interface {
	NormalizeText(text string) string
	TokeniseText(texts string) []string
	FixWords(words ...string) []string
}
