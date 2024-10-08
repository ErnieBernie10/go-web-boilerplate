package view

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/httplog/v2"
)

func GetLogger(r *http.Request) *slog.Logger {
	return httplog.LogEntry(r.Context())
}
