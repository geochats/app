package bot

import (
	"github.com/Arman92/go-tdlib"
)

func (b *Bot) ActionShowCoords(msg *tdlib.Message) error {
	t := msg.Content.(*tdlib.MessageLocation)
	return b.sendText(msg, "Координаты этого места %f, %f", t.Location.Latitude, t.Location.Longitude)
}
