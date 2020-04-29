package bot

import (
	"fmt"
	"geochats/pkg/types"
	"github.com/Arman92/go-tdlib"
	"github.com/boltdb/bolt"
)

func (b *Bot) ActionGroupSetLocation(msg *tdlib.Message) error {
	text := tryTextWithoutCommand(msg)
	return b.store.GetConn().Update(func(tx *bolt.Tx) error {
		group, err := b.store.GetGroup(tx, msg.ChatId)
		if err != nil {
			return fmt.Errorf("can't get group: %v", err)
		}
		if group == nil {
			group = &types.Group{ChatID: msg.ChatId}
		}

		if text == "" {
			if group.Latitude == 0 && group.Longitude == 0 {
				if err := b.sendText(msg, "Сейчас локация не указана"); err != nil {
					return err
				}
			} else {
				if err := b.sendText(msg, "Сейчас координаты такие: %f, %f", group.Latitude, group.Longitude); err != nil {
					return err
				}
			}
			return b.sendText(msg, "Чтобы изменить место, отправьте <pre>%s широта, долгота</pre>", locationCommand)
		}

		lat, long, found := extractCoords(text)
		if !found {
			b.logger.Warnf("can't extract group coords from user string `%s`", text)
			return b.sendText(msg, "Не могу вытащить координаты из строки, сорри")
		}
		group.Latitude = lat
		group.Longitude = long
		if err := b.store.SaveGroup(tx, group); err != nil {
			return fmt.Errorf("can't save group: %v", err)
		}
		return b.sendText(msg, "Местоположение митинга сохранено. Увидеть ее можно тут - https://miting.link/#g:%d", group.ChatID)
	})
}
