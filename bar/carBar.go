package bar

import "strings"

type carBar struct {
	numBlock int
}

func NewCarBar() *carBar {
	return &carBar{
		numBlock: getNumBlocks(),
	}
}

func (b *carBar) ToString(percentage float64) string {
	if percentage < 0 {
		percentage = 0
	}
	if percentage > 100 {
		percentage = 100
	}

	blockWidth := 100 / float64(b.numBlock)
	countFullBlock := int(percentage / blockWidth)
	if countFullBlock > 0 {
		countFullBlock--
	}
	var builder strings.Builder
	builder.Grow(b.numBlock)
	builder.WriteString(strings.Repeat(" ", b.numBlock-countFullBlock-1))
	builder.WriteRune('🚂')
	builder.WriteString(strings.Repeat("=", countFullBlock))

	return builder.String()
}
