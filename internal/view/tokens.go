package view

import (
	"context"
	"framer/internal/core"
	"net/http"

	"google.golang.org/grpc/metadata"
)

func GetTokens(r *http.Request) func() (string, string) {
	return func() (string, string) {
		accessToken, _ := r.Context().Value(core.TokenContextKey).(string)
		refreshToken, _ := r.Context().Value(core.RefreshContextKey).(string)

		return accessToken, refreshToken
	}
}

func ContextWithToken(r *http.Request) context.Context {
	accessToken, _ := GetTokens(r)()
	ctx := metadata.AppendToOutgoingContext(r.Context(), "authorization", "Bearer "+accessToken)
	return ctx
}
