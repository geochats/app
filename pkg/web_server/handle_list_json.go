package web_server

import (
	"fmt"
	"geochats/pkg/markdown"
	"geochats/pkg/types"
	"net/http"
	"os"
)

func (s *WebServer) handleListJSON() http.HandlerFunc {
	type respMarker struct {
		ID           string  `json:"id"`
		Username     string  `json:"username"`
		Title        string  `json:"title"`
		MembersCount int32   `json:"count"`
		Latitude     float64 `json:"latitude"`
		Longitude    float64 `json:"longitude"`
		Text         string  `json:"description"`
	}
	type respSpec struct {
		Groups []respMarker `json:"groups"`
		Points []respMarker `json:"points"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		resp := new(respSpec)
		points, err := s.store.ListPublishedPoints()
		if err != nil {
			s.responseWithErrorJSON(w, fmt.Errorf("can't get points list: %v", err))
			return
		}

		if os.Getenv("RANDOM_LIST") != "" {
			f := types.NewRandomFixturer("fake")
			for i := 0; i < 100; i++ {
				points = append(points, f.Single())
			}
			for i := 0; i < 10; i++ {
				points = append(points, f.Group())
			}
		}

		resp.Points = make([]respMarker, 0)
		resp.Groups = make([]respMarker, 0)
		for _, p := range points {
			if !p.Published {
				continue
			}
			if p.IsSingle {
				resp.Points = append(resp.Points, respMarker{
					ID:        p.HashedID(),
					Username:  p.Username,
					Title:     "",
					Latitude:  p.Latitude,
					Longitude: p.Longitude,
					Text:      markdown.ToHTML(p.Text),
				})
			} else {
				resp.Groups = append(resp.Groups, respMarker{
					ID:           fmt.Sprintf("%d", p.ChatID),
					Username:     p.Username,
					Latitude:     p.Latitude,
					Longitude:    p.Longitude,
					MembersCount: p.MembersCount,
					Text:         markdown.ToHTML(p.Text),
				})
			}
		}

		s.responseWithSuccessJSON(w, resp)
	}
}
