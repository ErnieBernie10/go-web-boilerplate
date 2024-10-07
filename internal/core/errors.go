package core

import (
	"errors"
	"net/http"
)

type ValidationError error

func NewValidationError(val string) ValidationError {
	return errors.New(val)
}

func HandleError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case ValidationError:
		w.Write([]byte(e.Error()))
	default:
		w.Write([]byte("Internal server error"))
	}
}
