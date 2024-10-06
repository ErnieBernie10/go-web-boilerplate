package view

import (
	"framer/internal/auth"
	"framer/internal/view/layout"
	"net/http"
)

func HandleGetIndex(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r)

	if user != nil {
		layout.Authenticated(user, indexPage(user)).Render(r.Context(), w)
		return
	}
	layout.Unauthenticated(RegisterPath, LoginPath, indexPage(user)).Render(r.Context(), w)
}
