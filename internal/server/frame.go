package server

import (
	"encoding/json"
	"framer/internal/database"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func frameResourceHandler(r chi.Router) {
	r.Get("/", getFramesHandler)
	r.Get("/{frameId}", getFrameHandler)
}

func getFrameHandler(w http.ResponseWriter, r *http.Request) {

}

func getFramesHandler(w http.ResponseWriter, r *http.Request) {
	fs, err := database.Service.GetFrames(r.Context())
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
