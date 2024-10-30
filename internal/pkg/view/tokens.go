package view

import (
	"context"
	"framer/internal/pkg"
	"net/http"

	"google.golang.org/grpc/metadata"
)

func GetTokens(r *http.Request) func() (string, string) {
	return func() (string, string) {
		accessToken, _ := r.Context().Value(pkg.TokenContextKey).(string)
		refreshToken, _ := r.Context().Value(pkg.RefreshContextKey).(string)

		return accessToken, refreshToken
	}
}

func ContextWithToken(r *http.Request) context.Context {
	accessToken, _ := GetTokens(r)()
	ctx := metadata.AppendToOutgoingContext(r.Context(), "authorization", "Bearer "+accessToken)
	return ctx
}
