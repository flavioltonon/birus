package normalization

// chain is a chain of string normalizers
type Chain []normalizer

// NewChain creates a new normalization chain
func NewChain(fns ...normalizer) Chain {
	return Chain(fns)
}

// Normalize normalizes a document according to the normalizers
func (c Chain) Normalize(document string) string {
	for _, normalizer := range c {
		document = normalizer.normalize(document)
	}

	return document
}
