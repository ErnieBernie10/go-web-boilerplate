package server

import (
	"framer/internal/core"
	"net/http"
)

func handleError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case core.ValidationError:
		w.Write([]byte(e.Error()))
	default:
		w.Write([]byte("Internal server error"))
	}
}
