package server

import (
	"framer/internal/features/auth"
	"framer/internal/features/home"
	"framer/internal/pkg"

	"github.com/go-chi/chi/v5"
)

type Route string

func viewRouteHandler(r chi.Router) {
	var publicRoutes = r.With(pkg.OptionalUserMiddleware)
	publicRoutes.Group(auth.AuthResourceHandler)
	publicRoutes.Group(home.HomeResourceHandler)
}
