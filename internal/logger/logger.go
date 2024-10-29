package logger

import (
	"context"
	"framer/internal/view"
	"net/http"
	"os"

	"go.uber.org/zap"
)

func NewLogger() *zap.Logger {
	var zapLogger *zap.Logger = nil
	var err error = nil
	if os.Getenv("APP_ENV") == "production" {
		zapLogger, err = zap.NewProduction()
	} else {
		zapLogger, err = zap.NewDevelopment()
	}
	if err != nil {
		panic(err)
	}
	zapLogger.Info("logger initialized")
	return zapLogger
}

func HttpLoggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info("HTTP request received",
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()),
				zap.String("remote_addr", r.RemoteAddr),
			)

			ctx := context.WithValue(r.Context(), view.LoggerKey, logger)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
