package bot

import (
	"github.com/Arman92/go-tdlib"
)

func (b *Bot) ActionShowCoords(msg *tdlib.Message) error {
	loc := new(tdlib.Location)
	if msg.Content.GetMessageContentEnum() == tdlib.MessageLocationType {
		loc = msg.Content.(*tdlib.MessageLocation).Location
	} else if msg.Content.GetMessageContentEnum() == tdlib.MessageVenueType {
		loc = msg.Content.(*tdlib.MessageVenue).Venue.Location
	}
	return b.sendText(msg, "Координаты этого места %f, %f", loc.Latitude, loc.Longitude)
}
