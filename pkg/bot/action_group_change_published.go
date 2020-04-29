package bot

import (
	"fmt"
	"geochats/pkg/types"
	"github.com/Arman92/go-tdlib"
	"github.com/boltdb/bolt"
)

func (b *Bot) ActionGroupChangeVisibilityEnable(msg *tdlib.Message, value bool) error {
	return b.store.GetConn().Update(func(tx *bolt.Tx) error {
		group, err := b.store.GetGroup(tx, msg.ChatId)
		if err != nil {
			return fmt.Errorf("can't get group: %v", err)
		}
		if group == nil {
			group = &types.Group{ChatID: msg.ChatId}
		}
		group.Published = value
		if value {
			exportedGroup, err := b.ch.Export(msg.ChatId, false)
			if err != nil {
				return fmt.Errorf("can't export group: %v", err)
			}
			group.Title = exportedGroup.Title
			group.Username = exportedGroup.Username
			group.Userpic = exportedGroup.Userpic
			group.MembersCount = exportedGroup.MembersCount
		}
		if err := b.store.SaveGroup(tx, group); err != nil {
			return fmt.Errorf("can't save group: %v", err)
		}
		if group.Published {
			return b.sendText(msg, "Митинг опубликован. Увидеть его можно тут - https://miting.link/#g:%d", group.ChatID)
		} else {
			return b.sendText(msg, "Митинг скрыт")
		}

	})
}
