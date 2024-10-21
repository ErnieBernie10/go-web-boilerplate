package frame

import (
	"framer/internal/api"
	"framer/internal/view"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func FrameViewHandler(r chi.Router) {
	r.Get(view.FramePath, handleGetFrame)
}

func handleGetFrame(w http.ResponseWriter, r *http.Request) {
	resp := &[]GetFrameResponseDto{}
	status, err := api.ApiClient.Request("GET", api.GetFramesApiPath, nil, &resp, view.GetTokens(r))
	user := api.GetUser(r)
	if err != nil || status != http.StatusOK {
		view.Render(w, r, view.Message(view.Error, "Failed to get frames"), user)
		return
	}

	view.Render(w, r, indexPage(resp), user)
}
