package server

import (
	"framer/internal/features/auth"
	"framer/internal/features/frame"
	"framer/internal/pkg"

	"github.com/go-chi/chi/v5"
)

func apiRouteHandler(r chi.Router) {
	var privateRoutes = r.With(pkg.AuthGuardMiddleware)
	privateRoutes.Group(frame.FrameResourceHandler)

	r.Group(auth.AuthApiResourceHandler)
}
