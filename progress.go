package progressbar

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/term"
)

type ProgressBar struct {
	TotalBytes    int64
	DownBytes     int64
	numBlock      int
	startTime     time.Time
	lastTime      time.Time
	inBuffer      int64
	downloadSpeed float64
	timeLeft      float64
}

func New(totalBytes int64) *ProgressBar {
	return &ProgressBar{
		TotalBytes:    totalBytes,
		DownBytes:     0,
		numBlock:      getNumBlocks(),
		startTime:     time.Now(),
		lastTime:      time.Now(),
		downloadSpeed: 0,
		timeLeft:      0,
	}
}

// io.Writer
func (p *ProgressBar) Write(data []byte) (int, error) {
	n := len(data)
	p.inBuffer += int64(n)
	p.DownBytes += int64(n)
	dif := time.Since(p.lastTime).Seconds()
	if dif >= 0.5 {
		p.downloadSpeed = float64(p.inBuffer) / dif
		if p.TotalBytes != -1 {
			p.timeLeft = float64(p.TotalBytes-p.DownBytes) / p.downloadSpeed
		}

		p.inBuffer = 0
		p.lastTime = time.Now()
	}
	p.Print()

	return n, nil
}

func (p *ProgressBar) Print() {
	fmt.Printf("\r%s", ToStringBar(p))

	if p.TotalBytes == p.DownBytes {
		fmt.Print(ToStringResult(p))
	}
}

func (p ProgressBar) percent() float64 {
	return float64(p.DownBytes) / float64(p.TotalBytes) * 100
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

	return min(100, width-30)
}
