package server

import (
	"framer/internal/view/index"
	"framer/internal/view/layout"
	"framer/internal/view/login"
	"framer/internal/view/register"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
)

func viewRouteHandler(r chi.Router) {
	r.Handle(index.Path, templ.Handler(layout.UnauthenticatedLayout(index.Page())))
	r.Handle(login.Path, templ.Handler(layout.UnauthenticatedLayout(login.Page())))
	r.Handle(register.Path, templ.Handler(layout.UnauthenticatedLayout(register.Page())))

	r.Post(login.Path, login.HandleLogin)
	r.Post(register.Path, register.HandleRegister)
}
