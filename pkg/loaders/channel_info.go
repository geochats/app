package loaders

import (
	"fmt"
	"geochats/pkg/client"
	"geochats/pkg/downloader"
	"geochats/pkg/types"
	"github.com/Arman92/go-tdlib"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

type ChannelInfoLoader struct {
	client  client.AbstractClient
	baseDir string
	baseUrl string
}

func NewChannelInfoLoader(client client.AbstractClient, baseDir string, baseUrl string) *ChannelInfoLoader {
	return &ChannelInfoLoader{
		client:  client,
		baseDir: baseDir,
		baseUrl: baseUrl,
	}
}

func (e *ChannelInfoLoader) Export(chatID int64, withImage bool) (*types.Group,  error) {
	chat, err := e.client.GetChat(chatID)
	if err != nil {
		return nil, fmt.Errorf("can't get chat from tg: %v", err)
	}
	info := &types.Group{
		ChatID: chat.Id,
		Title: chat.Title,
	}
	sgt, ok := chat.Type.(*tdlib.ChatTypeSupergroup)
	if !ok {
		return nil, fmt.Errorf("can't cast chat to supergroup")
	}
	sg, err := e.client.GetSupergroup(sgt.SupergroupId)
	if err != nil {
		return nil, fmt.Errorf("can't load supergroup: %v", err)
	}
	if sg.IsChannel {
		return nil, fmt.Errorf("it's a channel, not a group")
	}
	info.Username = sg.Username

	sgi, err := e.client.GetSupergroupFullInfo(sgt.SupergroupId)
	if err != nil {
		return nil, fmt.Errorf("can't load supergroup: %v", err)
	}
	info.MembersCount = sgi.MemberCount
	info.Text = sgi.Description

	if withImage {
		repl := strings.NewReplacer(string(os.PathSeparator), "", string(os.PathListSeparator), "")
		dirName := repl.Replace(info.Username)
		dl := downloader.NewSyncDownloader(e.client, fmt.Sprintf("%s/%s/", e.baseDir, dirName), fmt.Sprintf("%s/%s/", e.baseDir, dirName))
		if chat.Photo != nil {
			if chat.Photo.Small != nil {
				if err := dl.DownloadChannelFile(chat.Photo.Big, &info.Userpic.Path); err != nil {
					return nil, fmt.Errorf("can't download chat photo: %v", err)
				}
			}
		}
	}

	log.Debugf("super group channelInfo loaded by id `%d`", sgt.SupergroupId)

	return info, nil
}
