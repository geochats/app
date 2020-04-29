package bot

import (
	"encoding/json"
	"fmt"
	"geochats/pkg/client"
	"geochats/pkg/downloader"
	"geochats/pkg/loaders"
	"geochats/pkg/storage"
	"github.com/Arman92/go-tdlib"
	"github.com/sirupsen/logrus"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
)

const (
	startCommand = "/start"
	locationCommand = "/place"
	textCommand = "/text"
	publishCommand = "/show"
	hideCommand = "/hide"
)

type Bot struct {
	cl     client.AbstractClient
	dl     downloader.Downloader
	store  storage.Storage
	ch     *loaders.ChannelInfoLoader
	logger *logrus.Logger
	me *tdlib.User
}

type Flow struct {
	ChatID int64
	Status string
}

func New(cl client.AbstractClient, store storage.Storage, ch *loaders.ChannelInfoLoader, downloader downloader.Downloader, logger *logrus.Logger) *Bot {
	return &Bot{
		cl:     cl,
		store:  store,
		ch:     ch,
		dl:     downloader,
		logger: logger,
	}
}

func (b *Bot) Run() (err error) {
	b.me, err = b.cl.GetMe()
	if err != nil {
		return fmt.Errorf("can't get tg bot info: %v", err)
	}
	js, _ := json.MarshalIndent(b.me, "", "  ")
	b.logger.Infof("Bot info: %s", string(js))

	eventFilter := func(msg *tdlib.TdMessage) bool {
		updateMsg := (*msg).(*tdlib.UpdateNewMessage)
		return !updateMsg.Message.IsOutgoing
	}
	receiver := b.cl.AddEventReceiver(&tdlib.UpdateNewMessage{}, eventFilter, 5)
	for newMsg := range receiver.Chan {
		update := (newMsg).(*tdlib.UpdateNewMessage)
		b.logger.WithField("chatID", update.Message.ChatId).Infof("new telegram message `%s`", tryExtractText(update.Message))
		go func(update *tdlib.UpdateNewMessage) {
			lg := b.logger.WithField("chatID", update.Message.ChatId).WithField("msgID", update.Message.Id)
			err := b.Process(update.Message)
			if err != nil {
				lg.Errorf("can't process tg update: %v", err)
				_ = b.sendText(update.Message, "Что-то пошло не так :( Попробуйте чуть позже, пожалуйста.")
			}
		}(update)
	}
	return nil
}

func (b *Bot) Process(msg *tdlib.Message) error {
	defer func() {
		if r := recover(); r != nil {
			b.logger.Errorf("Panic in bot.Process():%#v\n%s", r, string(debug.Stack()))
		}
	}()

	chat, err := b.cl.GetChat(msg.ChatId)
	if err != nil {
		return fmt.Errorf("can't get processed message chat: %v", err)
	}

	text := tryExtractText(msg)
	switch chat.Type.GetChatTypeEnum() {
	case tdlib.ChatTypePrivateType:
		if err := b.processPrivateMessage(msg, text, chat); err != nil {
			return fmt.Errorf("can't process private tg message: %v", err)
		}
	case tdlib.ChatTypeSupergroupType:
		if err := b.processGroupMessage(msg, text, chat); err != nil {
			return fmt.Errorf("can't process group tg message: %v", err)
		}
	default:
		return b.sendText(msg, "Бот умеет работать только в публичных группах или в приватном чате")
	}

	return nil
}

func (b *Bot) processPrivateMessage(msg *tdlib.Message, text string, _ *tdlib.Chat) error {
	switch {
	case msg.Content.GetMessageContentEnum() == tdlib.MessageLocationType || msg.Content.GetMessageContentEnum() == tdlib.MessageVenueType:
		return b.ActionShowCoords(msg)
	case strings.HasPrefix(text, locationCommand):
		return b.ActionSingleSetLocation(msg)
	case strings.HasPrefix(text, textCommand):
		return b.ActionSingleSetText(msg)
	case strings.HasPrefix(text, publishCommand):
		return b.ActionSingleChangeVisibility(msg, true)
	case strings.HasPrefix(text, hideCommand):
		return b.ActionSingleChangeVisibility(msg, false)
	default:
		return b.ActionSingleShowHelp(msg)
	}
}

func (b *Bot) processGroupMessage(msg *tdlib.Message, text string, chat *tdlib.Chat) error {
	admins, err := b.cl.GetChatAdministrators(chat.Id)
	if err != nil {
		return fmt.Errorf("can't get chat admins: %v", err)
	}
	isAdmin := false
	for _, a := range admins.Administrators {
		if a.UserId == msg.SenderUserId {
			isAdmin = true
			break
		}
	}
	if !isAdmin {
		return b.sendText(msg, "Только администратор группы может управлять ботом")
	}

	switch {
	case strings.HasPrefix(text, locationCommand):
		return b.ActionGroupSetLocation(msg)
	case strings.HasPrefix(text, textCommand):
		return b.ActionGroupSetText(msg)
	case strings.HasPrefix(text, publishCommand):
		return b.ActionGroupChangeVisibilityEnable(msg, true)
	case strings.HasPrefix(text, hideCommand):
		return b.ActionGroupChangeVisibilityEnable(msg, false)
	default:
		return b.ActionGroupShowHelp(msg)
	}
}

func (b *Bot) sendText(msg *tdlib.Message, format string, a ...interface{}) error {
	text := fmt.Sprintf(format, a...)
	formatted, err := b.cl.ParseTextEntities(text, tdlib.NewTextParseModeHTML())
	if err != nil {
		return fmt.Errorf("can't parse text: %v\n%s", err, text)
	}
	inputMsgTxt := tdlib.NewInputMessageText(formatted, true, true)
	_, err = b.cl.SendMessage(msg.ChatId, 0, nil, nil, inputMsgTxt)
	return err
}

func tryExtractText(msg *tdlib.Message) string {
	text, ok := msg.Content.(*tdlib.MessageText)
	if ok {
		if text.Text != nil {
			return regexp.MustCompile(`(?m)<[^>]+>`).ReplaceAllLiteralString(text.Text.Text, "")
		}
	}
	return ""
}

func tryTextWithoutCommand(msg *tdlib.Message) string {
	text := tryExtractText(msg)
	commandEndPos := strings.Index(text, " ")
	if commandEndPos == -1 {
		if strings.HasPrefix(text, "/") {
			return ""
		}
		return text
	}
	return strings.Trim(text[commandEndPos:], " ")
}

func extractCoords(text string) (float64, float64, bool) {
	var re = regexp.MustCompile(`(?m)(-?[0-9.]+)[,\s]+(-?[0-9.]+)`)
	matches := re.FindAllStringSubmatch(text, 1)
	if len(matches) == 0 {
		return 0, 0, false
	}
	lat, _ := strconv.ParseFloat(matches[0][1], 64)
	long, _ := strconv.ParseFloat(matches[0][2], 64)
	return lat, long, true
}
