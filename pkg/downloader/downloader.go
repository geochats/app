package downloader

import (
	"github.com/Arman92/go-tdlib"
)

type Downloader interface {
	DownloadChannelFile(src *tdlib.File, dst *string) error
	WaitAllDownloads() error
}
