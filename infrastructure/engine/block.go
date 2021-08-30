package engine

import "github.com/otiai10/gosseract/v2"

// Block is a block of text extracted from an image
type Block struct {
	Text string
}

func newBlockFromGosseractBoundingBox(box gosseract.BoundingBox) Block {
	return Block{
		Text: box.Word,
	}
}

func newBlocksFromGosseractBoundingBoxes(boxes []gosseract.BoundingBox) []Block {
	blocks := make([]Block, 0, len(boxes))

	for _, box := range boxes {
		blocks = append(blocks, newBlockFromGosseractBoundingBox(box))
	}

	return blocks
}
