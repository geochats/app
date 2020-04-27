package web_server

import (
	"fmt"
	"geochats/pkg/types"
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

		points, err := s.store.ListPoint()
		if err != nil {
			s.responseWithErrorJSON(w, fmt.Errorf("can't load points: %v", err))
			return
		}
		resp.Points = make([]respPoint, 0)
		for _, p := range points {
			if p.Complete() {
				resp.Points = append(resp.Points, respPoint{
					ID:        p.PublicID(),
					Latitude:  p.Latitude,
					Longitude: p.Longitude,
				})
			}
		}

		groups, err := s.store.ListGroups()
		if err != nil {
			s.responseWithErrorJSON(w, fmt.Errorf("can't load groups: %v", err))
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
					Description:  g.Description,
				})
			}
		}

		s.responseWithSuccessJSON(w, resp)
	}
}
