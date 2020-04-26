package loaders

import (
	"fmt"
	"github.com/Arman92/go-tdlib"
	log "github.com/sirupsen/logrus"
	"geochats/pkg/client"
	"geochats/pkg/client/downloader"
	"geochats/pkg/types"
	"time"
)

type ChannelInfoLoader struct {
	client  client.AbstractClient
	rootDir string
}

func NewChannelInfoLoader(client client.AbstractClient, rootDir string) *ChannelInfoLoader {
	return &ChannelInfoLoader{
		client:  client,
		rootDir: rootDir,
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
	info.SuperGroupID = sg.Id
	info.Username = sg.Username
	info.RegistrationDate = time.Unix(int64(sg.Date), 0)
	log.Infof("super group loaded by id `%d`", sgt.SupergroupId)
	sgi, err := e.client.GetSupergroupFullInfo(sgt.SupergroupId)
	if err != nil {
		return nil, fmt.Errorf("can't load supergroup: %v", err)
	}
	info.InviteLink = sgi.InviteLink
	if info.InviteLink == "" {
		info.InviteLink = fmt.Sprintf("https://t.me/%s", info.Username)
	}
	info.MembersCount = sgi.MemberCount
	info.Description = sgi.Description

	dl := downloader.NewSyncDownloader(e.client, info, e.rootDir)
	if chat.Photo != nil {
		if chat.Photo.Small != nil {
			if err := dl.DownloadChannelFile(chat.Photo.Big, &info.Userpic); err != nil {
				return nil, fmt.Errorf("can't download chat photo: %v", err)
			}
		}
	}

	log.Infof("super group channelInfo loaded by id `%d`", sgt.SupergroupId)

	return info, nil
}
