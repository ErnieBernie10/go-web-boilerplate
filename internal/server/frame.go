package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) frameResourceHandler(r chi.Router) {
	r.Get("/", s.getFramesHandler)
	r.Get("/{frameId}", s.getFrameHandler)
}

func (s *Server) getFrameHandler(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) getFramesHandler(w http.ResponseWriter, r *http.Request) {
	fs, err := s.queries.GetFrames(s.ctx)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	dtos := []Frame{}
	for _, f := range fs {
		dtos = append(dtos, Frame{
			ID:          f.ID,
			Title:       f.Title,
			Description: f.Description,
		})
	}

	response, err := json.Marshal(dtos)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(response)
}
