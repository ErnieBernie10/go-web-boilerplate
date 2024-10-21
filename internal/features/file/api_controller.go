package file

import (
	"database/sql"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"framer/internal/api"
	"framer/internal/core"
	"framer/internal/database"
)

func FileApiHandler(r chi.Router) {
	r.Put(api.UploadFileApiPath, uploadRawFileHandler)
	r.Get(api.DownloadFileApiPath, downloadFileHandler)
}

var baseUploadDir = os.Getenv("UPLOAD_DIR")

// @Summary Upload File
// @Description Upload a file
// @Accept raw file
// @Produce json
// @Param filename path string true "Filename"
// @Success 200 {object} api.CreatedResponse
// @Router /api/file [Put]
func uploadRawFileHandler(w http.ResponseWriter, r *http.Request) {
	user := api.GetUser(r)
	filename := chi.URLParam(r, "filename")

	if filename == "" {
		api.HandleError(r, w, errors.Join(core.ErrValidation, errors.New("filename is required")))
		return
	}

	var id uuid.UUID
	err := database.Transactional(r.Context(), database.Db, func(tx *sql.Tx) error {
		var err error
		id, err = UploadFile(r.Context(), database.Service.WithTx(tx), UploadFileCommand{
			FileName: filename,
			UserID:   user.ID,
			Body:     r.Body,
		})
		return err
	})

	if err != nil {
		api.HandleError(r, w, err)
		return
	}

	api.WriteCreatedResponse(w, strings.Replace(api.DownloadFileApiPath, "{id}", id.String(), 1), api.CreatedResponse(id.String()))
}

// Handles file downloads
func downloadFileHandler(w http.ResponseWriter, r *http.Request) {
	user := api.GetUser(r)
	id := chi.URLParam(r, "id")
	userDir := filepath.Join(baseUploadDir, user.ID.String())
	files, err := os.ReadDir(userDir)
	if err != nil {
		api.HandleError(r, w, errors.Join(core.ErrNotFound, err))
		return
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), id) {
			http.ServeFile(w, r, filepath.Join(userDir, file.Name()))
			return
		}
	}

	api.HandleError(r, w, errors.Join(core.ErrNotFound, errors.New("file not found")))
}
