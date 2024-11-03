package api

import (
	core "framer/internal/pkg"
	"net/http"
)

func GetUser(r *http.Request) *core.Claims {
	user, ok := r.Context().Value(core.UserContextKey).(*core.Claims)

	if !ok {
		return nil
	}

	return user
}
