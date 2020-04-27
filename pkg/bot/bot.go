package bot

import (
	"encoding/json"
	"fmt"
	"geochats/pkg/client"
	"geochats/pkg/client/downloader"
	"geochats/pkg/collector/loaders"
	"geochats/pkg/storage"
	"github.com/Arman92/go-tdlib"
	"github.com/sirupsen/logrus"
)

type Bot struct {
	cl client.AbstractClient
	dl downloader.Downloader
	store storage.Storage
	ch *loaders.ChannelInfoLoader
	logger *logrus.Logger
	flows map[int64]Flow
}

type Flow struct {
	ChatID int64
	Status string
	Data map[string]interface{}
}

func New(cl client.AbstractClient, store storage.Storage, ch *loaders.ChannelInfoLoader, downloader downloader.Downloader, logger *logrus.Logger) *Bot {
	return &Bot{
		cl: cl,
		store: store,
		ch: ch,
		dl: downloader,
		logger: logger,
		flows: make(map[int64]Flow),
	}
}

func (b *Bot) Run() error {
	eventFilter := func(msg *tdlib.TdMessage) bool {
		updateMsg := (*msg).(*tdlib.UpdateNewMessage)
		return !updateMsg.Message.IsOutgoing
	}
	receiver := b.cl.AddEventReceiver(&tdlib.UpdateNewMessage{}, eventFilter, 5)
	for newMsg := range receiver.Chan {
		update := (newMsg).(*tdlib.UpdateNewMessage)
		lg := b.logger.WithField("chatID", update.Message.ChatId).WithField("msgID", update.Message.Id)
		currentFlow := b.flows[update.Message.ChatId]
		js, _ := json.MarshalIndent(update, "", "  ")
		fmt.Println(string(js))
		lg.WithField("status", currentFlow.Status).Infof("before: %#v", currentFlow)
		err := b.Process(&currentFlow, update.Message)
		lg.WithField("status", currentFlow.Status).Infof("after: %#v", currentFlow)
		if err != nil {
			lg.Errorf("can't process update: %v", err)
		} else {
			b.flows[update.Message.ChatId] = currentFlow
		}
	}
	return nil
}

func (b *Bot) Process(flow *Flow, msg *tdlib.Message) error {
	text := tryExtractText(msg)
	switch {
	case flow.Status == "":
		flow.ChatID = msg.ChatId
		flow.Status = "new"
		flow.Data = make(map[string]interface{})
		return b.ActionHelp(flow, msg)
	case text == "/single" || flow.Status == "single":
		return b.ActionSingle(flow, msg)
	case text == "/group" || flow.Status == "group":
		return b.ActionGroup(flow, msg)
	}


	return nil

	//var errs []error
	//point := &types.Point{
	//	ChatID: msg.ChatId,
	//}
	//switch msg.Content.GetMessageContentEnum() {
	//case tdlib.MessageLocationType:
	//	t := msg.Content.(*tdlib.MessageLocation)
	//	point.Latitude = t.Location.Latitude
	//	point.Longitude = t.Location.Longitude
	//case tdlib.MessageTextType:
	//	//t := msg.Content.(*tdlib.MessageText)
	//case tdlib.MessagePhotoType:
	//	t := msg.Content.(*tdlib.MessagePhoto)
	//	variant := t.Photo.Sizes[len(t.Photo.Sizes)-1]
	//	var p string
	//	if err := b.dl.DownloadChannelFile(variant.Photo, &p); err != nil {
	//		return nil, err
	//	}
	//	point.Photo = types.Image{Width: variant.Width, Height: variant.Height, Path: p}
	//case tdlib.MessagePollType:
	//	return nil, nil
	//case tdlib.MessageVideoType:
	//	return nil, nil
	//case tdlib.MessagePinMessageType:
	//	return nil, nil
	//case tdlib.MessageChatChangePhotoType:
	//	return nil, nil
	//case tdlib.MessageChatChangeTitleType:
	//	return nil, nil
	//case tdlib.MessageUnsupportedType:
	//	b.logger.Warn("tdlib.MessageUnsupportedType")
	//	return nil, nil
	//default:
	//	b.logger.Warnf("unsupported update with content type `%s`", msg.Content.GetMessageContentEnum())
	//}
	//return b.store.UpdatePoint(point)
}

func tryExtractText(msg *tdlib.Message) string {
	text, ok := msg.Content.(*tdlib.MessageText)
	if ok {
		if text.Text != nil {
			return text.Text.Text
		}
	}
	return ""
}
