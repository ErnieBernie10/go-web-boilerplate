package server

import (
	"github.com/go-chi/chi/v5"
)

func apiRouteHandler(r chi.Router) {
	r.Route("/frame", frameResourceHandler)
}
