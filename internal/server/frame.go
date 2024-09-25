package server

import (
	"encoding/json"
	"framer/internal/database"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/thoas/go-funk"
)

func frameResourceHandler(r chi.Router) {
	r.Get("/", getFramesHandler)
	r.Get("/{id}", getFrameHandler)
}

// @Summary Get Frame
// @Description Get Frame by id
// @Accept json
// @Produce json
// @Param id path string true "Frame Id"
// @Success 200
// @Router /frame/{id} [get]
func getFrameHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fs, err := database.Service.GetFrame(ctx, uuid.MustParse(chi.URLParam(r, "id")))
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	dto := mapFrameFromDb(fs)

	response, err := json.Marshal(dto)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	w.Write(response)
}

// @Summary Get Frames
// @Description Get Frames for logged in user
// @Accept json
// @Produce json
// @Success 200
// @Router /frame [get]
func getFramesHandler(w http.ResponseWriter, r *http.Request) {
	fs, err := database.Service.GetFrames(r.Context())
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	dtos := funk.Map(fs, mapFrameFromDb)

	response, err := json.Marshal(dtos)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(response)
}
