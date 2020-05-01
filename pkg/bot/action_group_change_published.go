package bot

import (
	"fmt"
	"geochats/pkg/types"
	"github.com/Arman92/go-tdlib"
	"github.com/jackc/pgx/v4"
)

func (b *Bot) ActionPointChangeVisibilityEnable(msg *tdlib.Message, tx pgx.Tx, point *types.Point, value bool) error {
	point.Published = value
	if value && !point.IsSingle {
		exportedGroup, err := b.ch.Export(msg.ChatId, false)
		if err != nil {
			return fmt.Errorf("can't export point: %v", err)
		}
		point.Username = exportedGroup.Username
		point.MembersCount = exportedGroup.MembersCount
	}
	if err := b.store.UpdatePoint(tx, point); err != nil {
		return fmt.Errorf("can't save point: %v", err)
	}
	if point.Published {
		return b.sendText(msg, "Митинг опубликован. Увидеть его можно тут - %s", point.PublicURI())
	} else {
		return b.sendText(msg, "Митинг скрыт")
	}
}
