package frame

import (
	"encoding/json"
	"errors"
	"framer/internal/pkg"
	"framer/internal/pkg/api"
	"framer/internal/pkg/database"
	"framer/internal/pkg/util"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func FrameApiHandler(r chi.Router) {
	r.Get(api.GetFramesApiPath, getFramesHandler)
	r.Get(api.GetFrameApiPath, getFrameHandler)
	r.Post(api.PostFrameApiPath, postFrameHandler)
	r.Put(api.PutFrameApiPath, putFrameHandler)
	r.Delete(api.DeleteFrameApiPath, deleteFrameHandler)
}

type GetFrameResponseDto struct {
	ID          uuid.UUID     `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
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

	e, err := database.Service.Queries.GetFrame(ctx, database.GetFrameParams{
		ID:     id,
		UserID: user.ID,
	})

	if err != nil {
		api.HandleError(r, w, err)
		return
	}

	dto := &GetFrameResponseDto{
		ID:          e.ID,
		Title:       e.Title,
		Description: e.Description,
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
	fs, err := database.Service.Queries.GetFrames(r.Context(), user.ID)
	if err != nil {
		api.HandleError(r, w, err)
		return
	}

	dtos := util.Map(fs, func(e database.GetFramesRow) *GetFrameResponseDto {
		return &GetFrameResponseDto{
			ID:          e.ID,
			Title:       e.Title,
			Description: e.Description,
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
		api.HandleError(r, w, errors.Join(pkg.ErrValidation, err))
		return
	}

	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		api.HandleError(r, w, errors.Join(pkg.ErrValidation, err))
		return
	}

	_, err = database.Service.Queries.GetFrame(r.Context(), database.GetFrameParams{
		ID:     id,
		UserID: user.ID,
	})
	if err != nil {
		api.HandleError(r, w, errors.Join(pkg.ErrNotFound, err))
		return
	}

	entity, err := fromDto(body, user.ID, uuid.NullUUID{UUID: id, Valid: true})
	if err != nil {
		api.HandleError(r, w, errors.Join(pkg.ErrValidation, err))
		return
	}

	uow, err := database.NewUnitOfWork()
	defer uow.Rollback()
	if err != nil {
		api.HandleError(r, w, errors.Join(err, errors.New("failed to save frame")))
		return
	}

	id, err = SaveFrame(r.Context(), uow, entity)
	if err != nil {
		api.HandleError(r, w, errors.Join(err, errors.New("failed to save frame")))
		return
	}
	entity.ID = id

	if err != nil {
		api.HandleError(r, w, errors.Join(err, errors.New("failed to save frame")))
		return
	}

	uow.Commit()

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
		api.HandleError(r, w, errors.Join(pkg.ErrValidation, err))
		return
	}

	entity, err := fromDto(body, user.ID, uuid.NullUUID{})
	if err != nil {
		api.HandleError(r, w, errors.Join(pkg.ErrValidation, err))
		return
	}

	uow, err := database.NewUnitOfWork()
	defer uow.Rollback()
	if err != nil {
		api.HandleError(r, w, errors.Join(err, errors.New("failed to create frame")))
		return
	}

	id, err := SaveFrame(r.Context(), uow, entity)
	if err != nil {
		api.HandleError(r, w, errors.Join(err, errors.New("failed to create frame")))
		return
	}

	entity.ID = id
	if err != nil {
		api.HandleError(r, w, errors.Join(pkg.ErrNotFound, err))
		return
	}

	uow.Commit()

	api.WriteCreatedResponse(w, strings.Replace(api.GetFrameApiPath, "{id}", entity.ID.String(), 1), api.CreatedResponse(entity.ID.String()))
}

func deleteFrameHandler(w http.ResponseWriter, r *http.Request) {
	user := api.GetUser(r)

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		api.HandleError(r, w, errors.Join(pkg.ErrValidation, err))
		return
	}

	_, err = database.Service.Queries.GetFrame(r.Context(), database.GetFrameParams{
		ID:     id,
		UserID: user.ID,
	})
	if err != nil {
		api.HandleError(r, w, errors.Join(pkg.ErrNotFound, err))
		return
	}

	uow, err := database.NewUnitOfWork()
	defer uow.Rollback()

	err = DeleteFrame(r.Context(), uow, DeleteFrameCommand{
		ID:     id,
		UserID: user.ID,
	})
	if err != nil {
		api.HandleError(r, w, errors.Join(errors.New("failed to delete frame"), err))
		return
	}

	uow.Commit()

	api.WriteOkResponse(w, nil)
}
