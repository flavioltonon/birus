package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/flavioltonon/birus/infrastructure/engine"
	"github.com/flavioltonon/birus/internal/utils"
	"github.com/flavioltonon/birus/pkg/dictionary"
	"github.com/flavioltonon/birus/pkg/shingling"
	"github.com/flavioltonon/birus/pkg/tokeniser"
	"github.com/spf13/viper"
)

func init() {
	// Server config
	viper.SetDefault("SERVER_ADDRESS", ":8000")
	viper.SetDefault("DEVELOPMENT_ENVIRONMENT", true)

	// OCR Engine config
	viper.SetDefault("OCR_ENGINE_TESSDATA_PREFIX", "/home/flavioltonon/Documents/dev/tools/tesseract/tessdata_best")
	viper.SetDefault("OCR_ENGINE_LANGUAGE", "por")
}

// Estabelecimento (nome, CNPJ, endereço)
// Consumidor (nome?, CPF)
// Produtos (nome, valor)
// Total da compra
func main() {
	// e, _ := engine.NewGosseract(engine.GosseractOptions{
	// 	TessdataPrefix: viper.GetString("OCR_ENGINE_TESSDATA_PREFIX"),
	// 	Language:       viper.GetString("OCR_ENGINE_LANGUAGE"),
	// })

	// paths := []string{
	// "./assets/images/1.jpg",
	// "./assets/images/2.jpg",
	// "./assets/images/3.jpg",
	// "./assets/images/4.png",
	// "./assets/images/5.jpg",
	// "./assets/images/6.jpg",
	// "./assets/images/7.jpg",
	// "./assets/images/8.jpg",
	// "./assets/images/9.png",
	// }

	// if err := readAndPersistFilesTexts(e, paths); err != nil {
	// 	log.Fatal(err)
	// }

	paths := []string{
		// "./assets/texts/1.txt",
		// "./assets/texts/2.txt",
		// "./assets/texts/3.txt",
		// "./assets/texts/4.txt",
		"./assets/texts/5.txt",
		// "./assets/texts/6.txt",
		// "./assets/texts/7.txt",
		// "./assets/texts/8.txt",
		// "./assets/texts/9.txt",
	}

	texts, err := utils.ReadTextsFromFiles(paths)
	if err != nil {
		log.Fatal(err)
	}

	allTokens := extractTokensFromTexts(texts)

	for i := range allTokens {
		s1 := shingling.FromTokens(allTokens[i], 1)

		for j := range allTokens {
			if j <= i {
				continue
			}

			s2 := shingling.FromTokens(allTokens[j], 1)

			similarity := shingling.JaccardSimilarity(s1, s2)

			fmt.Printf("paths[i]: %v | paths[j]: %v\n", paths[i], paths[j])
			fmt.Printf("similarity: %v\n", similarity)
		}
	}
}

func extractTextsFromFiles(e engine.Engine, paths []string) ([]string, error) {
	texts := make([]string, 0, len(paths))

	for _, path := range paths {
		text, err := engine.ExtractTextFromFile(e, path)
		if err != nil {
			return nil, err
		}

		texts = append(texts, text)
	}

	return texts, nil
}

func readAndPersistFilesTexts(e engine.Engine, paths []string) error {
	texts, err := extractTextsFromFiles(e, paths)
	if err != nil {
		return err
	}

	for i, text := range texts {
		imageName := filepath.Base(paths[i])

		extension := filepath.Ext(paths[i])

		filename := filepath.Join("./assets/texts", strings.Split(imageName, extension)[0]+".txt")

		if err := utils.WriteTextToFile(text, filename); err != nil {
			return err
		}
	}

	return nil
}

var (
	_modificationChain = tokeniser.NewStringModificationChain(
		tokeniser.RemoveAccents,
		tokeniser.RemoveLineBreaks,
		strings.ToLower,
		tokeniser.RemoveSpecialCharacters,
		tokeniser.RemoveMultipleWhitespaces,
	)

	_stopWords = []string{
		"e", "a", "as", "o", "os", "de", "da", "das", "do", "dos", "em",
	}

	_dictionary = dictionary.New(
		"cnpj", "codigo", "consumidor", "cpf", "credito", "cupom",
		"descricao", "dinheiro",
		"endereco",
		"federal", "fiscal", "fonte",
		"ibpt", "impostos", "item",
		"lei",
		"mensagem",
		"nome",
		"pagos", "produto", "produtos", "promocional",
		"qtd",
		"razao", "rua",
		"sistemas", "social",
		"total",
		"valor",
	)

	_similarityFunc = dictionary.LevenshteinDistanceThreshold(1)
)

func extractTokensFromTexts(texts []string) [][]string {
	allTokens := make([][]string, 0, len(texts))

	for i := range texts {
		texts[i] = _modificationChain.Modify(texts[i])

		tokens := tokeniser.Tokenise(texts[i], _stopWords...)

		for i := range tokens {
			if replacement, found := _dictionary.FindWordBySimilarity(tokens[i], _similarityFunc); found {
				tokens[i] = replacement
			}
		}

		allTokens = append(allTokens, tokens)
	}

	return allTokens
}
