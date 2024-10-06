package server

import (
	"framer/internal/auth"
	"framer/internal/view"

	"github.com/go-chi/chi/v5"
)

type Route string

func viewRouteHandler(r chi.Router) {
	var publicRoutes = r.With(auth.OptionalUserMiddleware)
	publicRoutes.Get(view.IndexPath, view.HandleGetIndex)

	publicRoutes.Get(view.RegisterPath, view.HandleGetRegister)
	publicRoutes.Post(view.RegisterPath, view.HandleRegister)

	publicRoutes.Get(view.LoginPath, view.HandleGetLogin)
	publicRoutes.Post(view.LoginPath, view.HandleLogin)
}
