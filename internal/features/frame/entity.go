package frame

import (
	"errors"
	"fmt"
	"framer/internal/core"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

const TitleMaxLength = 50
const DescriptionMaxLength = 500

type Title string
type Description string
type FrameStatus int16

const (
	Active  FrameStatus = 1
	Deleted FrameStatus = 0
)

type Model struct {
	ID          uuid.UUID
	Title       Title
	Description Description
	FileID      uuid.NullUUID
	FileName    string
	CreatedAt   time.Time
	ModifiedAt  time.Time
	UserID      uuid.UUID
	FrameStatus FrameStatus
}

func CreateTitle(title string) (Title, error) {
	if utf8.RuneCountInString(title) > TitleMaxLength {
		return "", errors.Join(core.ErrValidation, fmt.Errorf("title may not be longer than %d characters", TitleMaxLength))
	}

	if title == "" {
		return "", errors.Join(core.ErrValidation, fmt.Errorf("title may not be empty"))
	}

	return Title(title), nil
}

func CreateDescription(description string) (Description, error) {
	if utf8.RuneCountInString(description) > DescriptionMaxLength {
		return "", errors.Join(core.ErrValidation, fmt.Errorf("description may not be longer than %d characters", DescriptionMaxLength))
	}

	return Description(description), nil
}

func CreateUserID(userId string) (uuid.UUID, error) {
	id, err := uuid.Parse(userId)
	if err != nil {
		return uuid.Nil, errors.Join(core.ErrValidation, errors.New("invalid user id"))
	}

	return id, nil
}

func create(title, description string, userId uuid.UUID, id uuid.NullUUID) (*Model, error) {
	var errs []error

	t, err := CreateTitle(title)
	if err != nil {
		errs = append(errs, err)
	}

	d, err := CreateDescription(description)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	m := &Model{
		Title:       t,
		Description: d,
		CreatedAt:   time.Now(),
		ModifiedAt:  time.Now(),
		UserID:      userId,
		FrameStatus: Active,
	}

	if id.Valid {
		m.ID = id.UUID
	} else {
		m.ID = uuid.New()
	}

	return m, nil
}

func fromDto(dto *saveFrameDto, userId uuid.UUID, id uuid.NullUUID) (*Model, error) {
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

	m := &Model{
		Title:       title,
		Description: description,
		FileID:      dto.FileID,
		FileName:    dto.FileName,
		CreatedAt:   time.Now(),
		ModifiedAt:  time.Now(),
		UserID:      userId,
		FrameStatus: Active,
	}

	if id.Valid {
		m.ID = id.UUID
	} else {
		m.ID = uuid.New()
	}

	return m, nil
}
