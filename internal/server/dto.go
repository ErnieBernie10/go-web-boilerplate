package server

import "github.com/google/uuid"

type Frame struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}
