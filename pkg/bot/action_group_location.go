package bot

import (
	"fmt"
	"geochats/pkg/types"
	"github.com/Arman92/go-tdlib"
	"github.com/jackc/pgx/v4"
)

func (b *Bot) ActionPointSetLocation(msg *tdlib.Message, tx pgx.Tx, point *types.Point) error {
	text := tryTextWithoutCommand(msg)
	if text == "" {
		if point.Latitude == 0 && point.Longitude == 0 {
			if err := b.sendText(msg, "Сейчас локация не указана"); err != nil {
				return err
			}
		} else {
			if err := b.sendText(msg, "Сейчас координаты такие: %f, %f", point.Latitude, point.Longitude); err != nil {
				return err
			}
		}
		return b.sendText(msg, "Чтобы изменить место, отправьте <pre>%s широта, долгота</pre>", locationCommand)
	}

	lat, long, found := extractCoords(text)
	if !found {
		b.logger.Warnf("can't extract point coords from user string `%s`", text)
		return b.sendText(msg, "Не могу вытащить координаты из строки, сорри")
	}
	point.Latitude = lat
	point.Longitude = long
	if err := b.store.UpdatePoint(tx, point); err != nil {
		return fmt.Errorf("can't save point: %v", err)
	}
	return b.sendText(msg, "Местоположение митинга сохранено. Увидеть ее можно тут - %s", point.PublicURI())
}
