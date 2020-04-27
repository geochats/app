package bot

import (
	"fmt"
	"geochats/pkg/types"
	"github.com/Arman92/go-tdlib"
)

func (b *Bot) ActionSingle(flow *Flow, msg *tdlib.Message) error {
	if msg.Content.GetMessageContentEnum() == tdlib.MessageLocationType {
		t := msg.Content.(*tdlib.MessageLocation)
		point, err := b.store.UpdatePoint(&types.Point{
			ChatID:    msg.ChatId,
			Latitude:  t.Location.Latitude,
			Longitude: t.Location.Longitude,
		})
		if err != nil {
			return fmt.Errorf("can't update point from bot: %v", err)
		}
		flow.Status = ""
		text := fmt.Sprintf("Ваш пикет размещен по координатам %f, %f. Перейдите по ссылке, чтобы его увидеть: https://miting.link/#%s",
			point.Latitude,
			point.Longitude,
			point.PublicID())
		inputMsgTxt := tdlib.NewInputMessageText(tdlib.NewFormattedText(text, nil), true, true)
		_, err = b.cl.SendMessage(flow.ChatID, 0, nil, nil, inputMsgTxt)
		if err != nil {
			return fmt.Errorf("can't send picket greeting message: %v", err)
		}
		return nil
	}

	text := "Отправьте мне геолокацию пикета"
	inputMsgTxt := tdlib.NewInputMessageText(tdlib.NewFormattedText(text, nil), true, true)
	_, err := b.cl.SendMessage(flow.ChatID, 0, nil, nil, inputMsgTxt)
	flow.Status = "single"
	return err
}
