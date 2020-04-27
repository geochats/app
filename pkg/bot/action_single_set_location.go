package bot

import (
	"fmt"
	"geochats/pkg/types"
	"github.com/Arman92/go-tdlib"
)

func (b *Bot) ActionSingleSetLocation(msg *tdlib.Message, point *types.Point) error {
	text := tryExtractText(msg)
	var found bool
	point.Latitude, point.Longitude, found = extractCoords(text)
	if !found {
		b.logger.Warnf("can't extract point coords from user string `%s`", text)
		return b.sendText(msg, "Не могу вытащить координаты из строки, сорри")
	}
	if _, err := b.store.UpdatePoint(point); err != nil {
		return fmt.Errorf("can't update point with user coords: %v", err)
	}
	return b.sendText(msg, "Местоположение группы сохранено. Увидеть ее можно тут - https://miting.link/#p:%s", point.PublicID())
}
