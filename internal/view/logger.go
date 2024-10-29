package view

import (
	"net/http"

	"go.uber.org/zap"
)

type LoggerContextKey string

const LoggerKey LoggerContextKey = "httpLogger"

func GetLogger(r *http.Request) *zap.Logger {
	return r.Context().Value(LoggerKey).(*zap.Logger)
}
