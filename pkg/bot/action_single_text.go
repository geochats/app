package bot

import (
	"fmt"
	"geochats/pkg/types"
	"github.com/Arman92/go-tdlib"
	"github.com/boltdb/bolt"
	"strings"
)

func (b *Bot) ActionSingleSetText(msg *tdlib.Message) error {
	textWithCommand := tryExtractText(msg)
	if textWithCommand == "" {
		return b.sendText(msg, "Не вижу ваш текст, извините")
	}
	text := strings.Trim(strings.Replace(textWithCommand, textCommand, "", 1), " ")
	return b.store.GetConn().Update(func(tx *bolt.Tx) error {
		point, err := b.store.GetPoint(tx, msg.ChatId)
		if err != nil {
			return fmt.Errorf("can't get point: %v", err)
		}
		if point == nil {
			point = &types.Point{ChatID: msg.ChatId}
		}
		point.Text = text
		if err := b.store.SavePoint(tx, point); err != nil {
			return fmt.Errorf("can't save point: %v", err)
		}
		return b.sendText(msg, "Ок")
	})
}
