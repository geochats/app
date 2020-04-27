package downloader

import (
	"fmt"
	"geochats/pkg/client"
	"github.com/Arman92/go-tdlib"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"
)

type SyncDownloader struct {
	cl      client.AbstractClient
	baseDir string
	baseUrl string
}

func NewSyncDownloader(client client.AbstractClient, baseDir string, baseUrl string) Downloader {
	return &SyncDownloader{
		cl:      client,
		baseDir: baseDir,
		baseUrl: baseUrl,
	}
}

func (e *SyncDownloader) DownloadChannelFile(fileRef *tdlib.File, dst *string) error {
	const steps = 1
	for s := 0; s < steps; s++ {
		err := e.downloadAttempt(fileRef, dst)
		if err == nil {
			return nil
		}
		pause := time.Second
		if s > 5 {
			pause = 10 * time.Second
		}
		log.Infof("File `%d` on attempt %d/%d after pause `%s`: %v", fileRef.Id, s, steps, pause, err)
		time.Sleep(pause)
	}

	return fmt.Errorf("downloading incomplete in %d steps", steps)
}

func (e *SyncDownloader) downloadAttempt(fileRef *tdlib.File, dst *string) error {
	if fileRef == nil {
		return nil
	}

	publicPath, fullPath := e.buildChannelFilePaths(fileRef)
	if filepath.Ext(fullPath) != "" {
		if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
			*dst = publicPath
			return nil
		}
	}

	file, err := e.cl.DownloadFile(fileRef.Id, 1, 0, 0, false)
	if err != nil {
		return fmt.Errorf("file downloading error: %v", err)
	}

	if file.Local.IsDownloadingCompleted {
		publicPath, fullPath := e.buildChannelFilePaths(file)
		if err := os.MkdirAll(path.Dir(fullPath), 0777); err != nil {
			return fmt.Errorf("can't create channel files dir: %v", err)
		}
		if _, err := os.Stat(file.Local.Path); os.IsNotExist(err) {
			return fmt.Errorf("file not exists after download: %v", err)
		}
		if err := copyFile(file.Local.Path, fullPath); err != nil {
			return fmt.Errorf("can't copy downloaded channel file: %v", err)
		}
		*dst = publicPath
		return nil
	}

	return fmt.Errorf(
		"downloading incomplete: %s / %s = %.0f%%",
		ByteCountDecimal(int64(file.Local.DownloadedSize)),
		ByteCountDecimal(int64(file.Remote.UploadedSize)),
		float64(file.Local.DownloadedSize) / float64(file.Remote.UploadedSize) * 100,
	)
}

func ByteCountDecimal(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%dB", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%cB", float64(b)/float64(div), "kMGTPE"[exp])
}


func copyFile(src, dst string) error {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dst, input, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (e *SyncDownloader) buildChannelFilePaths(file *tdlib.File) (string, string) {
	publicPath := fmt.Sprintf("%s%s%s", e.baseUrl, file.Remote.Id, path.Ext(file.Local.Path))
	fullPath := fmt.Sprintf("%s%s%s", e.baseDir, file.Remote.Id, path.Ext(file.Local.Path))
	return publicPath, fullPath
}

func (e *SyncDownloader) WaitAllDownloads() error {
	return nil
}