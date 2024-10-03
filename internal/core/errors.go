package core

import "errors"

type ValidationError error

func NewValidationError(val string) ValidationError {
	return errors.New(val)
}
