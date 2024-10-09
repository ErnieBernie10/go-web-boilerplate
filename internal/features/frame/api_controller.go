package frame

import (
	"encoding/json"
	"framer/internal/api"
	"framer/internal/database"
	"framer/internal/util"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func FrameResourceHandler(r chi.Router) {
	r.Get("/api/frame", getFramesHandler)
	r.Get("/api/frame/{id}", getFrameHandler)
	r.Post("/api/frame", postFrameHandler)
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

func toDto(entity *Model) *getFrameDto {
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

func toEntity(dbModel database.Frame) (*Model, error) {
	return New(dbModel.ID,
		dbModel.UserID,
		dbModel.Title,
		dbModel.Description,
		dbModel.FrameStatus,
		dbModel.CreatedAt,
		dbModel.ModifiedAt)
}

// @Summary Get Frame
// @Description Get Frame by id
// @Accept json
// @Produce json
// @Param id path string true "Frame Id"
// @Success 200
// @Router /api/frame/{id} [get]
func getFrameHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	e, err := database.Service.GetFrame(ctx, uuid.MustParse(chi.URLParam(r, "id")))
	if err != nil {
		api.HandleError(r, w, err, http.StatusBadRequest)
		return
	}
	dto := &getFrameDto{
		ID:          e.ID,
		Title:       e.Title,
		Description: e.Description,
		CreatedAt:   e.CreatedAt,
		ModifiedAt:  e.ModifiedAt,
		UserID:      e.UserID,
		FrameStatus: int(e.FrameStatus),
	}

	response, err := json.Marshal(dto)
	if err != nil {
		api.HandleError(r, w, err, http.StatusInternalServerError)
		return
	}

	w.Write(response)
}

// @Summary Get Frames
// @Description Get Frames for logged in user
// @Accept json
// @Produce json
// @Success 200
func getFramesHandler(w http.ResponseWriter, r *http.Request) {
	fs, err := database.Service.GetFrames(r.Context())
	if err != nil {
		api.HandleError(r, w, err, http.StatusBadRequest)
		return
	}

	dtos := util.Map(fs, func(e database.Frame) *getFrameDto {
		return &getFrameDto{
			ID:          e.ID,
			Title:       e.Title,
			Description: e.Description,
			CreatedAt:   e.CreatedAt,
			ModifiedAt:  e.ModifiedAt,
			UserID:      e.UserID,
			FrameStatus: int(e.FrameStatus),
		}
	})

	jsonStr, err := json.Marshal(dtos)
	if err != nil {
		api.HandleError(r, w, err, http.StatusInternalServerError)
		return
	}

	w.Write(jsonStr)
}

func postFrameHandler(w http.ResponseWriter, r *http.Request) {
	body := &postFrameDto{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		api.HandleError(r, w, err, http.StatusInternalServerError)
		return
	}
}
