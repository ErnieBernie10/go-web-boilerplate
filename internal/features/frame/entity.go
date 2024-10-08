package frame

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

const TitleMaxLength = 50
const DescriptionMaxLength = 500

type Title string
type Description string
type FrameStatus int32

const (
	Active  FrameStatus = 1
	Deleted FrameStatus = 0
)

type Model struct {
	ID          uuid.UUID
	Title       Title
	Description Description
	CreatedAt   time.Time
	ModifiedAt  time.Time
	UserID      uuid.UUID
	FrameStatus FrameStatus
}

func New(id, userId uuid.UUID, title, description string, frameStatus int32, createdAt, modifiedAt time.Time) (*Model, error) {
	if err := Validate(title, description, frameStatus); err != nil {
		return nil, err
	}

	return &Model{
		ID:          id,
		Title:       Title(title),
		Description: Description(description),
		CreatedAt:   createdAt,
		ModifiedAt:  modifiedAt,
		UserID:      userId,
		FrameStatus: FrameStatus(frameStatus),
	}, nil
}

func Validate(title string, description string, frameStatus int32) error {
	if utf8.RuneCountInString(title) > TitleMaxLength {
		return fmt.Errorf("Title may not be longer than %d", TitleMaxLength)
	}

	if utf8.RuneCountInString(description) > DescriptionMaxLength {
		return fmt.Errorf("Description may not be longer than %d", DescriptionMaxLength)
	}

	return nil
}
