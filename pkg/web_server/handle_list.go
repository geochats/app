package web_server

import (
	"fmt"
	"geochats/pkg/types"
	"github.com/boltdb/bolt"
	"net/http"
)

func (s *WebServer) handleList() http.HandlerFunc {
	type respGroup struct {
		ChatID       int64       `json:"id"`
		Title        string      `json:"title"`
		Username     string      `json:"username"`
		Userpic      types.Image `json:"userpic"`
		MembersCount int32       `json:"count"`
		Latitude     float64     `json:"latitude"`
		Longitude    float64     `json:"longitude"`
		Description  string      `json:"description"`
	}
	type respPoint struct {
		ID        string  `json:"id"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}
	type respSpec struct {
		Groups []respGroup `json:"groups"`
		Points []respPoint `json:"points"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		resp := new(respSpec)

		if r.URL.Query().Get("random") != "" {
			f := types.NewRandomFixturer("fake")
			for i := 0; i < 100; i++ {
				p := f.Point()
				resp.Points = append(resp.Points, respPoint{
					ID:        p.PublicID(),
					Latitude:  p.Latitude,
					Longitude: p.Longitude,
				})
			}
			for i := 0; i < 10; i++ {
				g := f.Group()
				resp.Groups = append(resp.Groups, respGroup{
					ChatID:       g.ChatID,
					Title:        g.Title,
					Username:     g.Username,
					Userpic:      g.Userpic,
					MembersCount: g.MembersCount,
					Latitude:     g.Latitude,
					Longitude:    g.Longitude,
					Description:  g.Text,
				})
			}
			s.responseWithSuccessJSON(w, resp)
			return
		}

		points := make([]types.Point, 0)
		groups := make([]types.Group, 0)
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

		resp.Points = make([]respPoint, 0)
		for _, p := range points {
			if p.Published {
				resp.Points = append(resp.Points, respPoint{
					ID:        p.PublicID(),
					Latitude:  p.Latitude,
					Longitude: p.Longitude,
				})
			}
		}
		resp.Groups = make([]respGroup, 0)
		for _, g := range groups {
			if g.Complete() {
				resp.Groups = append(resp.Groups, respGroup{
					ChatID:       g.ChatID,
					Title:        g.Title,
					Username:     g.Username,
					Userpic:      g.Userpic,
					MembersCount: g.MembersCount,
					Latitude:     g.Latitude,
					Longitude:    g.Longitude,
					Description:  g.TextHTML(),
				})
			}
		}

		s.responseWithSuccessJSON(w, resp)
	}
}
