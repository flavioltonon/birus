package tokeniser

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	// _accentsRemover is a string transformer that removes accents from strings
	_accentsRemover = transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)

	// _specialCharactersMatcher is a regular expression to match all characters that are not letters or whitespaces
	_specialCharactersMatcher = regexp.MustCompile(`[^a-z\s]+`)

	// _multipleWhitespaceMatcher is a regular expression to match all occurrences of multiple sequential whitespaces
	_multipleWhitespaceMatcher = regexp.MustCompile(`\s+`)
)

type stringModificationChain []stringModificationFunc

// NewStringModificationChain creates a new chain of stringModificationFunc
func NewStringModificationChain(fns ...stringModificationFunc) stringModificationChain {
	return fns
}

// Modify modifies a string based on the modificationFuncs in the chain
func (c stringModificationChain) Modify(s string) string {
	for _, modifier := range c {
		s = modifier.modify(s)
	}

	return s
}

type stringModificationFunc func(s string) string

func (fn stringModificationFunc) modify(s string) string {
	return fn(s)
}

// RemoveAccents removes accents from a given string
func RemoveAccents(s string) string {
	s, _, err := transform.String(_accentsRemover, s)
	if err != nil {
		panic(err)
	}

	return s
}

// RemoveLineBreaks replaces line breaks with whitespaces
func RemoveLineBreaks(s string) string {
	return strings.Replace(s, "\n", " ", -1)
}

// RemoveMultipleWhitespaces replaces multiple sequential whitespaces with a single one
func RemoveMultipleWhitespaces(s string) string {
	return _multipleWhitespaceMatcher.ReplaceAllString(s, " ")
}

// RemoveSpecialCharacters removes special characters from a given string
func RemoveSpecialCharacters(s string) string {
	return _specialCharactersMatcher.ReplaceAllString(s, " ")
}
