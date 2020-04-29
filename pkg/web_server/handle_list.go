package web_server

import (
	"fmt"
	"geochats/pkg/types"
	"github.com/boltdb/bolt"
	"net/http"
)

func (s *WebServer) handleList() http.HandlerFunc {
	type respMarker struct {
		ID           string      `json:"id"`
		Username     string      `json:"username"`
		Title        string      `json:"title"`
		Userpic      types.Image `json:"userpic"`
		MembersCount int32       `json:"count"`
		Latitude     float64     `json:"latitude"`
		Longitude    float64     `json:"longitude"`
		Text         string      `json:"description"`
	}
	type respSpec struct {
		Groups []respMarker `json:"groups"`
		Points []respMarker `json:"points"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		resp := new(respSpec)

		points := make([]types.Point, 0)
		groups := make([]types.Group, 0)
		if r.URL.Query().Get("random") != "" {
			// TODO хорошо бы заменить на RandomStorage
			f := types.NewRandomFixturer("fake")
			for i := 0; i < 100; i++ {
				points = append(points, f.Point())
			}
			for i := 0; i < 10; i++ {
				groups = append(groups, f.Group())
			}
		} else {
			err := s.store.GetConn().View(func(tx *bolt.Tx) error {
				var err error
				points, err = s.store.ListPoint(tx)
				if err != nil {
					return fmt.Errorf("can't load points: %v", err)
				}
				groups, err = s.store.ListGroups(tx)
				if err != nil {
					return fmt.Errorf("can't load groups: %v", err)
				}
				return nil
			})
			if err != nil {
				s.responseWithErrorJSON(w, fmt.Errorf("can't load points: %v", err))
				return
			}
		}
		resp.Points = make([]respMarker, 0)
		for _, p := range points {
			if p.Published {
				resp.Points = append(resp.Points, respMarker{
					ID:        p.PublicID(),
					Username:  p.Username,
					Title:     p.Name,
					Latitude:  p.Latitude,
					Longitude: p.Longitude,
					Text:      p.TextHTML(),
				})
			}
		}
		resp.Groups = make([]respMarker, 0)
		for _, g := range groups {
			if g.Complete() {
				resp.Groups = append(resp.Groups, respMarker{
					ID:           fmt.Sprintf("%d", g.ChatID),
					Title:        g.Title,
					Username:     g.Username,
					Userpic:      g.Userpic,
					MembersCount: g.MembersCount,
					Latitude:     g.Latitude,
					Longitude:    g.Longitude,
					Text:         g.TextHTML(),
				})
			}
		}

		s.responseWithSuccessJSON(w, resp)
	}
}
