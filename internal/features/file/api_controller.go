package file

import (
	"errors"
	"framer/internal/pkg"
	"framer/internal/pkg/api"
	"framer/internal/pkg/database"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
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
		api.HandleError(r, w, errors.Join(pkg.ErrValidation, errors.New("filename is required")))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		api.HandleError(r, w, errors.Join(pkg.ErrValidation, err))
		return
	}

	uow, err := database.NewUnitOfWork()
	if err != nil {
		api.HandleError(r, w, err)
	}

	id, err := UploadFile(r.Context(), uow, UploadFileCommand{
		FileName: filename,
		UserID:   user.ID,
		Body:     body,
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
		api.HandleError(r, w, errors.Join(pkg.ErrNotFound, err))
		return
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), id) {
			http.ServeFile(w, r, filepath.Join(userDir, file.Name()))
			return
		}
	}

	api.HandleError(r, w, errors.Join(pkg.ErrNotFound, errors.New("file not found")))
}
