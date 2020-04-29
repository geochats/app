package bot

import (
	"fmt"
	"geochats/pkg/types"
	"github.com/Arman92/go-tdlib"
	"github.com/boltdb/bolt"
)

func (b *Bot) ActionSingleChangeVisibility(msg *tdlib.Message, value bool) error {
	return b.store.GetConn().Update(func(tx *bolt.Tx) error {
		point, err := b.store.GetPoint(tx, msg.ChatId)
		if err != nil {
			return fmt.Errorf("can't get point: %v", err)
		}
		if point == nil {
			point = &types.Point{ChatID: msg.ChatId}
		}
		point.Published = value
		if value {
			user, err := b.cl.GetUser(msg.SenderUserId)
			if err != nil {
				return fmt.Errorf("can't load user: %v", err)
			}
			point.Username = user.Username
			point.Name = user.FirstName
		}
		if err := b.store.SavePoint(tx, point); err != nil {
			return fmt.Errorf("can't save point: %v", err)
		}
		if point.Published {
			return b.sendText(msg, "Пикет опубликован. Увидеть его можно тут - https://miting.link/#g:%s", point.PublicID())
		} else {
			return b.sendText(msg, "Пикет скрыт")
		}

	})
}
