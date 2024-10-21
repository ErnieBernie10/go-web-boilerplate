package server

import (
	"framer/internal/features/auth"
	"framer/internal/features/frame"
	"framer/internal/features/home"

	"github.com/go-chi/chi/v5"
)

type Route string

func viewRouteHandler(r chi.Router) {
	var privateRoutes = r.With(AuthGuardMiddleware)
	privateRoutes.Group(frame.FrameViewHandler)

	var publicRoutes = r.With(OptionalUserMiddleware)
	publicRoutes.Group(auth.AuthResourceHandler)
	publicRoutes.Group(home.HomeResourceHandler)
}
