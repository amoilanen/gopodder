package http

import (
	"io"

	"github.com/amoilanen/gopodder/pkg/progressbar"
)

func NewDownloadProgressBar(total int64) *DownloadProgressBar {
	baseProgressBar := progressbar.NewProgressBar(total)
	return &DownloadProgressBar{
		ProgressBar: *baseProgressBar,
	}
}

type DownloadProgressBar struct {
	io.Writer
	progressbar.ProgressBar
}

func (p *DownloadProgressBar) Write(b []byte) (n int, err error) {
	n = len(b)
	completed := int64(len(b))
	p.OnProgress(completed)
	return n, nil
}
