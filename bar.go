package progressbar

import (
	"fmt"
	"strings"
	"time"
)

const (
	fullBlock  = '█'
	emptyBlock = ' '
	clearLine  = "\r\033[K" // ANSI escape code to clear line
)

var blocks = []struct {
	threshold float64
	char      rune
}{
	{1.0, '█'},   // Full block
	{0.875, '▉'}, // 7/8 block
	{0.75, '▊'},  // 3/4 block
	{0.625, '▋'}, // 5/8 block
	{0.5, '▌'},   // 1/2 block
	{0.375, '▍'}, // 3/8 block
	{0.25, '▎'},  // 1/4 block
	{0.125, '▏'}, // 1/8 block
}

func ToStringBar(p *ProgressBar) string {
	s := clearLine
	if p.TotalBytes == -1 {
		data := convertData(float64(p.DownBytes))
		t := convertTime(int64(time.Since(p.startTime).Seconds()))

		s = fmt.Sprintf("|%s has been downloaded in %s| %-10s", data, t, convertSpeed(p.downloadSpeed))
	} else {
		s += fmt.Sprintf("%-6.2f%%|", p.percent())
		s += stringBar(p)
		s += fmt.Sprintf("|%-10s %-10s", convertSpeed(p.downloadSpeed), convertTime(int64(p.timeLeft)))
	}
	return s
}

func ToStringResult(p *ProgressBar) string {
	data := convertData(float64(p.DownBytes))
	t := convertTime(int64(time.Since(p.startTime).Seconds()))

	s := fmt.Sprintf("\nFinished! %s has succesfully been downloaded in %s\n", data, t)
	return s
}

func stringBar(p *ProgressBar) string {
	per := p.percent()
	if per < 0 {
		per = 0
	}
	if per > 100 {
		per = 100
	}

	blockWidth := 100 / float64(p.numBlock)
	countFullBlock := int(per / blockWidth)

	var bar strings.Builder
	var count int
	for i := 0; i < countFullBlock; i++ {
		bar.WriteRune(fullBlock)
		count++
	}

	remainPer := per - blockWidth*float64(countFullBlock)
	for _, block := range blocks {
		if remainPer >= block.threshold {
			bar.WriteRune(block.char)
			count++
			break
		}
	}
	for i := count; i < p.numBlock; i++ {
		bar.WriteRune(' ')
	}
	return bar.String()
}

func convertSpeed(speed float64) string {
	switch {
	case speed > (1<<30)/10:
		return fmt.Sprintf("%0.2f GB/s", speed/(1<<30))
	case speed > (1<<20)/10:
		return fmt.Sprintf("%0.2f MB/s", speed/(1<<20))
	case speed > (1<<10)/10:
		return fmt.Sprintf("%0.2f KB/s", speed/(1<<10))
	default:
		return fmt.Sprintf("%0.2f B/s", speed)
	}
}

func convertData(data float64) string {
	switch {
	case data > (1<<30)/10:
		return fmt.Sprintf("%0.2f GB", data/(1<<30))
	case data > (1<<20)/10:
		return fmt.Sprintf("%0.2f MB", data/(1<<20))
	case data > (1<<10)/10:
		return fmt.Sprintf("%0.2f KB", data/(1<<10))
	default:
		return fmt.Sprintf("%0.2f B", data)
	}
}

func convertTime(t int64) string {
	var s strings.Builder

	if t >= 3600 {
		s.WriteString(fmt.Sprintf("%dh ", t/3600))
		t /= 3600
	}

	if t >= 60 {
		s.WriteString(fmt.Sprintf("%dm ", t/60))
		t /= 60
	}

	if t >= 0 {
		s.WriteString(fmt.Sprintf("%ds", t))
	}
	return s.String()
}
