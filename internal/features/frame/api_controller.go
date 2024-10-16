package frame

import (
	"encoding/json"
	"framer/internal/api"
	"framer/internal/database"
	"framer/internal/pkg"
	"framer/internal/util"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func FrameResourceHandler(r chi.Router) {
	r.Get(api.GetFramesApiPath, getFramesHandler)
	r.Get(api.GetFrameApiPath, getFrameHandler)
	r.Post(api.PostFrameApiPath, postFrameHandler)
}

type GetFrameDto struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	ModifiedAt  time.Time `json:"modifiedAt"`
	UserID      uuid.UUID `json:"userId"`
	FrameStatus int       `json:"frameStatus"`
}

type postFrameDto struct {
	Title       string `json:"title"`
	Description string `json:"description"`
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
	dto := &GetFrameDto{
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
// @Router /api/frame [get]
func getFramesHandler(w http.ResponseWriter, r *http.Request) {
	fs, err := database.Service.GetFrames(r.Context())
	if err != nil {
		api.HandleError(r, w, err, http.StatusBadRequest)
		return
	}

	dtos := util.Map(fs, func(e database.Frame) *GetFrameDto {
		return &GetFrameDto{
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

// @Summary Post Frame
// @Description Post Frame
// @Accept json
// @Produce json
// @Success 200
// @Param frame body postFrameDto true "Frame data"
// @Router /api/frame [post]
func postFrameHandler(w http.ResponseWriter, r *http.Request) {
	user := pkg.GetUser(r)
	body := &postFrameDto{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		api.HandleError(r, w, err, http.StatusInternalServerError)
		return
	}

	entity, err := fromDto(body, user.ID)
	if err != nil {
		api.HandleError(r, w, err, http.StatusBadRequest)
		return
	}

	id, err := database.Service.SaveFrame(r.Context(), database.SaveFrameParams{
		ID:          entity.ID,
		Title:       string(entity.Title),
		Description: string(entity.Description),
		UserID:      entity.UserID,
		FrameStatus: int32(entity.FrameStatus),
	})

	if err != nil {
		api.HandleError(r, w, err, http.StatusInternalServerError)
		return
	}

	pkg.WriteCreatedResponse(w, id, api.GetFramesApiPath)
}
