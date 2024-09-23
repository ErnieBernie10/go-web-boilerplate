// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"time"

	"github.com/google/uuid"
)

type Frame struct {
	ID          uuid.UUID
	Title       string
	Description string
	CreatedAt   time.Time
	ModifiedAt  time.Time
}

type SchemaMigration struct {
	Version string
}
