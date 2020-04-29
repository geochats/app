package bot

import (
	"fmt"
	"geochats/pkg/types"
	"github.com/Arman92/go-tdlib"
	"github.com/boltdb/bolt"
)

func (b *Bot) ActionSingleSetLocation(msg *tdlib.Message) error {
	text := tryExtractText(msg)
	lat, long, found := extractCoords(text)
	if !found {
		b.logger.Warnf("can't extract group coords from user string `%s`", text)
		return b.sendText(msg, "Не могу вытащить координаты из строки, сорри")
	}
	return b.store.GetConn().Update(func(tx *bolt.Tx) error {
		point, err := b.store.GetPoint(tx, msg.ChatId)
		if err != nil {
			return fmt.Errorf("can't get point: %v", err)
		}
		if point == nil {
			point = &types.Point{ChatID: msg.ChatId}
		}
		point.Latitude = lat
		point.Longitude = long
		if err := b.store.SavePoint(tx, point); err != nil {
			return fmt.Errorf("can't save point: %v", err)
		}
		return b.sendText(
			msg,
			"Ок. Если вы уже опубликовали пикет, командой %s, то увидеть его можно тут - https://miting.link/#g:%s",
			publishCommand,
			point.PublicID())
	})
}
