package home

import (
	"framer/internal/features/auth"
	"framer/internal/view"
	"framer/internal/view/layout"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func HomeResourceHandler(r chi.Router) {
	r.Get(view.IndexPath, handleGetIndex)
}

func handleGetIndex(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r)

	if user != nil {
		layout.Authenticated(&layout.AuthenticatedViewModel{
			Email: user.Email,
		}, indexPage(user)).Render(r.Context(), w)
		return
	}
	layout.Unauthenticated(view.RegisterPath, view.LoginPath, indexPage(user)).Render(r.Context(), w)
}
