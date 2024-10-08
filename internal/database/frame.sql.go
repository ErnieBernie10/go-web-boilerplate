// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: frame.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const getFrame = `-- name: GetFrame :one
select id, title, description, created_at, modified_at, user_id, frame_status
from frame
where id = $1
`

func (q *Queries) GetFrame(ctx context.Context, id uuid.UUID) (Frame, error) {
	row := q.db.QueryRowContext(ctx, getFrame, id)
	var i Frame
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.CreatedAt,
		&i.ModifiedAt,
		&i.UserID,
		&i.FrameStatus,
	)
	return i, err
}

const getFrames = `-- name: GetFrames :many
select id, title, description, created_at, modified_at, user_id, frame_status
from frame
`

func (q *Queries) GetFrames(ctx context.Context) ([]Frame, error) {
	rows, err := q.db.QueryContext(ctx, getFrames)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Frame
	for rows.Next() {
		var i Frame
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.CreatedAt,
			&i.ModifiedAt,
			&i.UserID,
			&i.FrameStatus,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const saveFrame = `-- name: SaveFrame :one
insert into frame (id, title, description, created_at)
values ($1, $2, $3, NOW()) on conflict (id) DO
UPDATE
set title = $2,
  description = $3,
  modified_at = NOW()
RETURNING id
`

type SaveFrameParams struct {
	ID          uuid.UUID
	Title       string
	Description string
}

func (q *Queries) SaveFrame(ctx context.Context, arg SaveFrameParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, saveFrame, arg.ID, arg.Title, arg.Description)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}
