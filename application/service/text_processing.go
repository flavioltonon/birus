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
			normalization.IsolateLineBreaks,
			strings.ToLower,
			normalization.RemoveSpecialCharacters,
			normalization.RemoveMultipleWhitespaces,
		),
		tokeniser: tokeniser.New(),
		dictionary: dictionary.New(
			"acesso", "auxiliar", "avenida",
			"bairro", "bermuda", "brasil",
			"cadastre", "cadastro", "caixa", "calca", "camiseta", "carros", "cartao", "cartoes", "chave", "cidade", "cnpj", "cod", "codigo", "comete", "compra", "compras", "comprovante", "comprovantes", "concorre", "consulta", "consumidor", "cpf", "credito", "crime", "cupom",
			"debito", "desc", "desconto", "descontos", "descricao", "dinheiro", "documento",
			"economizou", "eletronico", "eletronica", "emit", "emitida", "endereco", "estadual", "extrato",
			"federal", "fem", "fiscal", "fonte", "forma",
			"ibpt", "identificado", "identificados", "imposto", "impostos", "incidente", "incidentes", "inf", "item", "itens",
			"lei", "leis", "loja", "lojas",
			"macaquinho", "maguineta", "masc", "mensagem", "municipal", "municipais",
			"nao", "natal", "nome", "nota",
			"pagamento", "pagos", "pijama", "produto", "produtos", "promocional", "promocionais",
			"qtd", "quem",
			"razao", "regata", "rua",
			"sefaz", "shorts", "sistema", "sistemas", "social",
			"totais", "total", "tributo", "tributos", "troco",
			"valor", "venda", "vendas",
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
