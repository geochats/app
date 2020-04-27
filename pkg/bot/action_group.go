package bot

import (
	"fmt"
	"geochats/pkg/types"
	"github.com/Arman92/go-tdlib"
	"strings"
)

func (b *Bot) ActionGroup(flow *Flow, msg *tdlib.Message) error {
	text := tryExtractText(msg)
	switch {
	case strings.HasPrefix(text, "https://t.me/"):
		name := strings.Replace(text, "https://t.me/", "", -1)
		group, err := b.ch.Export(name)
		if err != nil {
			return fmt.Errorf("can't export group by name `%s`: %v", name, err)
		}
		text := fmt.Sprintf("Я вижу группу с названием `%s`. Чтобы добавить ее на карту, пришлите ее расположение.", group.Title)
		inputMsgTxt := tdlib.NewInputMessageText(tdlib.NewFormattedText(text, nil), true, true)
		_, err = b.cl.SendMessage(flow.ChatID, 0, nil, nil, inputMsgTxt)
		flow.Data["group"] = group
		return err
	case msg.Content.GetMessageContentEnum() == tdlib.MessageLocationType:
		t := msg.Content.(*tdlib.MessageLocation)
		group, ok := flow.Data["group"].(*types.Group)
		if !ok {
			return fmt.Errorf("group not found in flow")
		}
		err := b.store.AddGroup(group)
		if err != nil {
			return fmt.Errorf("can't add group from bot: %v", err)
		}
		flow.Status = ""
		text := fmt.Sprintf("Чтобы ваша группа показалась на карте, добавьте в ее описание ссылку `https://miting.link/#%f,%f`",
			t.Location.Latitude,
			t.Location.Longitude)
		inputMsgTxt := tdlib.NewInputMessageText(tdlib.NewFormattedText(text, nil), true, true)
		_, err = b.cl.SendMessage(flow.ChatID, 0, nil, nil, inputMsgTxt)
		if err != nil {
			return fmt.Errorf("can't send picket greeting message: %v", err)
		}
		return nil
	default:
		text := "Отправьте мне публичную ссылку группы. Посмотрите в описании группы. Она выглядит, как t.me/bla-bla-bla."
		inputMsgTxt := tdlib.NewInputMessageText(tdlib.NewFormattedText(text, nil), true, true)
		_, err := b.cl.SendMessage(flow.ChatID, 0, nil, nil, inputMsgTxt)
		flow.Status = "group"
		return err
	}
}
