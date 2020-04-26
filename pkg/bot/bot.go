package bot

import (
	"geochats/pkg/client"
	"geochats/pkg/client/downloader"
	"geochats/pkg/collector/loaders"
	"geochats/pkg/storage"
	"geochats/pkg/types"
	"github.com/Arman92/go-tdlib"
	"github.com/sirupsen/logrus"
)

type Bot struct {
	cl client.AbstractClient
	dl downloader.Downloader
	store storage.Storage
	logger *logrus.Logger
}

func New(cl client.AbstractClient, store storage.Storage, downloader downloader.Downloader, logger *logrus.Logger) *Bot {
	return &Bot{
		cl: cl,
		store: store,
		dl: downloader,
		logger: logger,
	}
}

func (b *Bot) Run() error {
	// Create an filter function which will be used to filter out unwanted tdlib messages
	eventFilter := func(msg *tdlib.TdMessage) bool {
		//updateMsg := (*msg).(*tdlib.UpdateNewMessage)
		//// For example, we want incomming messages from user with below id:
		//if updateMsg.Message.SenderUserId == 41507975 {
		//	return true
		//}
		//return false
		return true
	}

	// Here we can add a receiver to retreive any message type we want
	// We like to get UpdateNewMessage events and with a specific FilterFunc
	receiver := b.cl.AddEventReceiver(&tdlib.UpdateNewMessage{}, eventFilter, 5)
	for newMsg := range receiver.Chan {
		update := (newMsg).(*tdlib.UpdateNewMessage)
		point, err := b.Process(update.Message)
		if err != nil {
			b.logger.Errorf("can't process update: %v", err)
		} else {
			b.logger.Infof("Update %d processed, point is %#v", update.Message.Id, point)
		}
	}
	return nil
}

func (b *Bot) Process(msg *tdlib.Message) (*types.Point, error) {
	var errs []error
	point := types.NewPoint(msg.ChatId)
	switch msg.Content.GetMessageContentEnum() {
	case tdlib.MessageLocationType:
		t := msg.Content.(*tdlib.MessageLocation)
		point.Coords = []float64{t.Location.Latitude, t.Location.Longitude}
	case tdlib.MessageTextType:
		t := msg.Content.(*tdlib.MessageText)
		point.Description, errs = loaders.FormattedTextToMarkdown(t.Text)
		for _, err := range errs {
			b.logger.Warn(err)
		}
	case tdlib.MessagePhotoType:
		t := msg.Content.(*tdlib.MessagePhoto)
		variant := t.Photo.Sizes[len(t.Photo.Sizes)-1]
		var p string
		if err := b.dl.DownloadChannelFile(variant.Photo, &p); err != nil {
			return nil, err
		}
		point.Photo = types.Image{Width: variant.Width, Height: variant.Height, Path: p}
	case tdlib.MessagePollType:
		return nil, nil
	case tdlib.MessageVideoType:
		return nil, nil
	case tdlib.MessagePinMessageType:
		return nil, nil
	case tdlib.MessageChatChangePhotoType:
		return nil, nil
	case tdlib.MessageChatChangeTitleType:
		return nil, nil
	case tdlib.MessageUnsupportedType:
		b.logger.Warn("tdlib.MessageUnsupportedType")
		return nil, nil
	default:
		b.logger.Warn("unsupported update with content type `%s`", msg.Content.GetMessageContentEnum())
	}
	return b.store.UpdatePoint(point)
}
