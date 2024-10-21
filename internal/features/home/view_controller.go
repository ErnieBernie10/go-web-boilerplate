package home

import (
	"framer/internal/api"
	"framer/internal/view"
	"framer/internal/view/layout"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func HomeResourceHandler(r chi.Router) {
	r.Get(view.IndexPath, handleGetIndex)
}

func handleGetIndex(w http.ResponseWriter, r *http.Request) {
	user := api.GetUser(r)

	if user != nil {
		w.Header().Set("HX-Redirect", view.FramePath)
		http.Redirect(w, r, view.FramePath, http.StatusSeeOther)
		return
	}
	layout.Unauthenticated(view.RegisterPath, view.LoginPath, indexPage(user)).Render(r.Context(), w)
}
