package service

import (
	"birus/application/usecase"
	"birus/domain/entity/dictionary"
	"birus/domain/entity/normalization"
	"birus/domain/entity/tokeniser"
	"strings"
)

// TextProcessingService is a text normalization service
type TextProcessingService struct {
	normalizer normalization.Chain
	tokeniser  *tokeniser.Tokeniser
	dictionary *dictionary.Dictionary
}

// NewTextProcessingService creates a new TextProcessingService
func NewTextProcessingService() usecase.TextProcessingUsecase {
	return &TextProcessingService{
		normalizer: normalization.NewChain(
			normalization.RemoveAccents,
			normalization.RemoveLineBreaks,
			strings.ToLower,
			normalization.RemoveSpecialCharacters,
			normalization.RemoveMultipleWhitespaces,
		),
		tokeniser: tokeniser.New("e", "a", "as", "o", "os", "de", "da", "das", "do", "dos", "em"),
		dictionary: dictionary.New(
			"acesso", "auxiliar",
			"caixa", "cartao", "chave", "cnpj", "codigo", "comete", "comprovante", "consulta", "consumidor", "cpf", "credito", "crime", "cupom",
			"desconto", "descontos", "descricao", "dinheiro",
			"economizou", "eletronico", "eletronica", "endereco", "estadual", "extrato",
			"federal", "fiscal", "fonte", "forma",
			"ibpt", "impostos", "incidentes", "item", "itens",
			"lei", "loja",
			"mensagem", "municipal",
			"nome", "nota",
			"pagamento", "pagos", "produto", "produtos", "promocional",
			"qtd",
			"razao", "rua",
			"sefaz", "sistemas", "social",
			"totais", "total", "tributos", "troco",
			"valor", "venda",
		),
	}
}

// NormalizeText applies normalization functions over an input text
func (s *TextProcessingService) NormalizeText(text string) string {
	return s.normalizer.Normalize(text)
}

// TokeniseText tokenises an input text
func (s *TextProcessingService) TokeniseText(text string) []string {
	return s.tokeniser.Tokenise(text)
}

// FixWords tries to return the best match for all words in a given set by looking up into the dictionary
// for any words with a high level of similarity. If the word is known by the dictionary, the same word
// will be returned.
func (s *TextProcessingService) FixWords(words ...string) []string {
	result := make([]string, 0, len(words))

	for _, word := range words {
		w, _ := s.dictionary.FindWordBySimilarity(word, dictionary.LevenshteinDistance(1))
		result = append(result, w)
	}

	return result
}

func (s *TextProcessingService) ProcessText(text string) string {
	return strings.Join(s.FixWords(s.TokeniseText(s.NormalizeText(text))...), " ")
}
