package service

import (
	"context"
	"mime/multipart"
	"strings"

	"github.com/flavioltonon/birus/application/usecase"
	"github.com/flavioltonon/birus/domain/entity"
	"github.com/flavioltonon/birus/infrastructure/engine"
	"github.com/flavioltonon/birus/internal/dictionary"
	"github.com/flavioltonon/birus/internal/normalization"
	"github.com/flavioltonon/birus/internal/shingling"
	"github.com/flavioltonon/birus/internal/shingling/classifier"
	"github.com/flavioltonon/birus/internal/tokeniser"

	ozzo "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
)

// Ensure that ImageClassificationService implements usecase.TypificationUsecase
var _ usecase.ImageClassificationUsecase = (*ImageClassificationService)(nil)

// ImageClassificationService  interface
type ImageClassificationService struct {
	models     usecase.ModelsUsecase
	engine     engine.Engine
	normalizer normalization.Chain
	tokeniser  *tokeniser.Tokeniser
	dictionary *dictionary.Dictionary
}

// NewImageClassificationService creates new use case
func NewImageClassificationService(m usecase.ModelsUsecase, e engine.Engine) *ImageClassificationService {
	return &ImageClassificationService{
		models: m,
		engine: e,
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
			"caixa", "cartao", "chave", "cnpj", "codigo", "comprovante", "consulta", "consumidor", "cpf", "credito", "cupom",
			"desconto", "descontos", "descricao", "dinheiro",
			"economizou", "eletronica", "endereco", "estadual",
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

func (s *ImageClassificationService) newShinglingFromText(text string) *shingling.Shingling {
	normalizedText := s.normalizer.Normalize(text)

	tokens := s.dictionary.ReplaceWordsBySimilarity(
		s.tokeniser.Tokenise(normalizedText),
		dictionary.LevenshteinDistance(1),
	)

	return shingling.FromTokens(tokens, 1)
}

// CreateClassificationModel creates a new typification model for a given name and a set of images
func (s *ImageClassificationService) CreateClassificationModel(ctx context.Context, name string, files []*multipart.FileHeader) (string, error) {
	if err := ozzo.Required.Validate(files); err != nil {
		return "", errors.WithMessage(err, "failed to validate files")
	}

	images, err := entity.NewImages(files)
	if err != nil {
		return "", errors.WithMessage(err, "failed to read images from files")
	}

	shinglings := make([]*shingling.Shingling, 0, len(images))

	for _, image := range images {
		text, err := s.engine.ExtractTextFromImage(image.Bytes())
		if err != nil {
			return "", errors.WithMessage(err, "failed to extract text from image")
		}

		shinglings = append(shinglings, s.newShinglingFromText(text))
	}

	return s.models.CreateModel(ctx, name, shinglings)
}

func (s *ImageClassificationService) ClassifyImage(ctx context.Context, file *multipart.FileHeader) (string, error) {
	image, err := entity.NewImage(file)
	if err != nil {
		return "", errors.WithMessage(err, "failed to read image from file")
	}

	models, err := s.models.ListModels(ctx)
	if err != nil {
		return "", errors.WithMessage(err, "failed to list models")
	}

	set := classifier.NewSet()

	for _, model := range models {
		set.AddClassifier(model.Classifier)
	}

	text, err := s.engine.ExtractTextFromImage(image.Bytes())
	if err != nil {
		return "", errors.WithMessage(err, "failed to extract text from image")
	}

	var (
		highestScore           float64
		highestScoreClassifier string
	)

	scores := set.Classify(s.newShinglingFromText(text))

	for classifierName, score := range scores {
		if score > highestScore {
			highestScore = score
			highestScoreClassifier = classifierName
		}
	}

	return highestScoreClassifier, nil
}
