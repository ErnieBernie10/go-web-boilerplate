package view

import (
	"framer/internal/core"
	"framer/internal/view/layout"
	"net/http"

	"github.com/a-h/templ"
)

func Render(w http.ResponseWriter, r *http.Request, c templ.Component, user *core.Claims) {
	if IsHxRequest(r) {
		c.Render(r.Context(), w)
	} else {
		if user != nil {
			layout.Authenticated(&layout.AuthenticatedViewModel{Email: user.Email}, c).Render(r.Context(), w)
		} else {
			layout.Unauthenticated(RegisterPath, LoginPath, c).Render(r.Context(), w)
		}
	}
}

func IsHxRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}
