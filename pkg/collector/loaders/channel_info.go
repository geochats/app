package loaders

import (
	"fmt"
	"geochats/pkg/client"
	"geochats/pkg/client/downloader"
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

func (e *ChannelInfoLoader) Export(name string) (*types.Group, error) {
	info := &types.Group{}
	chat, err := e.client.SearchPublicChat(name)
	if err != nil {
		return nil, fmt.Errorf("can't find public chat: %v", err)
	}
	if chat == nil {
		return nil, fmt.Errorf("chat not found by name `%s`", name)
	}
	log.Infof("chat `%d` found by name `%s`", chat.Id, name)

	info.ChatID = chat.Id
	info.Title = chat.Title
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
	log.Infof("super group loaded by id `%d`", sgt.SupergroupId)
	sgi, err := e.client.GetSupergroupFullInfo(sgt.SupergroupId)
	if err != nil {
		return nil, fmt.Errorf("can't load supergroup: %v", err)
	}
	info.MembersCount = sgi.MemberCount
	info.Description = sgi.Description

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

	log.Infof("super group channelInfo loaded by id `%d`", sgt.SupergroupId)

	return info, nil
}
