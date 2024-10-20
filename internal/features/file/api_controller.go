package file

import (
	"database/sql"
	"errors"
	"io"
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

func FileResourceHandler(r chi.Router) {
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

	id := uuid.New().String()
	filename = id + "_" + filename

	uploadDir := filepath.Join(baseUploadDir, user.ID.String())

	// Ensure upload directory exists
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, 0755)
	}

	// Create a file locally to save the uploaded file
	dst, err := os.Create(filepath.Join(uploadDir, filename))
	if err != nil {
		api.HandleError(r, w, err)
		return
	}
	defer dst.Close()

	// Copy the request body (file data) to the destination file
	_, err = io.Copy(dst, r.Body)
	if err != nil {
		api.HandleError(r, w, err)
		os.Remove(filepath.Join(uploadDir, filename))
		return
	}

	err = database.Service.CreateFile(r.Context(), database.CreateFileParams{
		ID:       uuid.MustParse(id),
		FileName: sql.NullString{String: filename, Valid: true},
	})
	if err != nil {
		api.HandleError(r, w, err)
		os.Remove(filepath.Join(uploadDir, filename))
		return
	}

	api.WriteCreatedResponse(w, strings.Replace(api.DownloadFileApiPath, "{id}", id, 1), api.CreatedResponse(id))
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
