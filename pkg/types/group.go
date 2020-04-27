package types

type Group struct {
	ChatID           int64
	Title        string
	Username     string
	Userpic      Image
	MembersCount int32
	Latitude     float64
	Longitude    float64
	Description  string
}


func (g *Group) Complete() bool {
	return g.Username != "" && g.Latitude != 0
}