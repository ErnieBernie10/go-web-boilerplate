package frame

import (
	"context"
	"database/sql"
	"errors"
	"framer/internal/database"
	"framer/internal/features/file"

	"github.com/google/uuid"
)

type DeleteFrameCommand struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func DeleteFrame(ctx context.Context, q *database.Queries, cmd DeleteFrameCommand) error {
	frame, err := q.GetFrame(ctx, database.GetFrameParams{
		ID:     cmd.ID,
		UserID: cmd.UserID,
	})
	if err != nil {
		return err
	}

	err = q.DeleteFrame(ctx, database.DeleteFrameParams{
		ID:     cmd.ID,
		UserID: cmd.UserID,
	})
	if err != nil {
		return err
	}

	err = file.DeleteFile(ctx, q, file.DeleteFileCommand{
		ID:     frame.FileID.UUID,
		UserID: cmd.UserID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}

	return nil
}

func SaveFrame(ctx context.Context, q *database.Queries, cmd *Model) (uuid.UUID, error) {
	if cmd.FileID.Valid {
		_, err := q.GetFileByID(ctx, cmd.FileID.UUID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return uuid.Nil, errors.New("file not found")
			}
			return uuid.Nil, err
		}
	}

	id, err := q.SaveFrame(ctx, database.SaveFrameParams{
		ID:          cmd.ID,
		Title:       string(cmd.Title),
		Description: string(cmd.Description),
		FrameStatus: int32(cmd.FrameStatus),
		UserID:      cmd.UserID,
		FileID:      cmd.FileID,
	})
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
