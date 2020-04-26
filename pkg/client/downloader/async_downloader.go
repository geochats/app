package downloader

import (
	"fmt"
	"github.com/Arman92/go-tdlib"
	"github.com/cheggaaa/pb/v3"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"geochats/pkg/client"
	"geochats/pkg/types"
)

type AsyncDownloader struct {
	cl              client.AbstractClient
	channelInfo     *types.Group
	publicDir       string
	activeDownloads int64
	m sync.Mutex
}

func NewAsyncDownloader(client client.AbstractClient, channelInfo *types.Group, publicDir string) Downloader {
	return &AsyncDownloader{
		cl:          client,
		channelInfo: channelInfo,
		publicDir:   publicDir,
	}
}

func (e *AsyncDownloader) DownloadChannelFile(src *tdlib.File, dst *string) error {
	if src == nil {
		return nil
	}

	publicPath, fullPath := e.buildChannelFilePaths(src)
	if filepath.Ext(fullPath) != "" {
		if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
			*dst = publicPath
			return nil
		}
	}

	go func() {
		e.m.Lock()
		atomic.AddInt64(&e.activeDownloads, 1)
		e.m.Unlock()
		bar := pb.Full.Start(int(src.Remote.UploadedSize)).
			SetWriter(os.Stdout).
			SetRefreshRate(time.Second).
			Set(pb.Bytes, true)
		bar.Start()
		timeout := 600
		for s := 1; s < timeout; s++ {
			time.Sleep(500 * time.Millisecond)
			file, err := e.cl.DownloadFile(src.Id, 1, 0, 0, false)
			if err != nil {
				log.Warnf("can't resume file downloading: %v", err)
			}
			bar.SetCurrent(int64(file.Local.DownloadedSize))
			if file.Local.IsDownloadingCompleted {
				publicPath, fullPath := e.buildChannelFilePaths(file)
				if err := os.MkdirAll(path.Dir(fullPath), 0777); err != nil {
					log.Errorf("can't create channel files dir: %v", err)
				}
				if err := os.Rename(file.Local.Path, fullPath); err != nil {
					log.Errorf("can't move downloaded channel file: %v", err)
				}
				bar.Finish()
				*dst = publicPath
				e.m.Lock()
				atomic.AddInt64(&e.activeDownloads, -1)
				e.m.Unlock()
				return
			}
		}
		log.Error("downloading incomplete in 60 second")
	}()

	return nil
}

func (e *AsyncDownloader) buildChannelFilePaths(file *tdlib.File) (string, string) {
	repl := strings.NewReplacer(string(os.PathSeparator), "", string(os.PathListSeparator), "")
	dirName := repl.Replace(e.channelInfo.Username)
	publicPath := fmt.Sprintf("/c/%s/files/%s%s", dirName, file.Remote.Id, path.Ext(file.Local.Path))
	fullPath := path.Clean(e.publicDir + publicPath)
	return publicPath, fullPath
}

func (e *AsyncDownloader) WaitAllDownloads() error {
	timeout := 600
	for s := 1; s < timeout; s++ {
		time.Sleep(500 * time.Millisecond)
		e.m.Lock()
		v := atomic.LoadInt64(&e.activeDownloads)
		e.m.Unlock()
		if 0 == v {
			return nil
		}
	}
	return fmt.Errorf("not all downloads finished, remained=%d", atomic.LoadInt64(&e.activeDownloads))
}
