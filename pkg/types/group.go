package types

import (
	"time"
)

type Group struct {
	ChatID           int64     `json:"id"`
	SuperGroupID     int32     `json:"supergroup"`
	Title            string    `json:"title"`
	Username         string    `json:"username"`
	InviteLink       string    `json:"link"`
	Userpic          string    `json:"userpic"`
	MembersCount     int32     `json:"count"`
	RegistrationDate time.Time `json:"date"`
	Coords           []float64 `json:"coords"`
	Description      string    `json:"description"`
}
