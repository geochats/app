package bot

import (
	"fmt"
	"geochats/pkg/types"
	"github.com/Arman92/go-tdlib"
	"github.com/boltdb/bolt"
)

func (b *Bot) ActionSingleSetText(msg *tdlib.Message) error {
	text := tryTextWithoutCommand(msg)
	return b.store.GetConn().Update(func(tx *bolt.Tx) error {
		point, err := b.store.GetPoint(tx, msg.ChatId)
		if err != nil {
			return fmt.Errorf("can't get point: %v", err)
		}
		if point == nil {
			point = &types.Point{ChatID: msg.ChatId}
		}

		if text == "" {
			if point.Text == "" {
				if err := b.sendText(msg, "Сейчас текста нет"); err != nil {
					return err
				}
			} else {
				if err := b.sendText(msg, "Сейчас текст такой:"); err != nil {
					return err
				}
				if err := b.sendText(msg, "<pre>%s</pre>", point.Text); err != nil {
					return err
				}
			}
			return b.sendText(msg, "Чтобы изменить текст, отправьте <pre>%s ваш текст</pre>", textCommand)
		}

		point.Text = text
		if err := b.store.SavePoint(tx, point); err != nil {
			return fmt.Errorf("can't save point: %v", err)
		}
		return b.sendText(msg, "Ок")
	})
}
