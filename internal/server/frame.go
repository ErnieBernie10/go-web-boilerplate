package server

import (
	"encoding/json"
	"framer/internal/database"
	"framer/internal/frame"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/thoas/go-funk"
)

func frameResourceHandler(r chi.Router) {
	r.Get("/", getFramesHandler)
	r.Get("/{id}", getFrameHandler)
	r.Post("/", postFrameHandler)
}

type getFrameDto struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"Description"`
	CreatedAt   time.Time `json:"createdAt"`
	ModifiedAt  time.Time `json:"modifiedAt"`
	UserID      uuid.UUID `json:"userId"`
	FrameStatus int       `json:"frameStatus"`
}

type postFrameDto struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	UserId      string `json:"userId"`
}

func toDto(entity *frame.Model) *getFrameDto {
	return &getFrameDto{
		ID:          entity.ID,
		Title:       string(entity.Title),
		Description: string(entity.Description),
		CreatedAt:   entity.CreatedAt,
		ModifiedAt:  entity.ModifiedAt,
		UserID:      entity.UserID,
		FrameStatus: int(entity.FrameStatus),
	}
}

func toEntity(dbModel database.Frame) (*frame.Model, error) {
	return frame.New(dbModel.ID, dbModel.UserID, dbModel.Title, dbModel.Description, dbModel.FrameStatus, dbModel.CreatedAt, dbModel.ModifiedAt)
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
		handleError(w, err)
		return
	}

	e, err := toEntity(fs)
	if err != nil {
		handleError(w, err)
		return
	}

	// DO LOGIC HERE

	dto := toDto(e)

	response, err := json.Marshal(dto)
	if err != nil {
		handleError(w, err)
		return
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
		handleError(w, err)
		return
	}
	entities := funk.Map(fs, toEntity)

	// DO LOGIC HERE

	dtos := funk.Map(entities, toDto)

	response, err := json.Marshal(dtos)
	if err != nil {
		handleError(w, err)
		return
	}
	w.Write(response)
}

func postFrameHandler(w http.ResponseWriter, r *http.Request) {
	body := &postFrameDto{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		handleError(w, err)
		return
	}
}
