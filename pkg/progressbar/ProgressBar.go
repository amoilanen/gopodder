package progressbar

import (
	"fmt"
	"strings"
)

func NewProgressBar(total int64) *ProgressBar {
	return &ProgressBar{Total: total}
}

type ProgressBar struct {
	Total     int64
	completed int64
}

func (p *ProgressBar) OnProgress(completed int64) {
	p.completed += completed
	if p.Total > 0 {
		progress := float64(p.completed) / float64(p.Total)
		fmt.Printf("\r[%-80s] %.2f%%", strings.Repeat("#", int(progress*80)), progress*100)
		if p.completed >= p.Total {
			fmt.Printf("\r %-100s", strings.Repeat(" ", 100))
			fmt.Printf("\r")
		}
	} else {
		fmt.Printf("\rCompleted %d out of unknown total", p.completed)
	}
}
