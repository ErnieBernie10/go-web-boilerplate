package home

import (
	"fmt"
	"framer/internal/api"
	"framer/internal/features/frame"
	"framer/internal/pkg"
	"framer/internal/view"
	"framer/internal/view/layout"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func HomeResourceHandler(r chi.Router) {
	r.Get(view.IndexPath, handleGetIndex)
}

func handleGetIndex(w http.ResponseWriter, r *http.Request) {
	user := pkg.GetUser(r)

	response := &frame.GetFrameDto{}
	api.ApiClient.Request("GET", api.FramesApiPath, nil, response, pkg.GetTokens(r))

	fmt.Println(response)

	if user != nil {
		layout.Authenticated(&layout.AuthenticatedViewModel{
			Email: user.Email,
		}, indexPage(user)).Render(r.Context(), w)
		return
	}
	layout.Unauthenticated(view.RegisterPath, view.LoginPath, indexPage(user)).Render(r.Context(), w)
}
