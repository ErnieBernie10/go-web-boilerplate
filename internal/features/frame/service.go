package frame

import (
	"context"
	"database/sql"
	"errors"
	"framer/internal/core"
	"framer/internal/database"
	"framer/internal/features/file"

	"github.com/google/uuid"
)

type DeleteFrameCommand struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func DeleteFrame(ctx context.Context, uow *database.UnitOfWork, cmd DeleteFrameCommand) error {
	frame, err := uow.Queries.GetFrame(ctx, database.GetFrameParams{
		ID:     cmd.ID,
		UserID: cmd.UserID,
	})
	if err != nil {
		return err
	}

	err = uow.Queries.DeleteFrame(ctx, database.DeleteFrameParams{
		ID:     cmd.ID,
		UserID: cmd.UserID,
	})
	if err != nil {
		return err
	}

	err = file.DeleteFile(ctx, uow, file.DeleteFileCommand{
		ID:     frame.FileID.UUID,
		UserID: cmd.UserID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.Join(core.ErrNotFound, errors.New("frame not found"))
		}
		return err
	}

	return nil
}

func SaveFrameWithFile(ctx context.Context, uow *database.UnitOfWork, cmd *Model, f []byte, fileName string) (uuid.UUID, error) {

	fileID, err := file.UploadFile(ctx, uow, file.UploadFileCommand{
		FileName: fileName,
		UserID:   cmd.UserID,
		Body:     f,
	})

	if err != nil {
		return uuid.Nil, err
	}

	cmd.FileID = uuid.NullUUID{
		UUID:  fileID,
		Valid: true,
	}

	id, err := SaveFrame(ctx, uow, cmd)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func SaveFrame(ctx context.Context, uow *database.UnitOfWork, cmd *Model) (uuid.UUID, error) {
	if cmd.FileID.Valid {
		_, err := uow.Queries.GetFileByID(ctx, cmd.FileID.UUID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return uuid.Nil, errors.New("file not found")
			}
			return uuid.Nil, err
		}
	}

	id, err := uow.Queries.SaveFrame(ctx, database.SaveFrameParams{
		ID:          cmd.ID,
		Title:       string(cmd.Title),
		Description: string(cmd.Description),
		FrameStatus: int16(cmd.FrameStatus),
		UserID:      cmd.UserID,
		FileID:      cmd.FileID,
	})
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
