package progressbar

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/TheTeemka/progressbar/bar"
)

const (
	emptyBlock = ' '
	clearLine  = "\r\033[K"
)

type ProgressBar struct {
	TotalBytes    int64
	DownBytes     int64
	startTime     time.Time
	lastTime      time.Time
	inBuffer      int64
	downloadSpeed float64
	timeLeft      float64

	Bar bar.Bar
}

// nil bar is defaultBar
func New(totalBytes int64, b bar.Bar) (*ProgressBar, error) {
	p := ProgressBar{
		TotalBytes:    totalBytes,
		DownBytes:     0,
		startTime:     time.Now(),
		lastTime:      time.Now(),
		downloadSpeed: 0,
		timeLeft:      0,
	}
	if b == nil {
		p.Bar = bar.NewDefaultBar()
	} else {
		p.Bar = b
	}

	return &p, nil
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
	fmt.Print(p)

	return n, nil
}

func (p *ProgressBar) String() string {
	switch p.DownBytes {
	case 0:
		return "\n" + toStringProgress(p)
	case p.TotalBytes:
		return "\n" + toStringResult(p) + "\n"
	default:
		return "\r" + toStringProgress(p)
	}
}

func toStringProgress(p *ProgressBar) string {
	s := clearLine
	if p.TotalBytes == -1 {
		data := convertData(float64(p.DownBytes))
		t := convertTime(int64(time.Since(p.startTime).Seconds()))

		s = fmt.Sprintf("%s has been downloaded in %s| %-10s", data, t, convertSpeed(p.downloadSpeed))
	} else {
		s += fmt.Sprintf("%-7s|", strconv.FormatFloat(p.percent(), byte('f'), 2, 64)+"%") // 8 blocks
		s += toStringBar(p)
		s += fmt.Sprintf("|%-12s %-6s", convertSpeed(p.downloadSpeed), convertTime(int64(p.timeLeft))) //20block
	}
	return s
}

func toStringResult(p *ProgressBar) string {
	data := convertData(float64(p.DownBytes))
	t := convertTime(int64(time.Since(p.startTime).Seconds()))

	s := fmt.Sprintf("\nFinished! %s has succesfully been downloaded in %s", data, t)
	return s
}

func toStringBar(p *ProgressBar) string {
	return p.Bar.ToString(p.percent())
}

func (p ProgressBar) percent() float64 {
	return float64(p.DownBytes) / float64(p.TotalBytes) * 100
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
	case data > (1<<30)/3:
		return fmt.Sprintf("%0.2f GB", data/(1<<30))
	case data > (1<<20)/3:
		return fmt.Sprintf("%0.2f MB", data/(1<<20))
	case data > (1<<10)/3:
		return fmt.Sprintf("%0.2f KB", data/(1<<10))
	default:
		return fmt.Sprintf("%0.2f B", data)
	}
}

func convertTime(t int64) string {
	var s strings.Builder

	if t >= 3600 {
		s.WriteString(fmt.Sprintf("%dh ", t/3600))
		t %= 3600
	}

	if t >= 60 {
		s.WriteString(fmt.Sprintf("%dm ", t/60))
		t %= 60
	}

	if t >= 0 {
		s.WriteString(fmt.Sprintf("%ds", t))
	}
	return s.String()
}
