package bot

import (
	"fmt"
	"geochats/pkg/types"
	"github.com/Arman92/go-tdlib"
	"github.com/boltdb/bolt"
	"strings"
)

func (b *Bot) ActionGroupSetLocation(msg *tdlib.Message) error {
	text := strings.Trim(strings.Replace(tryExtractText(msg), locationCommand, "", 1), " ")
	if text == "" {
		return b.sendText(msg, "Вы, похоже, забыли указать координаты. Отправьте <pre>%s широта,долгота</pre>", locationCommand)
	}
	lat, long, found := extractCoords(text)
	if !found {
		b.logger.Warnf("can't extract group coords from user string `%s`", text)
		return b.sendText(msg, "Не могу вытащить координаты из строки, сорри")
	}
	return b.store.GetConn().Update(func(tx *bolt.Tx) error {
		group, err := b.store.GetGroup(tx, msg.ChatId)
		if err != nil {
			return fmt.Errorf("can't get group: %v", err)
		}
		if group == nil {
			group = &types.Group{ChatID: msg.ChatId}
		}
		group.Latitude = lat
		group.Longitude = long
		if err := b.store.SaveGroup(tx, group); err != nil {
			return fmt.Errorf("can't save group: %v", err)
		}
		return b.sendText(msg, "Местоположение митинга сохранено. Увидеть ее можно тут - https://miting.link/#g:%d", group.ChatID)
	})
}
