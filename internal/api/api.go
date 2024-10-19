package api

import (
	"framer/internal/core"
	"log/slog"
	"net/http"

	"github.com/go-chi/httplog/v2"
)

func GetLogger(r *http.Request) *slog.Logger {
	return httplog.LogEntry(r.Context())
}


func GetUser(r *http.Request) *core.Claims {
	user, ok := r.Context().Value(core.UserContextKey).(*core.Claims)

	if !ok {
		return nil
	}

	return user
}