package file

import (
	"context"
	"database/sql"
	"framer/internal/pkg/database"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type DeleteFileCommand struct {
	UserID uuid.UUID
	ID     uuid.UUID
}

type UploadFileCommand struct {
	FileName string
	UserID   uuid.UUID
	Body     []byte
}

func DeleteFile(ctx context.Context, uow *database.UnitOfWork, cmd DeleteFileCommand) error {
	file, err := uow.Queries.GetFileByID(ctx, cmd.ID)
	if err != nil {
		return err
	}

	err = uow.Queries.DeleteFile(ctx, cmd.ID)
	if err != nil {
		return err
	}

	err = os.Remove(filepath.Join(baseUploadDir, cmd.UserID.String(), file.FileName.String))
	if err != nil {
		return err
	}

	return nil
}

func UploadFile(ctx context.Context, uow *database.UnitOfWork, cmd UploadFileCommand) (uuid.UUID, error) {
	uploadDir := filepath.Join(baseUploadDir, cmd.UserID.String())

	id := uuid.New().String()
	filename := id + "_" + cmd.FileName

	// Ensure upload directory exists
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, 0755)
	}

	// Create a file locally to save the uploaded file
	err := os.WriteFile(filepath.Join(uploadDir, filename), cmd.Body, 0644)
	if err != nil {
		return uuid.Nil, err
	}

	err = uow.Queries.CreateFile(ctx, database.CreateFileParams{
		ID:       uuid.MustParse(id),
		FileName: sql.NullString{String: filename, Valid: true},
	})
	if err != nil {
		os.Remove(filepath.Join(uploadDir, filename))
		return uuid.Nil, err
	}

	return uuid.MustParse(id), nil
}
