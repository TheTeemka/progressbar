package bar

import (
	"os"
	"strings"

	"golang.org/x/term"
)

type Bar interface {
	ToString(percentage float64) string
}
type block struct {
	threshold float64
	char      rune
}
type defaultBar struct {
	numBlock int
	blocks   []block
}

func NewDefaultBar() *defaultBar {
	return &defaultBar{
		numBlock: getNumBlocks(),
		blocks: []block{
			{1.0, '█'},   // Full block
			{0.875, '▉'}, // 7/8 block
			{0.75, '▊'},  // 3/4 block
			{0.625, '▋'}, // 5/8 block
			{0.5, '▌'},   // 1/2 block
			{0.375, '▍'}, // 3/8 block
			{0.25, '▎'},  // 1/4 block
			{0.125, '▏'}, // 1/8 block
		},
	}
}

func (b *defaultBar) ToString(percentage float64) string {
	if percentage < 0 {
		percentage = 0
	}
	if percentage > 100 {
		percentage = 100
	}

	blockWidth := 100 / float64(b.numBlock)
	countFullBlock := int(percentage / blockWidth)

	var bar strings.Builder
	bar.Grow(b.numBlock)
	var count int
	for i := 0; i < countFullBlock; i++ {
		bar.WriteRune(b.blocks[0].char)
		count++
	}

	remainPer := percentage - blockWidth*float64(countFullBlock)
	for _, block := range b.blocks {
		if remainPer >= block.threshold {
			bar.WriteRune(block.char)
			count++
			break
		}
	}
	for i := count; i < b.numBlock; i++ {
		bar.WriteRune(' ')
	}
	return bar.String()
}

func getTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = 70
	}
	return width
}

func getNumBlocks() int {
	width := getTerminalWidth()

	return width - 35
}
