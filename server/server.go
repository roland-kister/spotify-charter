package server

import (
	"encoding/json"
	"net/http"
	"spotify-charter/db"
	"spotify-charter/model"
	"time"
)

type Server struct {
	Reader *db.Reader
}

func (s *Server) GetPlaylists(w http.ResponseWriter, r *http.Request) {
	dateNow := model.TimeToDatestamp(time.Now())

	chartTracks := s.Reader.GetChartTracksExt(dateNow)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(chartTracks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
