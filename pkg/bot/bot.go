package bot

import (
	"geochats/pkg/client"
	"geochats/pkg/client/downloader"
	"geochats/pkg/collector/loaders"
	"geochats/pkg/storage"
	"github.com/Arman92/go-tdlib"
	"github.com/sirupsen/logrus"
)

type Bot struct {
	cl     client.AbstractClient
	dl     downloader.Downloader
	store  storage.Storage
	ch     *loaders.ChannelInfoLoader
	logger *logrus.Logger
	flows  map[int64]Flow
}

type Flow struct {
	ChatID int64
	Status string
	Data   map[string]interface{}
}

func New(cl client.AbstractClient, store storage.Storage, ch *loaders.ChannelInfoLoader, downloader downloader.Downloader, logger *logrus.Logger) *Bot {
	return &Bot{
		cl:     cl,
		store:  store,
		ch:     ch,
		dl:     downloader,
		logger: logger,
		flows:  make(map[int64]Flow),
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
		b.logger.Infof("new message to bot `%d`", update.Message.Id)
		go func(update *tdlib.UpdateNewMessage) {
			lg := b.logger.WithField("chatID", update.Message.ChatId).WithField("msgID", update.Message.Id)
			currentFlow := b.flows[update.Message.ChatId]
			lg.WithField("status", currentFlow.Status).Infof("before: %#v", currentFlow)
			err := b.Process(&currentFlow, update.Message)
			lg.WithField("status", currentFlow.Status).Infof("after: %#v", currentFlow)
			if err != nil {
				lg.Errorf("can't process update: %v", err)
			} else {
				b.flows[update.Message.ChatId] = currentFlow
			}
		}(update)
	}
	return nil
}

func (b *Bot) Process(flow *Flow, msg *tdlib.Message) error {
	text := tryExtractText(msg)
	switch {
	case text == "/start" || flow.Status == "":
		flow.ChatID = msg.ChatId
		flow.Status = "new"
		flow.Data = make(map[string]interface{})
		return b.ActionHelp(flow, msg)
	case text == "/single":
		return b.ActionSingle(flow, msg)
	case text == "/group":
		return b.ActionGroup(flow, msg)
	case flow.Status == "single":
		return b.ActionSingle(flow, msg)
	case flow.Status == "group":
		return b.ActionGroup(flow, msg)
	default:
		return b.ActionHelp(flow, msg)
	}
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
