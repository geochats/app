package bot

import (
	"fmt"
	"geochats/pkg/types"
	"github.com/Arman92/go-tdlib"
)

func (b *Bot) ActionGroupSetLocation(msg *tdlib.Message, group *types.Group) error {
	text := tryExtractText(msg)
	var found bool
	group.Latitude, group.Longitude, found = extractCoords(text)
	if !found {
		b.logger.Warnf("can't extract group coords from user string `%s`", text)
		return b.sendText(msg, "Не могу вытащить координаты из строки, сорри")
	}
	if err := b.store.UpdateGroup(group); err != nil {
		return fmt.Errorf("can't update group with user coords: %v", err)
	}
	return b.sendText(msg, "Местоположение группы сохранено. Увидеть ее можно тут - https://miting.link/#g:%d", group.ChatID)
}
