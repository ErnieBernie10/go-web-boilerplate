package view

import (
	"framer/internal/core"
	"net/http"
)

func GetTokens(r *http.Request) func() (string, string) {
	return func() (string, string) {
		accessToken, _ := r.Context().Value(core.TokenContextKey).(string)
		refreshToken, _ := r.Context().Value(core.RefreshContextKey).(string)

		return accessToken, refreshToken
	}
}
