package tokeniser

import "strings"

// Tokeniser is a document tokeniser
type Tokeniser struct {
	stopWords []string
}

// New creates a new Tokeniser with an option set of stopwords
func New(stopWords ...string) *Tokeniser {
	return &Tokeniser{stopWords: stopWords}
}

// Tokenise extracts tokens from a document
func (t *Tokeniser) Tokenise(document string) []string {
	var tokens []string

	stopWordsMap := mapWords(t.stopWords)

	for _, token := range strings.Split(document, " ") {
		if _, exists := stopWordsMap[token]; exists {
			continue
		}

		tokens = append(tokens, token)
	}

	return tokens
}

func mapWords(words []string) map[string]struct{} {
	wordsMap := make(map[string]struct{}, len(words))

	for _, stopWord := range words {
		wordsMap[stopWord] = struct{}{}
	}

	return wordsMap
}
