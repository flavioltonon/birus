package tokeniser

import "strings"

// Tokenise extracts tokens from a document
func Tokenise(document string, stopWords ...string) []string {
	var tokens []string

	stopWordsMap := mapWords(stopWords)

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
