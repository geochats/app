package bot

import (
	"fmt"
	"geochats/pkg/types"
	"github.com/Arman92/go-tdlib"
	"github.com/jackc/pgx/v4"
)

func (b *Bot) ActionPointSetText(msg *tdlib.Message, tx pgx.Tx, point *types.Point) error {
	text := tryTextWithoutCommand(msg)
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
	if err := b.store.UpdatePoint(tx, point); err != nil {
		return fmt.Errorf("can't save point: %v", err)
	}
	return b.sendText(
		msg,
		"Текст сохранен. Если вы уже опубликовали митинг, командой <code>%s</code>, то увидеть его можно тут - %s",
		publishCommand,
		point.PublicURI())
}
