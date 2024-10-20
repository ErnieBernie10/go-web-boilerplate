package frame

import (
	"database/sql"
	"encoding/json"
	"errors"
	"framer/internal/api"
	"framer/internal/core"
	"framer/internal/database"
	"framer/internal/util"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func FrameResourceHandler(r chi.Router) {
	r.Get(api.GetFramesApiPath, getFramesHandler)
	r.Get(api.GetFrameApiPath, getFrameHandler)
	r.Post(api.PostFrameApiPath, postFrameHandler)
	r.Put(api.PutFrameApiPath, putFrameHandler)
	r.Delete(api.DeleteFrameApiPath, deleteFrameHandler)
}

type GetFrameDto struct {
	ID          uuid.UUID     `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	CreatedAt   time.Time     `json:"createdAt"`
	ModifiedAt  time.Time     `json:"modifiedAt"`
	UserID      uuid.UUID     `json:"userId"`
	FrameStatus int           `json:"frameStatus"`
	FileID      uuid.NullUUID `json:"fileId"`
}

type saveFrameDto struct {
	Title       string        `json:"title"`
	Description string        `json:"description"`
	FileID      uuid.NullUUID `json:"fileId"`
	FileName    string        `json:"fileName"`
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
	user := api.GetUser(r)

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		api.HandleError(r, w, err)
	}

	e, err := database.Service.GetFrame(ctx, database.GetFrameParams{
		ID:     id,
		UserID: user.ID,
	})

	if err != nil {
		api.HandleError(r, w, err)
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
		FileID:      e.FileID,
	}

	api.WriteOkResponse(w, dto)
}

// @Summary Get Frames
// @Description Get Frames for logged in user
// @Accept json
// @Produce json
// @Success 200
// @Router /api/frame [get]
func getFramesHandler(w http.ResponseWriter, r *http.Request) {
	user := api.GetUser(r)
	fs, err := database.Service.GetFrames(r.Context(), user.ID)
	if err != nil {
		api.HandleError(r, w, err)
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
			FileID:      e.FileID,
		}
	})

	api.WriteOkResponse(w, dtos)
}

// @Summary Put Frame
// @Description Put Frame
// @Accept json
// @Produce json
// @Success 200
// @Param frame body postFrameDto true "Frame data"
// @Router /api/frame [put]
func putFrameHandler(w http.ResponseWriter, r *http.Request) {
	user := api.GetUser(r)
	body := &saveFrameDto{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		api.HandleError(r, w, errors.Join(core.ErrValidation, err))
		return
	}

	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		api.HandleError(r, w, errors.Join(core.ErrValidation, err))
		return
	}

	_, err = database.Service.GetFrame(r.Context(), database.GetFrameParams{
		ID:     id,
		UserID: user.ID,
	})
	if err != nil {
		api.HandleError(r, w, errors.Join(core.ErrNotFound, err))
		return
	}

	entity, err := fromDto(body, user.ID, uuid.NullUUID{UUID: id, Valid: true})
	if err != nil {
		api.HandleError(r, w, errors.Join(core.ErrValidation, err))
		return
	}

	err = database.Transactional(r.Context(), database.Db, func(tx *sql.Tx) error {
		id, err := SaveFrame(r.Context(), database.Service.WithTx(tx), entity)
		if err != nil {
			return err
		}
		entity.ID = id
		return nil
	})

	if err != nil {
		api.HandleError(r, w, err)
		return
	}

	api.WriteUpdatedResponse(w, strings.Replace(api.GetFrameApiPath, "{id}", id.String(), 1), api.CreatedResponse(id.String()))
}

// @Summary Post Frame
// @Description Post Frame
// @Accept json
// @Produce json
// @Success 200
// @Param frame body postFrameDto true "Frame data"
// @Router /api/frame [post]
func postFrameHandler(w http.ResponseWriter, r *http.Request) {
	user := api.GetUser(r)
	body := &saveFrameDto{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		api.HandleError(r, w, errors.Join(core.ErrValidation, err))
		return
	}

	entity, err := fromDto(body, user.ID, uuid.NullUUID{})
	if err != nil {
		api.HandleError(r, w, errors.Join(core.ErrValidation, err))
		return
	}

	err = database.Transactional(r.Context(), database.Db, func(tx *sql.Tx) error {
		id, err := SaveFrame(r.Context(), database.Service.WithTx(tx), entity)
		if err != nil {
			return err
		}
		entity.ID = id
		return nil
	})
	if err != nil {
		api.HandleError(r, w, errors.Join(core.ErrNotFound, err))
		return
	}

	api.WriteCreatedResponse(w, strings.Replace(api.GetFrameApiPath, "{id}", entity.ID.String(), 1), api.CreatedResponse(entity.ID.String()))
}

func deleteFrameHandler(w http.ResponseWriter, r *http.Request) {
	user := api.GetUser(r)

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		api.HandleError(r, w, errors.Join(core.ErrValidation, err))
		return
	}

	_, err = database.Service.GetFrame(r.Context(), database.GetFrameParams{
		ID:     id,
		UserID: user.ID,
	})
	if err != nil {
		api.HandleError(r, w, errors.Join(core.ErrNotFound, err))
		return
	}

	err = database.Transactional(r.Context(), database.Db, func(tx *sql.Tx) error {
		return DeleteFrame(r.Context(), database.Service.WithTx(tx), DeleteFrameCommand{
			ID:     id,
			UserID: user.ID,
		})
	})

	if err != nil {
		api.HandleError(r, w, errors.Join(errors.New("failed to delete frame"), err))
		return
	}

	api.WriteOkResponse(w, nil)
}
