package api

import (
	"framer/internal/features/auth"
	"framer/internal/features/file"
	"framer/internal/features/frame"

	"github.com/go-chi/chi/v5"
)

func apiRouteHandler(r chi.Router) {
	var privateRoutes = r.With(AuthGuardMiddleware)
	privateRoutes.Group(frame.FrameApiHandler)

	privateRoutes.Group(file.FileApiHandler)

	r.Group(auth.AuthApiResourceHandler)
}
