package core

import (
	"errors"
	"net/http"
	"os"
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
		if os.Getenv("APP_ENV") == string(Development) {
			w.Write([]byte(err.Error()))
		} else {
			w.Write([]byte("Internal server error"))
		}
	}
}
