package bot

import (
	"github.com/Arman92/go-tdlib"
)

func (b *Bot) ActionHelp(flow *Flow, _ *tdlib.Message) error {
	text := "Вы можете добавить митинг или одиночный пикет. Чтобы добавить пикет, отправьте команду /single. Чтобы добавить митинг - /group"
	inputMsgTxt := tdlib.NewInputMessageText(tdlib.NewFormattedText(text, nil), true, true)
	_, err := b.cl.SendMessage(flow.ChatID, 0, nil, nil, inputMsgTxt)
	return err
}
