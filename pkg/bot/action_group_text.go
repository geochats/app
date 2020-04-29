package bot

import (
	"fmt"
	"geochats/pkg/types"
	"github.com/Arman92/go-tdlib"
	"github.com/boltdb/bolt"
)

func (b *Bot) ActionGroupSetText(msg *tdlib.Message) error {
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
			if group.Text == "" {
				if err := b.sendText(msg, "Сейчас текста нет"); err != nil {
					return err
				}
			} else {
				if err := b.sendText(msg, "Сейчас текст такой:"); err != nil {
					return err
				}
				if err := b.sendText(msg, "<pre>%s</pre>", group.Text); err != nil {
					return err
				}
			}
			return b.sendText(msg, "Чтобы изменить текст, отправьте <pre>%s ваш текст</pre>", textCommand)
		}

		group.Text = text
		if err := b.store.SaveGroup(tx, group); err != nil {
			return fmt.Errorf("can't save group: %v", err)
		}
		return b.sendText(
			msg,
			"Текст сохранен. Если вы уже опубликовали митинг, командой <code>%s</code>, то увидеть его можно тут - https://miting.link/#g:%d",
			publishCommand,
			group.ChatID)
	})
}
