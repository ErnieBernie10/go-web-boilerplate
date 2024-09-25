package server

import (
	"framer/internal/database"

	"github.com/google/uuid"
)

type Frame struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}

func mapFrameFromDb(f database.Frame) *Frame {
	return &Frame{
		ID:          f.ID,
		Title:       f.Title,
		Description: f.Description,
	}
}
