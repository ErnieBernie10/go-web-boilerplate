package frame

import (
	"errors"
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

func CreateTitle(title string) (Title, error) {
	if utf8.RuneCountInString(title) > TitleMaxLength {
		return "", fmt.Errorf("title may not be longer than %d characters", TitleMaxLength)
	}

	return Title(title), nil
}

func CreateDescription(description string) (Description, error) {
	if utf8.RuneCountInString(description) > DescriptionMaxLength {
		return "", fmt.Errorf("description may not be longer than %d characters", DescriptionMaxLength)
	}

	return Description(description), nil
}

func CreateUserID(userId string) (uuid.UUID, error) {
	id, err := uuid.Parse(userId)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user id")
	}

	return id, nil
}

func fromDto(dto *saveFrameDto, userId uuid.UUID) (*Model, error) {
	var errs []error

	title, err := CreateTitle(dto.Title)
	if err != nil {
		errs = append(errs, err)
	}

	description, err := CreateDescription(dto.Description)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return &Model{
		ID:          uuid.New(),
		Title:       title,
		Description: description,
		CreatedAt:   time.Now(),
		ModifiedAt:  time.Now(),
		UserID:      userId,
		FrameStatus: Active,
	}, nil
}
