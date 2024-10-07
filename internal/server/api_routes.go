package server

import (
	"framer/internal/features/auth"
	"framer/internal/features/frame"

	"github.com/go-chi/chi/v5"
)

func apiRouteHandler(r chi.Router) {
	r.Group(frame.FrameResourceHandler)
	r.Group(auth.AuthApiResourceHandler)
}
